package hash_test

import (
	"ecommerce-backend/internal/pkg/hash"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestHashPassword_And_CheckPassword(t *testing.T) {
	plainPassword := "SecretPassword123!"

	hashed, err := hash.HashPassword(plainPassword)
	if err != nil {
		t.Fatalf("failed hashing password: %v", err)
	}

	if hashed == plainPassword {
		t.Fatal("hashed password must not match plain text")
	}

	cost, err := bcrypt.Cost([]byte(hashed))
	if err != nil {
		t.Fatalf("failed reading bcrypt cost: %v", err)
	}
	if cost != hash.BcryptCost {
		t.Fatalf("expected bcrypt cost %d, got %d", hash.BcryptCost, cost)
	}

	// Correct password check
	if !hash.CheckPassword(plainPassword, hashed) {
		t.Fatal("expected CheckPassword to return true for correct password")
	}

	// Incorrect password check
	if hash.CheckPassword("WrongPassword123!", hashed) {
		t.Fatal("expected CheckPassword to return false for incorrect password")
	}
}
