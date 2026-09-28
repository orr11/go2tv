package subtitles

import (
	"context"
)

// SubtitleCandidate represents a downloadable subtitle match
type SubtitleCandidate struct {
	ID            string
	Language      string
	LanguageName  string
	Type          string  // "srt", "vtt", "ass", etc.
	DownloadURL   string
	Score         float64 // 0.0-1.0 match confidence
	FileName      string
	IsTrusted     bool
	UploadCount   int
	DownloadCount int
}

// SearchRequest parameters for subtitle search
type SearchRequest struct {
	Query     string // filename or free text
	Language  string // ISO 639-1: "en", "es", "fr"
	Type      string // "srt", "vtt", "all"
	MovieHash string // optional SHA1 hash
	FileSize  int64  // optional file size in bytes
}

// Provider defines the subtitle provider interface
type Provider interface {
	// Search finds matching subtitles
	Search(ctx context.Context, req SearchRequest) ([]SubtitleCandidate, error)

	// Download retrieves subtitle file, returns local file path
	Download(ctx context.Context, candidate SubtitleCandidate, destDir string) (string, error)
}

// ErrorNoResults returned when search finds no matches
type ErrorNoResults struct {
	Query string
}

func (e *ErrorNoResults) Error() string {
	return "no subtitles found for: " + e.Query
}

// ErrNoMatches checks if error is a no-results error
func ErrNoMatches(err error) bool {
	_, ok := err.(*ErrorNoResults)
	return ok
}
