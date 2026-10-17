package jwt_test

import (
	"ecommerce-backend/internal/pkg/jwt"
	"testing"
	"time"
)

func TestJWTService_GenerateAndValidateToken(t *testing.T) {
	secret := "test-secret-key-that-is-at-least-32-chars-long"
	svc := jwt.NewJWTService(secret, 2)

	userID := int64(42)
	email := "test@example.com"
	roles := []string{"Customer", "Super Admin"}

	// 1. Generate token
	tokenStr, err := svc.GenerateToken(userID, email, roles)
	if err != nil {
		t.Fatalf("expected no error generating token, got: %v", err)
	}
	if tokenStr == "" {
		t.Fatal("expected non-empty token string")
	}

	// 2. Validate token
	claims, err := svc.ValidateToken(tokenStr)
	if err != nil {
		t.Fatalf("expected valid token, got: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected email %s, got %s", email, claims.Email)
	}
	if len(claims.Roles) != len(roles) {
		t.Errorf("expected %d roles, got %d", len(roles), len(claims.Roles))
	}
}

func TestJWTService_InvalidSignature(t *testing.T) {
	svc1 := jwt.NewJWTService("secret-key-alpha-32-chars-length-1234", 1)
	svc2 := jwt.NewJWTService("secret-key-bravo-32-chars-length-5678", 1)

	tokenStr, err := svc1.GenerateToken(1, "user@example.com", []string{"Customer"})
	if err != nil {
		t.Fatalf("failed generating token: %v", err)
	}

	// Validating with a different secret should fail
	_, err = svc2.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected validation error due to mismatched secret, but got nil")
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	// Expired immediately (-1 hour)
	svc := jwt.NewJWTService("secret-key-charlie-32-chars-length-9999", -1)

	tokenStr, err := svc.GenerateToken(1, "expired@example.com", []string{"Customer"})
	if err != nil {
		t.Fatalf("failed generating token: %v", err)
	}

	time.Sleep(10 * time.Millisecond)

	_, err = svc.ValidateToken(tokenStr)
	if err == nil {
		t.Fatal("expected expired token error, got nil")
	}
}
