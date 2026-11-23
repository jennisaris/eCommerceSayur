package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/pkg/hash"
	"ecommerce-backend/internal/pkg/jwt"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"net/url"
	"strings"
	"time"
	"unicode"
)

var dummyPasswordHash string

func init() {
	if hashed, err := hash.HashPassword("BookVerse#TimingMitigation2026!"); err == nil {
		dummyPasswordHash = hashed
	}
}

type AuthUseCase struct {
	userRepo                 domain.UserRepository
	cartRepo                 domain.CartRepository
	jwtSvc                   *jwt.JWTService
	refreshTokenTTL          time.Duration
	requireEmailVerification bool
	verificationEmailSender  VerificationEmailSender
	passwordResetEmailSender PasswordResetEmailSender
	verificationBaseURL      string
	isProduction             bool
}

type VerificationEmailSender interface {
	SendVerificationEmail(ctx context.Context, recipient, recipientName, verificationURL string) error
}

type PasswordResetEmailSender interface {
	IsConfigured() bool
	SendPasswordResetEmail(ctx context.Context, recipient, recipientName, resetURL string) error
}

type AuthOptions struct {
	RefreshTokenTTL          time.Duration
	RequireEmailVerification bool
	VerificationEmailSender  VerificationEmailSender
	PasswordResetEmailSender PasswordResetEmailSender
	VerificationBaseURL      string
	IsProduction             bool
}

func NewAuthUseCase(userRepo domain.UserRepository, cartRepo domain.CartRepository, jwtSvc *jwt.JWTService, refreshTokenTTL ...time.Duration) domain.AuthUseCase {
	ttl := 30 * 24 * time.Hour
	if len(refreshTokenTTL) > 0 && refreshTokenTTL[0] > 0 {
		ttl = refreshTokenTTL[0]
	}
	return newAuthUseCase(userRepo, cartRepo, jwtSvc, AuthOptions{RefreshTokenTTL: ttl})
}

func NewAuthUseCaseWithOptions(userRepo domain.UserRepository, cartRepo domain.CartRepository, jwtSvc *jwt.JWTService, options AuthOptions) domain.AuthUseCase {
	return newAuthUseCase(userRepo, cartRepo, jwtSvc, options)
}

func newAuthUseCase(userRepo domain.UserRepository, cartRepo domain.CartRepository, jwtSvc *jwt.JWTService, options AuthOptions) domain.AuthUseCase {
	ttl := options.RefreshTokenTTL
	if ttl <= 0 {
		ttl = 30 * 24 * time.Hour
	}
	resetSender := options.PasswordResetEmailSender
	if resetSender == nil {
		if candidate, ok := options.VerificationEmailSender.(PasswordResetEmailSender); ok {
			resetSender = candidate
		}
	}
	return &AuthUseCase{
		userRepo:                 userRepo,
		cartRepo:                 cartRepo,
		jwtSvc:                   jwtSvc,
		refreshTokenTTL:          ttl,
		requireEmailVerification: options.RequireEmailVerification,
		verificationEmailSender:  options.VerificationEmailSender,
		passwordResetEmailSender: resetSender,
		verificationBaseURL:      options.VerificationBaseURL,
		isProduction:             options.IsProduction,
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func newRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

func hashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func hashVerificationToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func hasStrongPassword(password string) bool {
	characters := []rune(password)
	if len(characters) < 12 || len(password) > 72 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, char := range characters {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		default:
			hasSpecial = true
		}
	}
	return hasUpper && hasLower && hasDigit && hasSpecial
}

func (u *AuthUseCase) createAuthResponse(ctx context.Context, user *domain.User, roleNames []string) (*domain.AuthResponse, error) {
	accessToken, err := u.jwtSvc.GenerateToken(user.ID, user.Email, roleNames)
	if err != nil {
		return nil, err
	}
	refreshToken, err := newRefreshToken()
	if err != nil {
		return nil, err
	}
	if err := u.userRepo.CreateAuthSession(ctx, &domain.AuthSession{
		UserID:    user.ID,
		TokenHash: hashRefreshToken(refreshToken),
		ExpiresAt: time.Now().Add(u.refreshTokenTTL),
	}); err != nil {
		return nil, err
	}
	return &domain.AuthResponse{Token: accessToken, RefreshToken: refreshToken, User: user}, nil
}

func (u *AuthUseCase) Register(ctx context.Context, req *domain.RegisterRequest) (*domain.AuthResponse, error) {
	if !hasStrongPassword(req.Password) {
		return nil, domain.ErrWeakPassword
	}

	normalizedEmail := normalizeEmail(req.Email)

	// Check if email already taken
	existing, err := u.userRepo.GetByEmail(ctx, normalizedEmail)
	if err == nil && existing != nil {
		return nil, domain.ErrEmailAlreadyExists
	}

	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		Name:       strings.TrimSpace(req.Name),
		Email:      normalizedEmail,
		Password:   hashedPassword,
		Phone:      req.Phone,
		Address:    req.Address,
		IsVerified: !u.requireEmailVerification,
	}

	if err := u.userRepo.Create(ctx, newUser); err != nil {
		return nil, err
	}

	// Assign 'Customer' role (role_id = 2)
	_ = u.userRepo.AssignRole(ctx, newUser.ID, 2)

	// Pre-create empty cart for user
	_ = u.cartRepo.Create(ctx, &domain.Cart{UserID: newUser.ID})

	// Fetch fresh user with roles
	userWithRoles, err := u.userRepo.GetByID(ctx, newUser.ID)
	if err != nil {
		userWithRoles = newUser
	}

	var roleNames []string
	for _, r := range userWithRoles.Roles {
		roleNames = append(roleNames, r.Name)
	}
	if len(roleNames) == 0 {
		roleNames = []string{"Customer"}
	}

	if u.requireEmailVerification {
		if u.verificationEmailSender == nil || u.verificationBaseURL == "" {
			return nil, errors.New("email verification is not configured")
		}
		verificationToken, err := newRefreshToken()
		if err != nil {
			return nil, err
		}
		if err := u.userRepo.CreateVerificationToken(ctx, &domain.VerificationToken{
			UserID:    newUser.ID,
			Token:     hashVerificationToken(verificationToken),
			TokenType: "EMAIL_VERIFICATION",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}); err != nil {
			return nil, err
		}
		verificationURL := u.verificationBaseURL + "/verify-email?token=" + verificationToken
		if err := u.verificationEmailSender.SendVerificationEmail(ctx, newUser.Email, newUser.Name, verificationURL); err != nil {
			return nil, err
		}
		return &domain.AuthResponse{User: userWithRoles, VerificationRequired: true}, nil
	}

	return u.createAuthResponse(ctx, userWithRoles, roleNames)
}

func (u *AuthUseCase) Login(ctx context.Context, req *domain.LoginRequest) (*domain.AuthResponse, error) {
	user, err := u.userRepo.GetByEmail(ctx, normalizeEmail(req.Email))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			if dummyPasswordHash != "" {
				_ = hash.CheckPassword(req.Password, dummyPasswordHash)
			}
			return nil, domain.ErrInvalidCredentials
		}
		return nil, err
	}

	if !hash.CheckPassword(req.Password, user.Password) {
		return nil, domain.ErrInvalidCredentials
	}
	if u.requireEmailVerification && !user.IsVerified {
		return nil, domain.ErrUserNotVerified
	}

	var roleNames []string
	for _, r := range user.Roles {
		roleNames = append(roleNames, r.Name)
	}

	return u.createAuthResponse(ctx, user, roleNames)
}

func (u *AuthUseCase) Refresh(ctx context.Context, refreshToken string) (*domain.AuthResponse, error) {
	if refreshToken == "" {
		return nil, domain.ErrInvalidSession
	}

	tokenHash := hashRefreshToken(refreshToken)
	session, err := u.userRepo.GetActiveAuthSession(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	user, err := u.userRepo.GetByID(ctx, session.UserID)
	if err != nil {
		return nil, domain.ErrInvalidSession
	}
	if err := u.userRepo.RevokeAuthSession(ctx, tokenHash); err != nil {
		return nil, err
	}

	roleNames := make([]string, 0, len(user.Roles))
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}
	return u.createAuthResponse(ctx, user, roleNames)
}

func (u *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	return u.userRepo.RevokeAuthSession(ctx, hashRefreshToken(refreshToken))
}

func (u *AuthUseCase) GetProfile(ctx context.Context, userID int64) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, userID)
}

func (u *AuthUseCase) UpdateProfile(ctx context.Context, userID int64, req *domain.UpdateProfileRequest) (*domain.User, bool, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, false, err
	}

	newEmail := normalizeEmail(req.Email)
	emailChanged := newEmail != "" && newEmail != normalizeEmail(user.Email)

	if emailChanged {
		if req.CurrentPassword == "" || !hash.CheckPassword(req.CurrentPassword, user.Password) {
			return nil, false, domain.ErrInvalidCredentials
		}
		existing, lookupErr := u.userRepo.GetByEmail(ctx, newEmail)
		if lookupErr == nil && existing != nil && existing.ID != user.ID {
			return nil, false, domain.ErrEmailAlreadyExists
		}
		user.Email = newEmail
	}

	user.Name = strings.TrimSpace(req.Name)
	user.Phone = strings.TrimSpace(req.Phone)
	user.Address = strings.TrimSpace(req.Address)

	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, false, err
	}

	updated, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return user, emailChanged, nil
	}
	return updated, emailChanged, nil
}

func (u *AuthUseCase) VerifyEmail(ctx context.Context, tokenStr string) error {
	token, err := u.userRepo.GetVerificationToken(ctx, tokenStr, "EMAIL_VERIFICATION")
	if err != nil {
		return err
	}

	if token.ExpiresAt.Before(time.Now()) {
		return errors.New("verification token expired")
	}

	user, err := u.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return err
	}

	user.IsVerified = true
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	_ = u.userRepo.DeleteVerificationToken(ctx, token.ID)
	return nil
}

func (u *AuthUseCase) RequestPasswordReset(ctx context.Context, email string) error {
	normalized := normalizeEmail(email)
	if normalized == "" {
		return nil
	}

	user, err := u.userRepo.GetByEmail(ctx, normalized)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}

	resetToken, err := newRefreshToken()
	if err != nil {
		return err
	}

	_ = u.userRepo.DeleteVerificationTokensByUser(ctx, user.ID, "RESET_PASSWORD")
	if err := u.userRepo.CreateVerificationToken(ctx, &domain.VerificationToken{
		UserID:    user.ID,
		Token:     hashVerificationToken(resetToken),
		TokenType: "RESET_PASSWORD",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}); err != nil {
		return err
	}

	baseURL := strings.TrimRight(strings.TrimSpace(u.verificationBaseURL), "/")
	if baseURL == "" {
		baseURL = "http://localhost:3000"
	}
	resetURL := baseURL + "/reset-password?token=" + url.QueryEscape(resetToken)

	if u.passwordResetEmailSender != nil && u.passwordResetEmailSender.IsConfigured() {
		return u.passwordResetEmailSender.SendPasswordResetEmail(ctx, user.Email, user.Name, resetURL)
	}
	if u.isProduction {
		return errors.New("password reset email sender is not configured")
	}

	log.Printf("[DEV PASSWORD RESET] email=%s reset_url=%s", user.Email, resetURL)
	return nil
}

func (u *AuthUseCase) ResetPassword(ctx context.Context, tokenStr string, newPassword string) error {
	if !hasStrongPassword(newPassword) {
		return domain.ErrWeakPassword
	}

	trimmedToken := strings.TrimSpace(tokenStr)
	if trimmedToken == "" {
		return domain.ErrInvalidResetToken
	}

	token, err := u.userRepo.GetVerificationToken(ctx, trimmedToken, "RESET_PASSWORD")
	if err != nil {
		return domain.ErrInvalidResetToken
	}
	if token.ExpiresAt.Before(time.Now()) {
		_ = u.userRepo.DeleteVerificationToken(ctx, token.ID)
		return domain.ErrInvalidResetToken
	}

	user, err := u.userRepo.GetByID(ctx, token.UserID)
	if err != nil {
		return domain.ErrInvalidResetToken
	}

	hashedPassword, err := hash.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	if err := u.userRepo.Update(ctx, user); err != nil {
		return err
	}

	_ = u.userRepo.DeleteVerificationTokensByUser(ctx, user.ID, "RESET_PASSWORD")
	if err := u.userRepo.RevokeAllUserSessions(ctx, user.ID); err != nil {
		return err
	}
	return nil
}
