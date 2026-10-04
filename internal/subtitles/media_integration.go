package subtitles

import (
	"path/filepath"
	"strings"
)

// MediaSubtitle represents a subtitle attached to media
type MediaSubtitle struct {
	Path         string // Local file path
	Language     string
	Type         string // srt, vtt, etc.
	IsDownloaded bool   // true if fetched from OpenSubtitles
	SourceName   string // Original filename from OpenSubtitles
}

// MediaItemWithSubtitles extends media item with subtitle metadata
type MediaItemWithSubtitles struct {
	MediaPath   string
	Subtitles   []MediaSubtitle
	SelectedIdx int // Index of selected subtitle, -1 for none
}

// HasDownloadedSubtitle checks if any subtitle was auto-downloaded
func (m MediaItemWithSubtitles) HasDownloadedSubtitle() bool {
	for _, sub := range m.Subtitles {
		if sub.IsDownloaded {
			return true
		}
	}
	return false
}

// SelectedSubtitle returns the currently selected subtitle
func (m MediaItemWithSubtitles) SelectedSubtitle() *MediaSubtitle {
	if m.SelectedIdx < 0 || m.SelectedIdx >= len(m.Subtitles) {
		return nil
	}
	return &m.Subtitles[m.SelectedIdx]
}

// IsSubtitlePath checks if path is a subtitle file
func IsSubtitlePath(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".srt" || ext == ".vtt" || ext == ".ass" || ext == ".sub"
}
