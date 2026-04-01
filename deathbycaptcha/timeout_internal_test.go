package deathbycaptcha

import "testing"

func TestEffectiveTimeout_UsesTypeWhenProvided(t *testing.T) {
	t.Parallel()

	if got := effectiveTimeout(0, true, "0"); got != DefaultTimeout {
		t.Fatalf("type=0 should use DefaultTimeout=%d, got %d", DefaultTimeout, got)
	}
	if got := effectiveTimeout(0, false, "4"); got != DefaultTokenTimeout {
		t.Fatalf("type=4 should use DefaultTokenTimeout=%d, got %d", DefaultTokenTimeout, got)
	}
	if got := effectiveTimeout(0, false, "13"); got != DefaultTokenTimeout {
		t.Fatalf("type=13 should use DefaultTokenTimeout=%d, got %d", DefaultTokenTimeout, got)
	}
}

func TestEffectiveTimeout_FallbackBehavior(t *testing.T) {
	t.Parallel()

	if got := effectiveTimeout(0, false, ""); got != DefaultTimeout {
		t.Fatalf("image fallback should use DefaultTimeout=%d, got %d", DefaultTimeout, got)
	}
	if got := effectiveTimeout(0, true, ""); got != DefaultTokenTimeout {
		t.Fatalf("token fallback should use DefaultTokenTimeout=%d, got %d", DefaultTokenTimeout, got)
	}
	if got := effectiveTimeout(42, true, "0"); got != 42 {
		t.Fatalf("explicit timeout should win, expected 42 got %d", got)
	}
}
