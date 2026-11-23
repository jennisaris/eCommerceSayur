package seclog

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestMaskEmail(t *testing.T) {
	cases := map[string]string{
		"budi@example.com": "bu***@example.com",
		"ab@example.com":   "a*@example.com",
		"a@example.com":    "*@example.com",
		"":                 "",
	}
	for in, want := range cases {
		if got := MaskEmail(in); got != want {
			t.Errorf("MaskEmail(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRepeatedLoginFailuresTriggerAlertAndRespectCooldown(t *testing.T) {
	var webhookHits int32
	alertServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload AlertPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err == nil && payload.Category == "security_alert" {
			atomic.AddInt32(&webhookHits, 1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer alertServer.Close()

	var buf bytes.Buffer
	mon := NewMonitor(&buf, alertServer.URL)
	baseTime := time.Date(2026, 11, 23, 12, 0, 0, 0, time.UTC)
	mon.nowFunc = func() time.Time { return baseTime }

	// First 4 failures must not trigger alert
	for i := 1; i <= 4; i++ {
		if alerted := mon.RecordLoginFailure("203.0.113.55", "victim@example.com", "/api/v1/auth/login", "invalid_credentials"); alerted {
			t.Fatalf("attempt %d should not trigger alert yet", i)
		}
	}

	// 5th failure within window must trigger security_alert
	if alerted := mon.RecordLoginFailure("203.0.113.55", "victim@example.com", "/api/v1/auth/login", "invalid_credentials"); !alerted {
		t.Fatal("expected 5th failed login within 5 minutes to trigger security alert")
	}

	// 6th failure inside cooldown must not spam duplicate alert
	if alerted := mon.RecordLoginFailure("203.0.113.55", "victim@example.com", "/api/v1/auth/login", "invalid_credentials"); alerted {
		t.Fatal("expected cooldown to suppress immediate duplicate alert")
	}

	logs := buf.String()
	if !strings.Contains(logs, `"category":"security_alert"`) {
		t.Fatalf("expected structured security_alert log, got: %s", logs)
	}
	if !strings.Contains(logs, `"email_masked":"vi***@example.com"`) {
		t.Fatalf("expected masked email in log, got: %s", logs)
	}
	if strings.Contains(logs, "victim@example.com") {
		t.Fatalf("raw email must not appear in security log: %s", logs)
	}

	// Wait briefly for async webhook goroutine
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) && atomic.LoadInt32(&webhookHits) == 0 {
		time.Sleep(20 * time.Millisecond)
	}
	if atomic.LoadInt32(&webhookHits) == 0 {
		t.Fatal("expected security alert webhook to receive payload")
	}
}
