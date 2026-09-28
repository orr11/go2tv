package subtitles

import (
	"fmt"
	"strings"
)

// ValidLanguage checks if language code is valid
func ValidLanguage(lang string) bool {
	if lang == "" || lang == "all" {
		return true
	}
	// ISO 639-1 two-letter codes
	lang = strings.ToLower(lang)
	valid := map[string]bool{
		"en": true, "es": true, "fr": true, "de": true, "it": true,
		"pt": true, "ru": true, "ja": true, "ko": true, "zh": true,
		"ar": true, "hi": true, "pl": true, "tr": true, "nl": true,
		"sv": true, "no": true, "da": true, "fi": true, "el": true,
		"cs": true, "sk": true, "hu": true, "ro": true, "th": true,
	}
	return valid[lang]
}

// ValidType checks if subtitle type is valid
func ValidType(typ string) bool {
	if typ == "" || typ == "all" {
		return true
	}
	typ = strings.ToLower(typ)
	return typ == "srt" || typ == "vtt" || typ == "ass" || typ == "sub"
}

// ValidateConfig checks config for errors
func ValidateConfig(cfg Config) error {
	if !ValidLanguage(cfg.Language) {
		return fmt.Errorf("invalid language code: %s", cfg.Language)
	}
	if !ValidType(cfg.Type) {
		return fmt.Errorf("invalid subtitle type: %s", cfg.Type)
	}
	return nil
}
