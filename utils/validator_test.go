package utils

import "testing"

type testStruct struct {
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required"`
	Age      int    `json:"age" validate:"gte=0"`
}

func TestValidateStruct_AllValid(t *testing.T) {
	payload := testStruct{
		Email:    "user@example.com",
		Username: "budi",
		Age:      25,
	}

	errs := ValidateStruct(payload)
	if errs != nil {
		t.Errorf("expected no errors for valid struct, got: %v", errs)
	}
}

func TestValidateStruct_AllInvalid(t *testing.T) {
	payload := testStruct{
		// All fields empty / invalid.
		Email:    "not-an-email",
		Username: "",
		Age:      -1,
	}

	errs := ValidateStruct(payload)
	if errs == nil {
		t.Fatal("expected errors for invalid struct, got nil")
	}
	// Errors should be keyed by the JSON tag name (RegisterTagNameFunc maps to JSON).
	if _, ok := errs["email"]; !ok {
		t.Errorf("expected 'email' key in errors, got: %v", errs)
	}
	if _, ok := errs["username"]; !ok {
		t.Errorf("expected 'username' key in errors, got: %v", errs)
	}
	if _, ok := errs["age"]; !ok {
		t.Errorf("expected 'age' key in errors, got: %v", errs)
	}
}

func TestValidateStruct_PartialInvalid(t *testing.T) {
	payload := testStruct{
		Email:    "valid@example.com",
		Username: "budi",
		Age:      -5, // invalid (gte=0)
	}

	errs := ValidateStruct(payload)
	if errs == nil {
		t.Fatal("expected errors for invalid age, got nil")
	}
	if _, ok := errs["age"]; !ok {
		t.Errorf("expected 'age' key, got: %v", errs)
	}
	if _, ok := errs["email"]; ok {
		t.Errorf("did not expect 'email' key for valid email, got: %v", errs)
	}
}
