package midtrans

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	serverKey := "SB-Mid-server-TestKey123"
	client := NewClient(serverKey, "SB-Mid-client-TestKey123", false)

	orderID := "INV-20261123-ABC123"
	statusCode := "200"
	grossAmount := "115000.00"

	sum := sha512.Sum512([]byte(orderID + statusCode + grossAmount + serverKey))
	validSig := hex.EncodeToString(sum[:])

	if !client.VerifySignature(orderID, statusCode, grossAmount, validSig) {
		t.Fatal("expected valid signature to pass")
	}
	if !client.VerifySignature(orderID, statusCode, grossAmount, strings.ToUpper(validSig)) {
		t.Fatal("expected uppercase valid hex signature to pass")
	}
	if client.VerifySignature(orderID, statusCode, "1000.00", validSig) {
		t.Fatal("expected tampered grossAmount to fail verification")
	}
	if client.VerifySignature(orderID, statusCode, grossAmount, "") {
		t.Fatal("expected empty signature to fail verification")
	}

	emptyKeyClient := NewClient("", "", false)
	if emptyKeyClient.VerifySignature(orderID, statusCode, grossAmount, validSig) {
		t.Fatal("expected empty serverKey client to reject signatures")
	}
}
