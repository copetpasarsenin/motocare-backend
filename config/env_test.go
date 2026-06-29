package config

import "testing"

func TestEnvInt_DefaultWhenUnset(t *testing.T) {
	t.Setenv("MC_TEST_INT", "")
	if got := envInt("MC_TEST_INT", 42); got != 42 {
		t.Errorf("want default 42, got %d", got)
	}
}

func TestEnvInt_ParsedValue(t *testing.T) {
	t.Setenv("MC_TEST_INT", "17")
	if got := envInt("MC_TEST_INT", 42); got != 17 {
		t.Errorf("want 17, got %d", got)
	}
}

func TestEnvInt_FallsBackOnInvalid(t *testing.T) {
	t.Setenv("MC_TEST_INT", "not-a-number")
	if got := envInt("MC_TEST_INT", 42); got != 42 {
		t.Errorf("want fallback 42 for invalid value, got %d", got)
	}
}

func TestEnvInt_FallsBackOnNonPositive(t *testing.T) {
	t.Setenv("MC_TEST_INT", "0")
	if got := envInt("MC_TEST_INT", 42); got != 42 {
		t.Errorf("want fallback 42 for zero value, got %d", got)
	}
	t.Setenv("MC_TEST_INT", "-5")
	if got := envInt("MC_TEST_INT", 42); got != 42 {
		t.Errorf("want fallback 42 for negative value, got %d", got)
	}
}

func TestEnvDuration_DefaultWhenUnset(t *testing.T) {
	t.Setenv("MC_TEST_DUR", "")
	got := envDuration("MC_TEST_DUR", 5_000_000_000) // 5s in ns
	if got != 5_000_000_000 {
		t.Errorf("want default 5s, got %v", got)
	}
}

func TestEnvDuration_ParsedValue(t *testing.T) {
	t.Setenv("MC_TEST_DUR", "45s")
	if got := envDuration("MC_TEST_DUR", 5_000_000_000); got.String() != "45s" {
		t.Errorf("want 45s, got %v", got)
	}
}

func TestEnvDuration_FallsBackOnInvalid(t *testing.T) {
	t.Setenv("MC_TEST_DUR", "not-a-duration")
	if got := envDuration("MC_TEST_DUR", 5_000_000_000); got != 5_000_000_000 {
		t.Errorf("want fallback for invalid value, got %v", got)
	}
}

func TestEnvDuration_FallsBackOnNonPositive(t *testing.T) {
	t.Setenv("MC_TEST_DUR", "-1h")
	if got := envDuration("MC_TEST_DUR", 5_000_000_000); got != 5_000_000_000 {
		t.Errorf("want fallback for negative value, got %v", got)
	}
}
