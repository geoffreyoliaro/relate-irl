package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type xanoAuthResponse struct {
	AuthToken string `json:"authToken"`
	User      struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"user"`
}

// Login validates credentials against Xano, then issues our own short-lived JWT.
// This pattern decouples our graph API from Xano's session model.
func Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Step 1: Validate with Xano
	xanoUser, err := validateWithXano(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Step 2: Issue our own JWT
	token, err := issueJWT(fmt.Sprintf("%d", xanoUser.User.ID), xanoUser.User.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not issue token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_in": 3600,
		"user": gin.H{
			"id":    xanoUser.User.ID,
			"name":  xanoUser.User.Name,
			"email": xanoUser.User.Email,
		},
	})
}

func validateWithXano(email, password string) (*xanoAuthResponse, error) {
	xanoEndpoint := os.Getenv("XANO_API_BASE") + "/auth/login"

	body, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
	})

	resp, err := http.Post(xanoEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("xano returned %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var xanoResp xanoAuthResponse
	if err := json.Unmarshal(data, &xanoResp); err != nil {
		return nil, err
	}
	return &xanoResp, nil
}

func issueJWT(userID, email string) (string, error) {
	secret := []byte(os.Getenv("JWT_SECRET"))
	claims := jwt.MapClaims{
		"sub":   userID,
		"email": email,
		"iss":   "mininexus",
		"exp":   time.Now().Add(time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
