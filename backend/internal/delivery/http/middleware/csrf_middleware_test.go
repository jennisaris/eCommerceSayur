package middleware

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestCSRFProtectionRequiresMatchingHeaderAndValidOrigin(t *testing.T) {
	app := fiber.New()
	app.Use(CSRFProtection(false, "http://localhost:3000,https://bookverse.example"))
	app.Post("/unsafe", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusNoContent) })

	firstResponse, err := app.Test(httptest.NewRequest("POST", "/unsafe", nil))
	if err != nil {
		t.Fatalf("first request failed: %v", err)
	}
	if firstResponse.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected missing CSRF header to be forbidden, got %d", firstResponse.StatusCode)
	}
	csrfCookie := firstResponse.Header.Get("Set-Cookie")
	if csrfCookie == "" {
		t.Fatal("expected CSRF cookie to be issued")
	}
	lowerCookie := strings.ToLower(csrfCookie)
	if !strings.Contains(lowerCookie, "httponly") || !strings.Contains(lowerCookie, "samesite=lax") {
		t.Fatalf("expected CSRF cookie to set HttpOnly and SameSite=Lax, got %q", csrfCookie)
	}

	csrfValue := csrfCookieValue(csrfCookie)

	// 1. Valid token + allowed Origin -> 204
	validRequest := httptest.NewRequest("POST", "/unsafe", nil)
	validRequest.Header.Set("Origin", "https://bookverse.example")
	validRequest.Header.Set("Cookie", csrfCookieName+"="+csrfValue)
	validRequest.Header.Set("X-CSRF-Token", csrfValue)
	validResponse, err := app.Test(validRequest)
	if err != nil {
		t.Fatalf("valid request failed: %v", err)
	}
	if validResponse.StatusCode != fiber.StatusNoContent {
		t.Fatalf("expected matching CSRF header and allowed origin to pass, got %d", validResponse.StatusCode)
	}

	// 2. Valid token + disallowed Origin -> 403
	badOriginReq := httptest.NewRequest("POST", "/unsafe", nil)
	badOriginReq.Header.Set("Origin", "https://evil.example")
	badOriginReq.Header.Set("Cookie", csrfCookieName+"="+csrfValue)
	badOriginReq.Header.Set("X-CSRF-Token", csrfValue)
	badOriginResp, err := app.Test(badOriginReq)
	if err != nil {
		t.Fatalf("bad origin request failed: %v", err)
	}
	if badOriginResp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected disallowed Origin to be blocked with 403, got %d", badOriginResp.StatusCode)
	}

	// 3. Valid token + disallowed Referer (when Origin omitted) -> 403
	badRefReq := httptest.NewRequest("POST", "/unsafe", nil)
	badRefReq.Header.Set("Referer", "https://evil.example/phish")
	badRefReq.Header.Set("Cookie", csrfCookieName+"="+csrfValue)
	badRefReq.Header.Set("X-CSRF-Token", csrfValue)
	badRefResp, err := app.Test(badRefReq)
	if err != nil {
		t.Fatalf("bad referer request failed: %v", err)
	}
	if badRefResp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected disallowed Referer to be blocked with 403, got %d", badRefResp.StatusCode)
	}

	// 4. Valid token + Sec-Fetch-Site: cross-site -> 403
	crossSiteReq := httptest.NewRequest("POST", "/unsafe", nil)
	crossSiteReq.Header.Set("Sec-Fetch-Site", "cross-site")
	crossSiteReq.Header.Set("Cookie", csrfCookieName+"="+csrfValue)
	crossSiteReq.Header.Set("X-CSRF-Token", csrfValue)
	crossSiteResp, err := app.Test(crossSiteReq)
	if err != nil {
		t.Fatalf("cross-site request failed: %v", err)
	}
	if crossSiteResp.StatusCode != fiber.StatusForbidden {
		t.Fatalf("expected Sec-Fetch-Site: cross-site to be blocked with 403, got %d", crossSiteResp.StatusCode)
	}
}

func csrfCookieValue(setCookie string) string {
	const prefix = csrfCookieName + "="
	start := len(prefix)
	for index := start; index < len(setCookie); index++ {
		if setCookie[index] == ';' {
			return setCookie[start:index]
		}
	}
	return setCookie[start:]
}
