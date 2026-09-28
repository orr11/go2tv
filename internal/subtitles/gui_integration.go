//go:build !(android || ios)

package subtitles

import (
	"context"
	"fmt"
)

// GUISubtitleSearchResult holds results for GUI display
type GUISubtitleSearchResult struct {
	Candidates []GUICandidateDisplay
	BestIdx    int // Index of best candidate, -1 if none
	Query      string
	Error      error
}

// GUICandidateDisplay formats candidate for GUI list
type GUICandidateDisplay struct {
	FileName string
	Language string
	Type     string
	Trusted  bool
	Score    float64
	Internal *SubtitleCandidate // Keep reference for download
}

// FormatCandidateForDisplay creates GUI-friendly display
func FormatCandidateForDisplay(c SubtitleCandidate) GUICandidateDisplay {
	trustedMark := ""
	if c.IsTrusted {
		trustedMark = " ✓"
	}

	return GUICandidateDisplay{
		FileName: fmt.Sprintf("%s [%s]%s (%.0f%%)", c.FileName, c.Language, trustedMark, c.Score*100),
		Language: c.Language,
		Type:     c.Type,
		Trusted:  c.IsTrusted,
		Score:    c.Score,
		Internal: &c,
	}
}

// PerformGUISearch executes search and formats for GUI
func PerformGUISearch(ctx context.Context, svc *Service, mediaPath string) *GUISubtitleSearchResult {
	result, err := svc.SearchAndSelect(ctx, mediaPath)
	if err != nil {
		return &GUISubtitleSearchResult{Error: err}
	}

	displayItems := make([]GUICandidateDisplay, len(result.Candidates))
	for i, c := range result.Candidates {
		displayItems[i] = FormatCandidateForDisplay(c)
	}

	bestIdx := -1
	if result.Best != nil {
		for i, c := range result.Candidates {
			if c.ID == result.Best.ID {
				bestIdx = i
				break
			}
		}
	}

	return &GUISubtitleSearchResult{
		Candidates: displayItems,
		BestIdx:    bestIdx,
		Query:      result.Query,
	}
}
