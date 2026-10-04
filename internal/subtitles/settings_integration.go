package subtitles

import (
	"path/filepath"
	"os"
)

// AppSettings extends app settings with subtitle download config
type AppSettings struct {
	SubtitleMode     DownloadMode
	SubtitleLanguage string
	SubtitleType     string
	SubtitleCacheDir string
}

// DefaultAppSettings returns sensible defaults
func DefaultAppSettings() AppSettings {
	tmpDir := os.TempDir()
	subDir := filepath.Join(tmpDir, "go2tv-subtitles")

	return AppSettings{
		SubtitleMode:     DownloadDisabled,
		SubtitleLanguage: "en",
		SubtitleType:     "srt",
		SubtitleCacheDir: subDir,
	}
}

// ToConfig converts AppSettings to subtitle service Config
func (as AppSettings) ToConfig() Config {
	return Config{
		Mode:     as.SubtitleMode,
		Language: as.SubtitleLanguage,
		Type:     as.SubtitleType,
	}
}
