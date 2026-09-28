package subtitles

// DownloadMode controls subtitle auto-download behavior
type DownloadMode int

const (
	DownloadDisabled DownloadMode = iota
	DownloadEnabled
	DownloadAuto
)

// Config holds subtitle download settings
type Config struct {
	Mode     DownloadMode
	Language string // ISO 639-1 code: "en", "es", "fr", etc.
	Type     string // "srt", "vtt", "all"
	SearchText string // optional free text override
}

// ModeString returns a human-readable string for the mode
func (m DownloadMode) String() string {
	switch m {
	case DownloadDisabled:
		return "Disabled"
	case DownloadEnabled:
		return "Enabled"
	case DownloadAuto:
		return "Auto"
	default:
		return "Unknown"
	}
}
