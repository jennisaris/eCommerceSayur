package http

import (
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"ecommerce-backend/internal/delivery/http/handler"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/hash"
	"ecommerce-backend/internal/pkg/jwt"
	"ecommerce-backend/internal/pkg/midtrans"
	"ecommerce-backend/internal/usecase"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// --- In-memory repositories for deterministic full-router E2E security tests (no DB touched) ---

type memStore struct {
	mu           sync.Mutex
	nextUserID   int64
	nextCartID   int64
	nextOrderID  int64
	nextPayID    int64
	nextTokenID  int64
	users        map[int64]*domain.User
	sessions     map[string]*domain.AuthSession
	tokens       map[int64]*domain.VerificationToken
	carts        map[int64]*domain.Cart
	products     map[int64]*domain.Product
	categories   map[int]*domain.Category
	ordersByID   map[int64]*domain.Order
	ordersByCode map[string]*domain.Order
	payments     map[int64]*domain.Payment
	paymentLogs  map[string]*domain.PaymentLog
	resetURLs    []string
}

func newMemStore() *memStore {
	cat := &domain.Category{ID: 1, Name: "Teknologi", Slug: "teknologi", Status: true}
	prod := &domain.Product{
		ID:           1,
		CategoryID:   1,
		Category:     cat,
		Name:         "Clean Architecture",
		Slug:         "clean-architecture",
		RegularPrice: 120000,
		SalePrice:    100000,
		Weight:       500,
		Stock:        20,
		Status:       "ACTIVE",
	}
	return &memStore{
		nextUserID:   1,
		nextCartID:   1,
		nextOrderID:  1,
		nextPayID:    1,
		nextTokenID:  1,
		users:        map[int64]*domain.User{},
		sessions:     map[string]*domain.AuthSession{},
		tokens:       map[int64]*domain.VerificationToken{},
		carts:        map[int64]*domain.Cart{},
		products:     map[int64]*domain.Product{prod.ID: prod},
		categories:   map[int]*domain.Category{cat.ID: cat},
		ordersByID:   map[int64]*domain.Order{},
		ordersByCode: map[string]*domain.Order{},
		payments:     map[int64]*domain.Payment{},
		paymentLogs:  map[string]*domain.PaymentLog{},
	}
}

func (s *memStore) IsConfigured() bool { return true }
func (s *memStore) SendVerificationEmail(_ context.Context, _, _, _ string) error {
	return nil
}
func (s *memStore) SendPasswordResetEmail(_ context.Context, _, _, resetURL string) error {
	s.mu.Lock()
	s.resetURLs = append(s.resetURLs, resetURL)
	s.mu.Unlock()
	return nil
}

// UserRepository
func (s *memStore) Create(_ context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.ID = s.nextUserID
	s.nextUserID++
	cp := *user
	s.users[cp.ID] = &cp
	return nil
}
func (s *memStore) GetByID(_ context.Context, id int64) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *u
	cp.Roles = append([]domain.Role(nil), u.Roles...)
	return &cp, nil
}
func (s *memStore) GetByEmail(_ context.Context, email string) (*domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	target := strings.ToLower(strings.TrimSpace(email))
	for _, u := range s.users {
		if strings.ToLower(u.Email) == target {
			cp := *u
			cp.Roles = append([]domain.Role(nil), u.Roles...)
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (s *memStore) Update(_ context.Context, user *domain.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *user
	cp.Roles = append([]domain.Role(nil), user.Roles...)
	s.users[cp.ID] = &cp
	return nil
}
func (s *memStore) AssignRole(_ context.Context, userID int64, roleID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	u, ok := s.users[userID]
	if !ok {
		return domain.ErrNotFound
	}
	name := "Customer"
	if roleID == 1 {
		name = "Super Admin"
	}
	u.Roles = append(u.Roles, domain.Role{ID: roleID, Name: name})
	return nil
}
func (s *memStore) CreateVerificationToken(_ context.Context, token *domain.VerificationToken) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	token.ID = s.nextTokenID
	s.nextTokenID++
	cp := *token
	s.tokens[cp.ID] = &cp
	return nil
}
func (s *memStore) GetVerificationToken(_ context.Context, rawToken string, tokenType string) (*domain.VerificationToken, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sum := sha256.Sum256([]byte(rawToken))
	tokenHash := hex.EncodeToString(sum[:])
	for _, tok := range s.tokens {
		if tok.Token == tokenHash && tok.TokenType == tokenType && tok.ExpiresAt.After(time.Now()) {
			cp := *tok
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (s *memStore) DeleteVerificationToken(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tokens, id)
	return nil
}
func (s *memStore) DeleteVerificationTokensByUser(_ context.Context, userID int64, tokenType string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, tok := range s.tokens {
		if tok.UserID == userID && tok.TokenType == tokenType {
			delete(s.tokens, id)
		}
	}
	return nil
}
func (s *memStore) CreateAuthSession(_ context.Context, session *domain.AuthSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := *session
	s.sessions[cp.TokenHash] = &cp
	return nil
}
func (s *memStore) GetActiveAuthSession(_ context.Context, tokenHash string) (*domain.AuthSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[tokenHash]
	if !ok || sess.RevokedAt != nil || !sess.ExpiresAt.After(time.Now()) {
		return nil, domain.ErrInvalidSession
	}
	cp := *sess
	return &cp, nil
}
func (s *memStore) RevokeAuthSession(_ context.Context, tokenHash string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[tokenHash]; ok && sess.RevokedAt == nil {
		now := time.Now()
		sess.RevokedAt = &now
	}
	return nil
}
func (s *memStore) RevokeAllUserSessions(_ context.Context, userID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for _, sess := range s.sessions {
		if sess.UserID == userID && sess.RevokedAt == nil {
			sess.RevokedAt = &now
		}
	}
	return nil
}

// CartRepository
type memCartRepo struct{ s *memStore }

func (r *memCartRepo) GetByUserID(_ context.Context, userID int64) (*domain.Cart, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	c, ok := r.s.carts[userID]
	if !ok {
		c = &domain.Cart{ID: r.s.nextCartID, UserID: userID, Items: []domain.CartItem{}}
		r.s.nextCartID++
		r.s.carts[userID] = c
	}
	cp := *c
	cp.Items = make([]domain.CartItem, len(c.Items))
	for i, it := range c.Items {
		if p, exists := r.s.products[it.ProductID]; exists {
			pcp := *p
			it.Product = &pcp
		}
		cp.Items[i] = it
	}
	return &cp, nil
}
func (r *memCartRepo) Create(_ context.Context, cart *domain.Cart) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	if _, exists := r.s.carts[cart.UserID]; !exists {
		cart.ID = r.s.nextCartID
		r.s.nextCartID++
		r.s.carts[cart.UserID] = &domain.Cart{ID: cart.ID, UserID: cart.UserID, Items: []domain.CartItem{}}
	}
	return nil
}
func (r *memCartRepo) AddItem(_ context.Context, item *domain.CartItem) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, c := range r.s.carts {
		if c.ID == item.CartID {
			for i := range c.Items {
				if c.Items[i].ProductID == item.ProductID {
					c.Items[i].Quantity += item.Quantity
					return nil
				}
			}
			c.Items = append(c.Items, *item)
			return nil
		}
	}
	return domain.ErrNotFound
}
func (r *memCartRepo) UpdateItemQuantity(_ context.Context, cartID, productID int64, qty int) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, c := range r.s.carts {
		if c.ID == cartID {
			for i := range c.Items {
				if c.Items[i].ProductID == productID {
					c.Items[i].Quantity = qty
					return nil
				}
			}
		}
	}
	return nil
}
func (r *memCartRepo) RemoveItem(_ context.Context, cartID, productID int64) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, c := range r.s.carts {
		if c.ID == cartID {
			filtered := c.Items[:0]
			for _, it := range c.Items {
				if it.ProductID != productID {
					filtered = append(filtered, it)
				}
			}
			c.Items = filtered
		}
	}
	return nil
}
func (r *memCartRepo) ClearCart(_ context.Context, cartID int64) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, c := range r.s.carts {
		if c.ID == cartID {
			c.Items = []domain.CartItem{}
		}
	}
	return nil
}

// ProductRepository & CategoryRepository
type memProductRepo struct{ s *memStore }

func (r *memProductRepo) Create(_ context.Context, p *domain.Product) error { return nil }
func (r *memProductRepo) GetByID(_ context.Context, id int64) (*domain.Product, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	p, ok := r.s.products[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *p
	return &cp, nil
}
func (r *memProductRepo) GetByIDForUpdate(_ context.Context, _ *gorm.DB, id int64) (*domain.Product, error) {
	return r.GetByID(context.Background(), id)
}
func (r *memProductRepo) GetBySlug(_ context.Context, slug string) (*domain.Product, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, p := range r.s.products {
		if p.Slug == slug {
			cp := *p
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (r *memProductRepo) GetAll(_ context.Context, _ domain.ProductFilter) ([]domain.Product, int64, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	out := make([]domain.Product, 0, len(r.s.products))
	for _, p := range r.s.products {
		out = append(out, *p)
	}
	return out, int64(len(out)), nil
}
func (r *memProductRepo) Update(_ context.Context, p *domain.Product) error { return nil }
func (r *memProductRepo) UpdateStock(_ context.Context, _ *gorm.DB, _ int64, _ int) error {
	return nil
}
func (r *memProductRepo) Delete(_ context.Context, _ int64) error { return nil }

type memCategoryRepo struct{ s *memStore }

func (r *memCategoryRepo) Create(_ context.Context, _ *domain.Category) error { return nil }
func (r *memCategoryRepo) GetByID(_ context.Context, id int) (*domain.Category, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	c, ok := r.s.categories[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *c
	return &cp, nil
}
func (r *memCategoryRepo) GetBySlug(_ context.Context, slug string) (*domain.Category, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	for _, c := range r.s.categories {
		if c.Slug == slug {
			cp := *c
			return &cp, nil
		}
	}
	return nil, domain.ErrNotFound
}
func (r *memCategoryRepo) GetAll(_ context.Context, onlyActive bool) ([]domain.Category, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	var out []domain.Category
	for _, c := range r.s.categories {
		if !onlyActive || c.Status {
			out = append(out, *c)
		}
	}
	return out, nil
}
func (r *memCategoryRepo) Update(_ context.Context, _ *domain.Category) error { return nil }
func (r *memCategoryRepo) Delete(_ context.Context, _ int) error              { return nil }

// OrderRepository
type memOrderRepo struct{ s *memStore }

func (r *memOrderRepo) CreateWithItems(_ context.Context, order *domain.Order, items []domain.OrderItem) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	var subtotal int64
	totalWeight := 0
	for i := range items {
		p, ok := r.s.products[items[i].ProductID]
		if !ok {
			return domain.ErrNotFound
		}
		if p.Stock < items[i].Quantity {
			return domain.ErrInsufficientStock
		}
		p.Stock -= items[i].Quantity
		price := p.EffectivePrice()
		items[i].ProductName = p.Name
		items[i].Price = price
		items[i].Subtotal = price * int64(items[i].Quantity)
		subtotal += items[i].Subtotal
		totalWeight += p.Weight * items[i].Quantity
	}
	order.ID = r.s.nextOrderID
	r.s.nextOrderID++
	order.ShippingFee = domain.CalculateShippingFee(totalWeight)
	order.TotalAmount = subtotal + order.ShippingFee
	order.Items = items
	cp := *order
	r.s.ordersByID[cp.ID] = &cp
	r.s.ordersByCode[cp.OrderCode] = &cp

	pay := &domain.Payment{
		ID:            r.s.nextPayID,
		OrderID:       cp.ID,
		UserID:        cp.BuyerID,
		PaymentStatus: "pending",
		GrossAmount:   cp.TotalAmount,
	}
	r.s.nextPayID++
	r.s.payments[cp.ID] = pay
	if c, ok := r.s.carts[order.BuyerID]; ok {
		c.Items = []domain.CartItem{}
	}
	return nil
}
func (r *memOrderRepo) GetByID(_ context.Context, id int64) (*domain.Order, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.ordersByID[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *o
	return &cp, nil
}
func (r *memOrderRepo) GetByCode(_ context.Context, code string) (*domain.Order, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.ordersByCode[code]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *o
	return &cp, nil
}
func (r *memOrderRepo) GetByCodeWithDetails(ctx context.Context, code string) (*domain.Order, error) {
	return r.GetByCode(ctx, code)
}
func (r *memOrderRepo) GetAll(_ context.Context, filter domain.OrderFilter) ([]domain.Order, int64, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	var out []domain.Order
	for _, o := range r.s.ordersByID {
		if filter.BuyerID != nil && o.BuyerID != *filter.BuyerID {
			continue
		}
		out = append(out, *o)
	}
	return out, int64(len(out)), nil
}
func (r *memOrderRepo) UpdateStatus(_ context.Context, id int64, status string, tracking string) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.ordersByID[id]
	if !ok {
		return domain.ErrNotFound
	}
	o.Status = status
	if tracking != "" {
		o.TrackingNumber = tracking
	}
	return nil
}
func (r *memOrderRepo) RestockCancelledOrder(_ context.Context, orderID int64) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	o, ok := r.s.ordersByID[orderID]
	if !ok {
		return domain.ErrNotFound
	}
	if o.Status == string(domain.OrderStatusCancelled) {
		return nil
	}
	for _, it := range o.Items {
		if p, exists := r.s.products[it.ProductID]; exists {
			p.Stock += it.Quantity
		}
	}
	o.Status = string(domain.OrderStatusCancelled)
	return nil
}

// PaymentRepository
type memPaymentRepo struct{ s *memStore }

func (r *memPaymentRepo) Create(_ context.Context, p *domain.Payment) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	p.ID = r.s.nextPayID
	r.s.nextPayID++
	cp := *p
	r.s.payments[cp.OrderID] = &cp
	return nil
}
func (r *memPaymentRepo) GetByID(_ context.Context, _ int64) (*domain.Payment, error) {
	return nil, domain.ErrNotFound
}
func (r *memPaymentRepo) GetByOrderID(_ context.Context, orderID int64) (*domain.Payment, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	p, ok := r.s.payments[orderID]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *p
	return &cp, nil
}
func (r *memPaymentRepo) GetByGatewayID(_ context.Context, _ string) (*domain.Payment, error) {
	return nil, domain.ErrNotFound
}
func (r *memPaymentRepo) Update(_ context.Context, p *domain.Payment) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	cp := *p
	r.s.payments[cp.OrderID] = &cp
	return nil
}
func (r *memPaymentRepo) CreateLog(_ context.Context, l *domain.PaymentLog) error {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	cp := *l
	r.s.paymentLogs[cp.EventKey] = &cp
	return nil
}
func (r *memPaymentRepo) GetLogByEventKey(_ context.Context, key string) (*domain.PaymentLog, error) {
	r.s.mu.Lock()
	defer r.s.mu.Unlock()
	l, ok := r.s.paymentLogs[key]
	if !ok {
		return nil, domain.ErrNotFound
	}
	cp := *l
	return &cp, nil
}

// --- Test harness helpers ---

type testEnv struct {
	app        *fiber.App
	store      *memStore
	midtransSK string
}

func newSecurityTestEnv() *testEnv {
	store := newMemStore()
	cartRepo := &memCartRepo{s: store}
	prodRepo := &memProductRepo{s: store}
	catRepo := &memCategoryRepo{s: store}
	orderRepo := &memOrderRepo{s: store}
	payRepo := &memPaymentRepo{s: store}

	jwtSvc := jwt.NewJWTServiceWithDuration("e2e-security-test-secret-minimum-32-chars!!", 15*time.Minute)
	midtransSK := "SB-Mid-server-E2ETestKey999"
	midtransSvc := midtrans.NewClient(midtransSK, "SB-Mid-client-E2ETestKey999", false)

	authUC := usecase.NewAuthUseCaseWithOptions(store, cartRepo, jwtSvc, usecase.AuthOptions{
		RefreshTokenTTL:          24 * time.Hour,
		RequireEmailVerification: false,
		VerificationEmailSender:  store,
		PasswordResetEmailSender: store,
		VerificationBaseURL:      "http://localhost:3000",
		IsProduction:             false,
	})
	catUC := usecase.NewCategoryUseCase(catRepo, nil)
	prodUC := usecase.NewProductUseCase(prodRepo, catRepo, nil)
	cartUC := usecase.NewCartUseCase(cartRepo, prodRepo, nil)
	orderUC := usecase.NewOrderUseCase(orderRepo, cartRepo, prodRepo, store, nil, nil)
	payUC := usecase.NewPaymentUseCase(payRepo, orderRepo, store, nil, midtransSvc)

	app := fiber.New(fiber.Config{
		BodyLimit:               1 * 1024 * 1024,
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1", "0.0.0.0/0"},
		ProxyHeader:             "X-Real-IP",
	})

	SetupRoutes(app, &RouterConfig{
		AuthHandler:       handler.NewAuthHandler(authUC, false, 15*time.Minute, 24*time.Hour),
		CategoryHandler:   handler.NewCategoryHandler(catUC),
		ProductHandler:    handler.NewProductHandler(prodUC),
		CartHandler:       handler.NewCartHandler(cartUC),
		OrderHandler:      handler.NewOrderHandler(orderUC),
		PaymentHandler:    handler.NewPaymentHandler(payUC),
		AdminOrderHandler: handler.NewAdminOrderHandler(orderUC),
		UserRepository:    store,
		JWTService:        jwtSvc,
		AllowedOrigins:    "http://localhost:3000",
		CookieSecure:      false,
		IsProduction:      true,
	})

	return &testEnv{app: app, store: store, midtransSK: midtransSK}
}

type clientJar struct {
	env       *testEnv
	cookies   map[string]string
	csrfToken string
	ip        string
}

func newClientJar(env *testEnv, ip string) *clientJar {
	return &clientJar{env: env, cookies: map[string]string{}, ip: ip}
}

func (j *clientJar) do(t *testing.T, method, path string, body any, origin string, includeCSRF bool) (*http.Response, map[string]any) {
	t.Helper()
	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal failed: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Real-IP", j.ip)
	if origin != "" {
		req.Header.Set("Origin", origin)
	}
	if includeCSRF && j.csrfToken != "" {
		req.Header.Set("X-CSRF-Token", j.csrfToken)
	}
	var pairs []string
	for k, v := range j.cookies {
		pairs = append(pairs, k+"="+v)
	}
	if len(pairs) > 0 {
		req.Header.Set("Cookie", strings.Join(pairs, "; "))
	}

	resp, err := j.env.app.Test(req, -1)
	if err != nil {
		t.Fatalf("%s %s failed: %v", method, path, err)
	}
	for _, sc := range resp.Cookies() {
		if sc.MaxAge < 0 || sc.Value == "" {
			delete(j.cookies, sc.Name)
		} else {
			j.cookies[sc.Name] = sc.Value
		}
	}

	var parsed map[string]any
	_ = json.NewDecoder(resp.Body).Decode(&parsed)
	return resp, parsed
}

func (j *clientJar) initCSRF(t *testing.T) {
	t.Helper()
	resp, data := j.do(t, http.MethodGet, "/api/v1/auth/csrf", nil, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on /auth/csrf, got %d", resp.StatusCode)
	}
	payload := data["data"].(map[string]any)
	j.csrfToken = payload["token"].(string)
}

// --- E2E Security Test Cases (Point 16) ---

func TestSecurityE2E_LoginSessionLogoutAndPasswordReset(t *testing.T) {
	env := newSecurityTestEnv()
	client := newClientJar(env, "198.51.100.10")
	client.initCSRF(t)

	// 1. Weak password rejected
	resp, _ := client.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "weak",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for weak password, got %d", resp.StatusCode)
	}

	// 2. Strong password registration sets HttpOnly + SameSite=Lax cookies and does NOT leak tokens/password in JSON
	resp, body := client.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Alice",
		"email":    "alice@example.com",
		"password": "AliceStrong#2026",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201 on register, got %d (%v)", resp.StatusCode, body)
	}
	rawBody, _ := json.Marshal(body)
	if strings.Contains(string(rawBody), "refresh_token") || strings.Contains(string(rawBody), "AliceStrong#2026") {
		t.Fatalf("auth response must not leak tokens or passwords: %s", string(rawBody))
	}
	if client.cookies["access_token"] == "" || client.cookies["refresh_token"] == "" {
		t.Fatal("expected HttpOnly access_token and refresh_token cookies to be set")
	}

	// 3. Authenticated /auth/me works
	resp, _ = client.do(t, http.MethodGet, "/api/v1/auth/me", nil, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on /auth/me, got %d", resp.StatusCode)
	}

	// 4. Refresh rotates session; old refresh token can no longer be reused
	oldRefresh := client.cookies["refresh_token"]
	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/refresh", nil, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on /auth/refresh, got %d", resp.StatusCode)
	}
	if client.cookies["refresh_token"] == oldRefresh {
		t.Fatal("expected refresh token rotation to issue a new token")
	}

	replayJar := newClientJar(env, "198.51.100.11")
	replayJar.initCSRF(t)
	replayJar.cookies["refresh_token"] = oldRefresh
	resp, _ = replayJar.do(t, http.MethodPost, "/api/v1/auth/refresh", nil, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 when replaying rotated refresh token, got %d", resp.StatusCode)
	}

	// 5. Logout revokes session & clears cookies
	currentRefresh := client.cookies["refresh_token"]
	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/logout", nil, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on logout, got %d", resp.StatusCode)
	}
	if client.cookies["access_token"] != "" || client.cookies["refresh_token"] != "" {
		t.Fatal("expected auth cookies to be cleared on logout")
	}
	client.cookies["refresh_token"] = currentRefresh
	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/refresh", nil, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for revoked session after logout, got %d", resp.StatusCode)
	}

	// 6. Password reset flow updates password and invalidates old password
	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/forgot-password", map[string]string{
		"email": "alice@example.com",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK || len(env.store.resetURLs) != 1 {
		t.Fatalf("expected password reset email to be queued, status=%d urls=%d", resp.StatusCode, len(env.store.resetURLs))
	}
	parsedURL, _ := url.Parse(env.store.resetURLs[0])
	resetToken := parsedURL.Query().Get("token")

	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/reset-password", map[string]string{
		"token":        resetToken,
		"new_password": "AliceNewPass#2026",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on reset-password, got %d", resp.StatusCode)
	}

	// Old password fails, new password succeeds
	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "alice@example.com",
		"password": "AliceStrong#2026",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 with old password after reset, got %d", resp.StatusCode)
	}

	resp, _ = client.do(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "alice@example.com",
		"password": "AliceNewPass#2026",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 with new password after reset, got %d", resp.StatusCode)
	}
}

func TestSecurityE2E_TwoAccountDataIsolationAndAdminRBAC(t *testing.T) {
	env := newSecurityTestEnv()

	// Seed Admin account directly in store
	adminHash, _ := hash.HashPassword("SuperAdmin#2026!")
	_ = env.store.Create(context.Background(), &domain.User{
		Name:       "Admin",
		Email:      "admin@example.com",
		Password:   adminHash,
		IsVerified: true,
		Roles:      []domain.Role{{ID: 1, Name: "Super Admin"}},
	})

	// Register Account A and Account B (2 distinct customer accounts)
	acctA := newClientJar(env, "198.51.100.21")
	acctA.initCSRF(t)
	resp, _ := acctA.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Customer A",
		"email":    "a@example.com",
		"password": "CustomerA#2026!",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed registering Account A: %d", resp.StatusCode)
	}

	acctB := newClientJar(env, "198.51.100.22")
	acctB.initCSRF(t)
	resp, _ = acctB.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Customer B",
		"email":    "b@example.com",
		"password": "CustomerB#2026!",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("failed registering Account B: %d", resp.StatusCode)
	}

	// Account A adds product #1 (sale_price=100,000, weight=500g) qty=2 -> subtotal 200,000 + shipping 15,000 = 215,000
	resp, _ = acctA.do(t, http.MethodPost, "/api/v1/cart/items", map[string]any{
		"product_id": 1,
		"quantity":   2,
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Account A add to cart failed: %d", resp.StatusCode)
	}

	// Account A checks out (frontend only sends delivery info, NEVER price/total)
	resp, orderBody := acctA.do(t, http.MethodPost, "/api/v1/orders", map[string]any{
		"recipient_name":   "Customer A",
		"recipient_phone":  "081234567890",
		"shipping_address": "Jl. Merdeka No. 1, Jakarta",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Account A create order failed: %d (%v)", resp.StatusCode, orderBody)
	}
	orderData := orderBody["data"].(map[string]any)
	orderCode := orderData["order_code"].(string)
	if int64(orderData["total_amount"].(float64)) != 215000 {
		t.Fatalf("expected server-computed total 215000, got %v", orderData["total_amount"])
	}

	// Account A can read own order -> 200
	resp, _ = acctA.do(t, http.MethodGet, "/api/v1/orders/"+orderCode, nil, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected Account A to read own order (200), got %d", resp.StatusCode)
	}

	// Account B CANNOT read Account A's order (IDOR blocked -> 403)
	resp, _ = acctB.do(t, http.MethodGet, "/api/v1/orders/"+orderCode, nil, "", false)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected Account B IDOR order read to be 403, got %d", resp.StatusCode)
	}

	// Account B's order list does NOT contain Account A's order
	resp, listBody := acctB.do(t, http.MethodGet, "/api/v1/orders", nil, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 on Account B order list, got %d", resp.StatusCode)
	}
	if orders, ok := listBody["data"].([]any); ok && len(orders) != 0 {
		t.Fatalf("expected Account B to see 0 orders, got %d", len(orders))
	}

	// Account B CANNOT trigger Snap payment for Account A's order (IDOR blocked -> 403)
	resp, _ = acctB.do(t, http.MethodPost, "/api/v1/payments/"+orderCode+"/snap", nil, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected Account B IDOR payment snap to be 403, got %d", resp.StatusCode)
	}

	// Account B CANNOT access admin routes -> 403
	for _, adminPath := range []string{"/api/v1/admin/orders", "/api/v1/admin/products", "/api/v1/admin/categories"} {
		resp, _ = acctB.do(t, http.MethodGet, adminPath, nil, "", false)
		if resp.StatusCode != http.StatusForbidden {
			t.Fatalf("expected 403 for customer on %s, got %d", adminPath, resp.StatusCode)
		}
	}

	// Admin logs in -> can access admin routes & Account A's order -> 200
	adminJar := newClientJar(env, "198.51.100.23")
	adminJar.initCSRF(t)
	resp, _ = adminJar.do(t, http.MethodPost, "/api/v1/auth/login", map[string]string{
		"email":    "admin@example.com",
		"password": "SuperAdmin#2026!",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("admin login failed: %d", resp.StatusCode)
	}

	resp, _ = adminJar.do(t, http.MethodGet, "/api/v1/admin/orders", nil, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for Super Admin on /admin/orders, got %d", resp.StatusCode)
	}
	resp, _ = adminJar.do(t, http.MethodGet, "/api/v1/orders/"+orderCode, nil, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for Super Admin viewing order, got %d", resp.StatusCode)
	}

	// Revoke Admin role in DB -> immediately blocked (403) even with unexpired JWT cookie
	env.store.mu.Lock()
	env.store.users[1].Roles = []domain.Role{{ID: 2, Name: "Customer"}}
	env.store.mu.Unlock()

	resp, _ = adminJar.do(t, http.MethodGet, "/api/v1/admin/orders", nil, "", false)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 on /admin/orders immediately after role revocation, got %d", resp.StatusCode)
	}
	resp, _ = adminJar.do(t, http.MethodGet, "/api/v1/orders/"+orderCode, nil, "", false)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 on other user's order immediately after role revocation, got %d", resp.StatusCode)
	}
}

func TestSecurityE2E_PaymentWebhookAndCSRFProtection(t *testing.T) {
	env := newSecurityTestEnv()
	acct := newClientJar(env, "198.51.100.30")
	acct.initCSRF(t)

	// 1. CSRF & Origin enforcement on mutating endpoint
	resp, _ := acct.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Buyer",
		"email":    "buyer@example.com",
		"password": "BuyerStrong#2026!",
	}, "http://localhost:3000", false)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 when X-CSRF-Token is missing, got %d", resp.StatusCode)
	}

	resp, _ = acct.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Buyer",
		"email":    "buyer@example.com",
		"password": "BuyerStrong#2026!",
	}, "https://evil.example", true)
	if resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 when Origin is disallowed, got %d", resp.StatusCode)
	}

	// Mutating verify-email must not work via GET
	resp, _ = acct.do(t, http.MethodGet, "/api/v1/auth/verify-email?token=1234567890abcdef", nil, "", false)
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for GET /auth/verify-email, got %d", resp.StatusCode)
	}

	// Register validly and create an order
	resp, _ = acct.do(t, http.MethodPost, "/api/v1/auth/register", map[string]string{
		"name":     "Buyer",
		"email":    "buyer@example.com",
		"password": "BuyerStrong#2026!",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("register failed: %d", resp.StatusCode)
	}
	_, _ = acct.do(t, http.MethodPost, "/api/v1/cart/items", map[string]any{"product_id": 1, "quantity": 1}, "http://localhost:3000", true)
	resp, orderBody := acct.do(t, http.MethodPost, "/api/v1/orders", map[string]any{
		"recipient_name":   "Buyer",
		"recipient_phone":  "081234567890",
		"shipping_address": "Jl. Sudirman No. 9, Jakarta",
	}, "http://localhost:3000", true)
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create order failed: %d", resp.StatusCode)
	}
	orderCode := orderBody["data"].(map[string]any)["order_code"].(string)
	// Total = 100,000 + 15,000 shipping = 115,000

	sign := func(orderID, statusCode, gross string) string {
		sum := sha512.Sum512([]byte(orderID + statusCode + gross + env.midtransSK))
		return hex.EncodeToString(sum[:])
	}

	webhookClient := newClientJar(env, "198.51.100.99")

	// 2. Forged webhook signature -> 401
	resp, _ = webhookClient.do(t, http.MethodPost, "/api/v1/payments/midtrans-callback", map[string]string{
		"order_id":           orderCode,
		"status_code":        "200",
		"gross_amount":       "115000.00",
		"signature_key":      "invalid-signature",
		"transaction_id":     "tx-e2e-1",
		"transaction_status": "settlement",
		"payment_type":       "qris",
		"currency":           "IDR",
	}, "", false)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for invalid Midtrans webhook signature, got %d", resp.StatusCode)
	}

	// 3. Tampered gross_amount -> 400
	resp, _ = webhookClient.do(t, http.MethodPost, "/api/v1/payments/midtrans-callback", map[string]string{
		"order_id":           orderCode,
		"status_code":        "200",
		"gross_amount":       "1000.00",
		"signature_key":      sign(orderCode, "200", "1000.00"),
		"transaction_id":     "tx-e2e-2",
		"transaction_status": "settlement",
		"payment_type":       "qris",
		"currency":           "IDR",
	}, "", false)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 for tampered gross_amount, got %d", resp.StatusCode)
	}

	// 4. Valid settlement webhook -> 200 and order status becomes PAID
	resp, _ = webhookClient.do(t, http.MethodPost, "/api/v1/payments/midtrans-callback", map[string]string{
		"order_id":           orderCode,
		"status_code":        "200",
		"gross_amount":       "115000.00",
		"signature_key":      sign(orderCode, "200", "115000.00"),
		"transaction_id":     "tx-e2e-3",
		"transaction_status": "settlement",
		"payment_type":       "qris",
		"currency":           "IDR",
	}, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for valid settlement webhook, got %d", resp.StatusCode)
	}

	resp, updatedOrder := acct.do(t, http.MethodGet, "/api/v1/orders/"+orderCode, nil, "", false)
	if resp.StatusCode != http.StatusOK || updatedOrder["data"].(map[string]any)["status"] != string(domain.OrderStatusPaid) {
		t.Fatalf("expected order to be PAID after settlement webhook, got %v", updatedOrder)
	}

	// 5. Late expire webhook must NOT cancel or restock already PAID order
	resp, _ = webhookClient.do(t, http.MethodPost, "/api/v1/payments/midtrans-callback", map[string]string{
		"order_id":           orderCode,
		"status_code":        "200",
		"gross_amount":       "115000.00",
		"signature_key":      sign(orderCode, "200", "115000.00"),
		"transaction_id":     "tx-e2e-4",
		"transaction_status": "expire",
		"payment_type":       "qris",
		"currency":           "IDR",
	}, "", false)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 for safely ignored late expire webhook, got %d", resp.StatusCode)
	}
	resp, finalOrder := acct.do(t, http.MethodGet, "/api/v1/orders/"+orderCode, nil, "", false)
	if resp.StatusCode != http.StatusOK || finalOrder["data"].(map[string]any)["status"] != string(domain.OrderStatusPaid) {
		t.Fatalf("expected PAID order to stay PAID after late expire webhook, got %v", finalOrder)
	}
}
