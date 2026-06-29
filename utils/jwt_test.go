package utils

import (
	"testing"
	"time"

	"motocare-dashboard/backend/models"

	"github.com/golang-jwt/jwt/v5"
)

func TestParseToken_MissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := ParseToken("any.token.here")
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is unset, got nil")
	}
}

func TestGetJWTSecretMissingYieldsError(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := getJWTSecret()
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is unset, got nil")
	}
}

func TestGenerateToken_MissingSecret(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	_, err := GenerateToken(models.User{ID: 7, Role: "user"})
	if err == nil {
		t.Fatal("expected error when JWT_SECRET is unset, got nil")
	}
}

func TestTokenRoundTrip(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-secret-key-for-roundtrip")

	user := models.User{
		ID:       42,
		Username: "tester",
		Email:    "tester@motocare.test",
		Role:     "user",
	}
	tokenString, err := GenerateToken(user)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if tokenString == "" {
		t.Fatal("expected non-empty token string")
	}

	claims, err := ParseToken(tokenString)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}

	if claims.UserID != user.ID {
		t.Errorf("UserID: want %d, got %d", user.ID, claims.UserID)
	}
	if claims.Username != user.Username {
		t.Errorf("Username: want %q, got %q", user.Username, claims.Username)
	}
	if claims.Email != user.Email {
		t.Errorf("Email: want %q, got %q", user.Email, claims.Email)
	}
	if claims.Role != user.Role {
		t.Errorf("Role: want %q, got %q", user.Role, claims.Role)
	}
	if claims.ExpiresAt == nil {
		t.Error("expected ExpiresAt to be set")
	}
	// Sanity: expiry should be ~24h from now (matches GenerateToken config).
	if time.Until(claims.ExpiresAt.Time) > 25*time.Hour || time.Until(claims.ExpiresAt.Time) < 23*time.Hour {
		t.Errorf("ExpiresAt outside expected 24h window, got %v", claims.ExpiresAt.Time)
	}
}

func TestParseToken_InvalidSigningMethod(t *testing.T) {
	t.Setenv("JWT_SECRET", "")

	claims := jwt.MapClaims{
		"user_id": 1,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	// Sign with HMAC SHA-512 to test the signing-method whitelist.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte("doesnt-matter"))
	if err != nil {
		t.Fatalf("could not sign test token: %v", err)
	}

	_, err = ParseToken(tokenString)
	if err == nil {
		t.Fatal("expected error for non-HS256 signing method, got nil")
	}
	if err.Error() != "invalid token signing method" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestParseToken_MalformedToken(t *testing.T) {
	t.Setenv("JWT_SECRET", "anything")

	if _, err := ParseToken("not-a-real-jwt"); err == nil {
		t.Fatal("expected error for malformed token, got nil")
	}
}
