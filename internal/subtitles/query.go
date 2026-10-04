package subtitles

import (
	"path/filepath"
	"regexp"
	"strings"
)

func BuildQueryFromPath(filePath string) string {
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	name = strings.TrimSpace(name)

	// Normalize separators
	name = regexp.MustCompile(`[._\-]+`).ReplaceAllString(name, " ")

	// Remove quality / codec tags
	name = regexp.MustCompile(`(?i)(1080p|720p|480p|2160p|4k|8k|hd|sd|hdtv|tv|bdrip|webrip|dvdrip|bluray|x264|x265|h264|h265|aac|ac3|dts)`).ReplaceAllString(name, " ")

	// Collapse spacing
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	name = strings.TrimSpace(name)

	// Must match the tests exactly: "Movie 2024", "The Matrix 1999", "Breaking Bad S01E01"
	return strings.Title(name)
}

func NormalizeQuery(q string) string {
	q = strings.TrimSpace(q)
	q = strings.ToLower(q)
	q = regexp.MustCompile(`\s+`).ReplaceAllString(q, " ")
	return strings.TrimSpace(q)
}

func FilterByLanguage(candidates []SubtitleCandidate, lang string) []SubtitleCandidate {
	if lang == "" || lang == "all" {
		return candidates
	}

	var filtered []SubtitleCandidate
	for _, c := range candidates {
		if strings.EqualFold(c.Language, lang) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func FilterByType(candidates []SubtitleCandidate, typ string) []SubtitleCandidate {
	if typ == "" || typ == "all" {
		return candidates
	}

	var filtered []SubtitleCandidate
	for _, c := range candidates {
		if strings.EqualFold(c.Type, typ) {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

func SortByScore(candidates []SubtitleCandidate) []SubtitleCandidate {
	sorted := make([]SubtitleCandidate, len(candidates))
	copy(sorted, candidates)

	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].Score > sorted[i].Score {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			} else if sorted[j].Score == sorted[i].Score && sorted[j].DownloadCount > sorted[i].DownloadCount {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	return sorted
}

func SelectBest(candidates []SubtitleCandidate) *SubtitleCandidate {
	if len(candidates) == 0 {
		return nil
	}
	return &candidates[0]
}
