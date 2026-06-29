package utils

import "testing"

func TestHashPassword_ProducesBcryptHash(t *testing.T) {
	hashed, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}

	if hashed == "" {
		t.Fatal("expected non-empty hash")
	}
	if hashed == "password123" {
		t.Fatal("hash should not equal plaintext")
	}
	// bcrypt hashes start with $2a$, $2b$, or $2y$
	if hashed[:4] != "$2a$" && hashed[:4] != "$2b$" && hashed[:4] != "$2y$" {
		t.Errorf("unexpected hash prefix, got: %q", hashed[:4])
	}
}

func TestHashPassword_Empty(t *testing.T) {
	hashed, err := HashPassword("")
	if err != nil {
		t.Fatalf("HashPassword empty should not error, got %v", err)
	}
	if hashed == "" {
		t.Fatal("expected non-empty hash even for empty password")
	}
}

func TestCheckPasswordHash_ValidPassword(t *testing.T) {
	hashed, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if !CheckPasswordHash("correct horse battery staple", hashed) {
		t.Fatal("CheckPasswordHash returned false for valid password")
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	hashed, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword returned error: %v", err)
	}
	if CheckPasswordHash("wrong password", hashed) {
		t.Fatal("CheckPasswordHash returned true for incorrect password")
	}
}

func TestCheckPasswordHash_DifferentHashesForSamePassword(t *testing.T) {
	// bcrypt salts randomly, so identical plaintexts should produce different hashes.
	hash1, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("first HashPassword returned error: %v", err)
	}
	hash2, err := HashPassword("same-password")
	if err != nil {
		t.Fatalf("second HashPassword returned error: %v", err)
	}
	if hash1 == hash2 {
		t.Fatal("expected bcrypt to salt hashes uniquely, but two hashes matched")
	}
	// Both should still verify.
	if !CheckPasswordHash("same-password", hash1) || !CheckPasswordHash("same-password", hash2) {
		t.Fatal("expected both hashes to verify the same plaintext")
	}
}

func TestCheckPasswordHash_InvalidHash(t *testing.T) {
	// A malformed hash should return false rather than panic.
	if CheckPasswordHash("any-password", "not-a-real-bcrypt-hash") {
		t.Fatal("CheckPasswordHash returned true for malformed hash")
	}
}
