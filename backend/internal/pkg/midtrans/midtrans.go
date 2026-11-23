package midtrans

import (
	"bytes"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var ErrTransactionNotFound = errors.New("midtrans transaction not found")

type Client struct {
	serverKey     string
	clientKey     string
	isProduction  bool
	apiBaseURL    string
	httpClient    *http.Client
}

func NewClient(serverKey, clientKey string, isProduction bool) *Client {
	return &Client{
		serverKey:    serverKey,
		clientKey:    clientKey,
		isProduction: isProduction,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// SetAPIBaseURL overrides the Midtrans Core API base URL (used in tests).
func (c *Client) SetAPIBaseURL(baseURL string) {
	c.apiBaseURL = strings.TrimRight(baseURL, "/")
}

func (c *Client) getBaseURL() string {
	if c.isProduction {
		return "https://app.midtrans.com"
	}
	return "https://app.sandbox.midtrans.com"
}

func (c *Client) getAPIBaseURL() string {
	if c.apiBaseURL != "" {
		return c.apiBaseURL
	}
	if c.isProduction {
		return "https://api.midtrans.com"
	}
	return "https://api.sandbox.midtrans.com"
}

type SnapTransactionDetails struct {
	OrderID     string `json:"order_id"`
	GrossAmount int64  `json:"gross_amount"`
}

type SnapCustomerDetails struct {
	FirstName       string               `json:"first_name"`
	Email           string               `json:"email"`
	Phone           string               `json:"phone"`
	ShippingAddress *SnapShippingAddress `json:"shipping_address,omitempty"`
}

type SnapShippingAddress struct {
	FirstName   string `json:"first_name"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	City        string `json:"city,omitempty"`
	PostalCode  string `json:"postal_code,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

type SnapItemDetails struct {
	ID       string `json:"id"`
	Price    int64  `json:"price"`
	Quantity int    `json:"quantity"`
	Name     string `json:"name"`
}

type CreateSnapRequest struct {
	TransactionDetails SnapTransactionDetails `json:"transaction_details"`
	CustomerDetails    SnapCustomerDetails    `json:"customer_details"`
	ItemDetails        []SnapItemDetails      `json:"item_details,omitempty"`
}

type CreateSnapResponse struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages,omitempty"`
}

func (c *Client) CreateSnapTransaction(req *CreateSnapRequest) (*CreateSnapResponse, error) {
	url := fmt.Sprintf("%s/snap/v1/transactions", c.getBaseURL())

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snap request: %w", err)
	}

	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// Basic auth: base64(serverKey + ":")
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.serverKey+":"))
	httpReq.Header.Set("Authorization", authHeader)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans snap request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read midtrans response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("midtrans error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var snapResp CreateSnapResponse
	if err := json.Unmarshal(respBody, &snapResp); err != nil {
		return nil, fmt.Errorf("failed to decode snap response: %w", err)
	}

	if snapResp.Token == "" {
		return nil, errors.New("midtrans returned empty snap token")
	}

	return &snapResp, nil
}

// GetTransactionStatus queries Midtrans Core API v2/{order_id}/status server-to-server using ServerKey.
func (c *Client) GetTransactionStatus(orderID string) ([]byte, error) {
	if strings.TrimSpace(c.serverKey) == "" {
		return nil, errors.New("midtrans server key is not configured")
	}
	endpoint := fmt.Sprintf("%s/v2/%s/status", c.getAPIBaseURL(), url.PathEscape(orderID))
	httpReq, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create status request: %w", err)
	}
	authHeader := "Basic " + base64.StdEncoding.EncodeToString([]byte(c.serverKey+":"))
	httpReq.Header.Set("Authorization", authHeader)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("midtrans status request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read midtrans status response: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrTransactionNotFound
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("midtrans status error (status %d): %s", resp.StatusCode, string(respBody))
	}

	// Midtrans sometimes returns HTTP 200 with {"status_code":"404",...} when transaction does not exist yet
	var probe struct {
		StatusCode string `json:"status_code"`
	}
	if json.Unmarshal(respBody, &probe) == nil && probe.StatusCode == "404" {
		return nil, ErrTransactionNotFound
	}

	return respBody, nil
}

// VerifySignature validates SHA-512 signature from Midtrans callback using constant-time comparison:
// SHA512(order_id + status_code + gross_amount + ServerKey)
func (c *Client) VerifySignature(orderID, statusCode, grossAmount, incomingSignature string) bool {
	if strings.TrimSpace(c.serverKey) == "" || strings.TrimSpace(incomingSignature) == "" {
		return false
	}
	raw := orderID + statusCode + grossAmount + c.serverKey
	hash := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(hash[:])
	candidate := strings.ToLower(strings.TrimSpace(incomingSignature))
	if len(expected) != len(candidate) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(candidate)) == 1
}
