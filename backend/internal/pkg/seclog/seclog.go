package seclog

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type Event struct {
	Event    string
	Outcome  string
	IP       string
	UserID   int64
	Email    string
	Method   string
	Path     string
	Reason   string
	TargetID int64
}

type AlertPayload struct {
	Category  string    `json:"category"`
	AlertType string    `json:"alert_type"`
	Key       string    `json:"key"`
	Count     int       `json:"count"`
	WindowSec int       `json:"window_seconds"`
	IP        string    `json:"ip,omitempty"`
	Email     string    `json:"email_masked,omitempty"`
	Path      string    `json:"path,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

type Monitor struct {
	mu             sync.Mutex
	logger         *slog.Logger
	attempts       map[string][]time.Time
	lastAlert      map[string]time.Time
	loginThreshold int
	denyThreshold  int
	rateThreshold  int
	window         time.Duration
	cooldown       time.Duration
	webhookURL     string
	httpClient     *http.Client
	nowFunc        func() time.Time
}

func NewMonitor(out io.Writer, webhookURL string) *Monitor {
	if out == nil {
		out = os.Stdout
	}
	return &Monitor{
		logger:         slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo})),
		attempts:       make(map[string][]time.Time),
		lastAlert:      make(map[string]time.Time),
		loginThreshold: 5,
		denyThreshold:  10,
		rateThreshold:  5,
		window:         5 * time.Minute,
		cooldown:       5 * time.Minute,
		webhookURL:     sanitizeWebhookURL(webhookURL),
		httpClient:     &http.Client{Timeout: 5 * time.Second},
		nowFunc:        time.Now,
	}
}

var defaultMonitor = NewMonitor(os.Stdout, os.Getenv("SECURITY_ALERT_WEBHOOK_URL"))

func sanitizeWebhookURL(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}
	parsed, err := url.ParseRequestURI(trimmed)
	if err != nil {
		return ""
	}
	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "https" && scheme != "http" {
		return ""
	}
	if parsed.Host == "" {
		return ""
	}
	return trimmed
}

// MaskEmail hides most of the local part while preserving domain for incident triage.
func MaskEmail(email string) string {
	cleaned := strings.ToLower(strings.TrimSpace(email))
	if cleaned == "" {
		return ""
	}
	parts := strings.SplitN(cleaned, "@", 2)
	if len(parts) != 2 {
		return "***"
	}
	local, domain := parts[0], parts[1]
	if len(local) <= 1 {
		return "*@" + domain
	}
	if len(local) == 2 {
		return local[:1] + "*@" + domain
	}
	return local[:2] + "***@" + domain
}

func shortHash(value string) string {
	sum := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(value))))
	return hex.EncodeToString(sum[:8])
}

func (m *Monitor) Audit(ev Event) {
	attrs := []any{
		slog.String("category", "security_audit"),
		slog.String("event", ev.Event),
		slog.String("outcome", ev.Outcome),
	}
	if ev.IP != "" {
		attrs = append(attrs, slog.String("ip", ev.IP))
	}
	if ev.UserID > 0 {
		attrs = append(attrs, slog.Int64("user_id", ev.UserID))
	}
	if ev.TargetID > 0 {
		attrs = append(attrs, slog.Int64("target_id", ev.TargetID))
	}
	if masked := MaskEmail(ev.Email); masked != "" {
		attrs = append(attrs, slog.String("email_masked", masked))
	}
	if ev.Method != "" {
		attrs = append(attrs, slog.String("method", ev.Method))
	}
	if ev.Path != "" {
		attrs = append(attrs, slog.String("path", ev.Path))
	}
	if ev.Reason != "" {
		attrs = append(attrs, slog.String("reason", ev.Reason))
	}

	if ev.Outcome == "failure" || ev.Outcome == "denied" || ev.Outcome == "rate_limited" {
		m.logger.Warn("security_event", attrs...)
		return
	}
	m.logger.Info("security_event", attrs...)
}

func (m *Monitor) RecordLoginFailure(ip, email, path, reason string) bool {
	m.Audit(Event{
		Event:   "auth_login",
		Outcome: "failure",
		IP:      ip,
		Email:   email,
		Method:  http.MethodPost,
		Path:    path,
		Reason:  reason,
	})

	alerted := false
	if ip != "" {
		if m.trackAndAlert("login_fail:ip:"+ip, "repeated_login_failures_ip", m.loginThreshold, ip, email, path) {
			alerted = true
		}
	}
	if strings.TrimSpace(email) != "" {
		if m.trackAndAlert("login_fail:email:"+shortHash(email), "repeated_login_failures_account", m.loginThreshold, ip, email, path) {
			alerted = true
		}
	}
	return alerted
}

func (m *Monitor) RecordAccessDenied(ip string, userID int64, method, path, reason string) bool {
	m.Audit(Event{
		Event:   "access_denied",
		Outcome: "denied",
		IP:      ip,
		UserID:  userID,
		Method:  method,
		Path:    path,
		Reason:  reason,
	})
	if ip == "" {
		return false
	}
	return m.trackAndAlert("access_denied:ip:"+ip, "repeated_access_denied", m.denyThreshold, ip, "", path)
}

func (m *Monitor) RecordRateLimit(ip string, userID int64, method, path string) bool {
	m.Audit(Event{
		Event:   "rate_limit_exceeded",
		Outcome: "rate_limited",
		IP:      ip,
		UserID:  userID,
		Method:  method,
		Path:    path,
		Reason:  "http_429",
	})
	if ip == "" {
		return false
	}
	return m.trackAndAlert("rate_limit:ip:"+ip, "repeated_rate_limit_hits", m.rateThreshold, ip, "", path)
}

func (m *Monitor) AuditRoleChange(actorID, targetUserID int64, oldRole, newRole, ip string) {
	m.Audit(Event{
		Event:    "role_change",
		Outcome:  "success",
		IP:       ip,
		UserID:   actorID,
		TargetID: targetUserID,
		Reason:   oldRole + "->" + newRole,
	})
}

func (m *Monitor) trackAndAlert(key, alertType string, threshold int, ip, email, path string) bool {
	now := m.nowFunc()
	cutoff := now.Add(-m.window)

	m.mu.Lock()
	existing := m.attempts[key]
	pruned := existing[:0]
	for _, ts := range existing {
		if ts.After(cutoff) {
			pruned = append(pruned, ts)
		}
	}
	pruned = append(pruned, now)
	m.attempts[key] = pruned
	count := len(pruned)

	shouldAlert := false
	if count >= threshold {
		if last, ok := m.lastAlert[key]; !ok || now.Sub(last) >= m.cooldown {
			m.lastAlert[key] = now
			shouldAlert = true
		}
	}
	webhookURL := m.webhookURL
	m.mu.Unlock()

	if !shouldAlert {
		return false
	}

	payload := AlertPayload{
		Category:  "security_alert",
		AlertType: alertType,
		Key:       key,
		Count:     count,
		WindowSec: int(m.window.Seconds()),
		IP:        ip,
		Email:     MaskEmail(email),
		Path:      path,
		Timestamp: now.UTC(),
	}

	m.logger.Error(
		"security_alert",
		slog.String("category", payload.Category),
		slog.String("alert_type", payload.AlertType),
		slog.Int("count", payload.Count),
		slog.Int("window_seconds", payload.WindowSec),
		slog.String("ip", payload.IP),
		slog.String("email_masked", payload.Email),
		slog.String("path", payload.Path),
	)

	if webhookURL != "" {
		go m.dispatchWebhook(webhookURL, payload)
	}
	return true
}

func (m *Monitor) dispatchWebhook(targetURL string, payload AlertPayload) {
	body, err := json.Marshal(payload)
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, targetURL, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return
	}
	_ = resp.Body.Close()
}

// Package-level convenience helpers using defaultMonitor.
func ConfigureAlertWebhook(rawURL string) {
	defaultMonitor.mu.Lock()
	defaultMonitor.webhookURL = sanitizeWebhookURL(rawURL)
	defaultMonitor.mu.Unlock()
}

func Audit(ev Event) { defaultMonitor.Audit(ev) }

func RecordLoginFailure(ip, email, path, reason string) bool {
	return defaultMonitor.RecordLoginFailure(ip, email, path, reason)
}

func RecordAccessDenied(ip string, userID int64, method, path, reason string) bool {
	return defaultMonitor.RecordAccessDenied(ip, userID, method, path, reason)
}

func RecordRateLimit(ip string, userID int64, method, path string) bool {
	return defaultMonitor.RecordRateLimit(ip, userID, method, path)
}

func AuditRoleChange(actorID, targetUserID int64, oldRole, newRole, ip string) {
	defaultMonitor.AuditRoleChange(actorID, targetUserID, oldRole, newRole, ip)
}
