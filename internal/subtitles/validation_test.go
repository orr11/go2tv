package subtitles

import (
	"testing"
)

func TestValidLanguage(t *testing.T) {
	tt := []struct {
		lang     string
		expected bool
	}{
		{"en", true},
		{"es", true},
		{"fr", true},
		{"de", true},
		{"zh", true},
		{"xx", false},
		{"all", true},
		{"", true},
	}

	for _, tc := range tt {
		if got := ValidLanguage(tc.lang); got != tc.expected {
			t.Errorf("ValidLanguage(%q) = %v, want %v", tc.lang, got, tc.expected)
		}
	}
}

func TestValidType(t *testing.T) {
	tt := []struct {
		typ      string
		expected bool
	}{
		{"srt", true},
		{"vtt", true},
		{"ass", true},
		{"sub", true},
		{"xxx", false},
		{"all", true},
		{"", true},
	}

	for _, tc := range tt {
		if got := ValidType(tc.typ); got != tc.expected {
			t.Errorf("ValidType(%q) = %v, want %v", tc.typ, got, tc.expected)
		}
	}
}

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := Config{Language: "en", Type: "srt"}
		if err := ValidateConfig(cfg); err != nil {
			t.Errorf("got error %v, want nil", err)
		}
	})

	t.Run("invalid language", func(t *testing.T) {
		cfg := Config{Language: "xx", Type: "srt"}
		if err := ValidateConfig(cfg); err == nil {
			t.Error("expected error for invalid language")
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		cfg := Config{Language: "en", Type: "xxx"}
		if err := ValidateConfig(cfg); err == nil {
			t.Error("expected error for invalid type")
		}
	})
}
