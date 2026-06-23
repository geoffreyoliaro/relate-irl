package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func TestLogin_ValidCredentials(t *testing.T) {
	// Set env vars for testing
	os.Setenv("JWT_SECRET", "test-secret-key-for-unit-tests")

	// Create a mock Xano server
	mockXano := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/auth/login" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		var req map[string]string
		json.NewDecoder(r.Body).Decode(&req)

		// Accept test credentials
		if req["email"] == "test@example.com" && req["password"] == "testpass" {
			resp := map[string]interface{}{
				"authToken": "mock-token",
				"user": map[string]interface{}{
					"id":    123,
					"name":  "Test User",
					"email": "test@example.com",
				},
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}

		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer mockXano.Close()

	os.Setenv("XANO_API_BASE", mockXano.URL)

	// Create request
	loginReq := loginRequest{
		Email:    "test@example.com",
		Password: "testpass",
	}
	body, _ := json.Marshal(loginReq)

	// Test endpoint
	router := gin.Default()
	router.POST("/auth/login", Login)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(w.Body).Decode(&resp)

	if _, ok := resp["token"]; !ok {
		t.Error("Expected 'token' in response")
	}

	if resp["expires_in"] != float64(3600) {
		t.Errorf("Expected expires_in 3600, got %v", resp["expires_in"])
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")

	mockXano := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}))
	defer mockXano.Close()

	os.Setenv("XANO_API_BASE", mockXano.URL)

	loginReq := loginRequest{
		Email:    "wrong@example.com",
		Password: "wrongpass",
	}
	body, _ := json.Marshal(loginReq)

	router := gin.Default()
	router.POST("/auth/login", Login)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401, got %d", w.Code)
	}
}

func TestIssueJWT_ValidToken(t *testing.T) {
	secret := "test-secret"
	os.Setenv("JWT_SECRET", secret)

	token, err := issueJWT("user123", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to issue JWT: %v", err)
	}

	// Verify token structure
	if token == "" {
		t.Error("JWT token is empty")
	}

	// Parse and validate token
	parsed, err := jwt.Parse(token, func(tok *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		t.Fatalf("Failed to parse token: %v", err)
	}

	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to extract claims")
	}

	if claims["sub"] != "user123" {
		t.Errorf("Expected sub 'user123', got %v", claims["sub"])
	}

	if claims["email"] != "user@example.com" {
		t.Errorf("Expected email 'user@example.com', got %v", claims["email"])
	}

	if claims["iss"] != "mininexus" {
		t.Errorf("Expected iss 'mininexus', got %v", claims["iss"])
	}

	// Verify expiration is ~1 hour in future
	exp := time.Unix(int64(claims["exp"].(float64)), 0)
	now := time.Now()
	diff := exp.Sub(now)
	if diff < 50*time.Minute || diff > 70*time.Minute {
		t.Errorf("Token expiration time is unexpected: %v", diff)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	router := gin.Default()
	router.POST("/auth/login", Login)

	// Missing password
	loginReq := loginRequest{Email: "test@example.com"}
	body, _ := json.Marshal(loginReq)

	req := httptest.NewRequest("POST", "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}
