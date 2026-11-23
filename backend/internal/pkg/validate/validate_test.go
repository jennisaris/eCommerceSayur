package validate

import (
	"ecommerce-backend/internal/domain"
	"strings"
	"testing"
)

func TestSafeURLValidation(t *testing.T) {
	v := New()

	validCases := []string{
		"",
		"/images/products/book-1.jpg",
		"/assets/images/category/cat.svg",
		"https://cdn.example.com/books/cover.png",
		"http://localhost:3000/img.png",
	}
	for _, img := range validCases {
		req := domain.CreateProductRequest{
			CategoryID:   1,
			Name:         "Buku Bagus",
			Slug:         "buku-bagus",
			Image:        img,
			RegularPrice: 50000,
			Weight:       300,
			Stock:        10,
		}
		if err := v.Struct(&req); err != nil {
			t.Errorf("expected image %q to be valid, got error: %v", img, err)
		}
	}

	invalidCases := []string{
		"javascript:alert(1)",
		"JAVASCRIPT:alert(1)",
		"data:text/html,<script>alert(1)</script>",
		"//evil.com/payload.js",
		"ftp://example.com/file.png",
		"/images/bad path.png",
	}
	for _, img := range invalidCases {
		req := domain.CreateProductRequest{
			CategoryID:   1,
			Name:         "Buku Bagus",
			Slug:         "buku-bagus",
			Image:        img,
			RegularPrice: 50000,
			Weight:       300,
			Stock:        10,
		}
		if err := v.Struct(&req); err == nil {
			t.Errorf("expected image %q to be rejected", img)
		}
	}
}

func TestSlugPhoneAndNumericBoundsValidation(t *testing.T) {
	v := New()

	// Invalid negative ProductID and excessive quantity in Cart
	badCart := domain.AddToCartRequest{ProductID: -1, Quantity: 5000}
	if err := v.Struct(&badCart); err == nil {
		t.Fatal("expected negative product_id / excessive quantity to fail validation")
	}

	// Valid checkout payload
	goodCheckout := domain.CheckoutRequest{
		RecipientName:   "Budi Santoso",
		RecipientPhone:  "+62 812-3456-7890",
		ShippingAddress: "Jl. Sudirman No. 10, Jakarta",
		ShippingLat:     "-6.2088",
		ShippingLng:     "106.8456",
		Remarks:         "Tolong dibungkus rapi",
	}
	if err := v.Struct(&goodCheckout); err != nil {
		t.Fatalf("expected valid checkout to pass, got %v", err)
	}

	// Invalid phone & oversized remarks
	badCheckout := domain.CheckoutRequest{
		RecipientName:   "Budi",
		RecipientPhone:  "not-a-phone",
		ShippingAddress: "Jl. Sudirman No. 10",
		Remarks:         strings.Repeat("a", 300),
	}
	err := v.Struct(&badCheckout)
	if err == nil {
		t.Fatal("expected bad checkout to fail validation")
	}
	msg := FormatError(err)
	if strings.Contains(msg, "CheckoutRequest") {
		t.Fatalf("formatted validation error must not leak internal Go struct names: %s", msg)
	}

	// Invalid slug with uppercase/spaces
	badCat := domain.CreateCategoryRequest{
		Name: "Fiksi",
		Slug: "Fiksi Baru!",
	}
	if err := v.Struct(&badCat); err == nil {
		t.Fatal("expected invalid slug to fail validation")
	}
}
