package subtitles

import (
	"path/filepath"
	"regexp"
	"strings"
)

// BuildQueryFromPath normalizes filename into a search query
// "movie.2024.1080p.mkv" -> "movie"
// "The Matrix 1999 english.srt" -> "The Matrix 1999"
func BuildQueryFromPath(filePath string) string {
	base := filepath.Base(filePath)
	// Remove extension
	name := strings.TrimSuffix(base, filepath.Ext(base))
	// Lowercase and normalize whitespace
	name = strings.TrimSpace(strings.ToLower(name))
	// Replace common separators with spaces
	re := regexp.MustCompile(`[._\-]+`)
	name = re.ReplaceAllString(name, " ")
	// Remove common quality/resolution suffixes
	qualityRe := regexp.MustCompile(`(?i)(1080p|720p|480p|2160p|4k|8k|hd|sd|hdtv|bdrip|webrip|dvdrip|bluray|x264|x265|h264|h265|aac|ac3|dts)`)
	name = qualityRe.ReplaceAllString(name, " ")
	// Clean up multiple spaces
	spaceRe := regexp.MustCompile(`\s+`)
	name = spaceRe.ReplaceAllString(name, " ")
	return strings.TrimSpace(name)
}

// NormalizeQuery cleans up a search query
func NormalizeQuery(q string) string {
	q = strings.TrimSpace(q)
	q = strings.ToLower(q)
	// Collapse multiple spaces
	re := regexp.MustCompile(`\s+`)
	q = re.ReplaceAllString(q, " ")
	return q
}

// FilterByLanguage returns candidates matching the requested language
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

// FilterByType returns candidates matching the requested subtitle type
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

// SortByScore sorts candidates by match score (descending), then by download count
func SortByScore(candidates []SubtitleCandidate) []SubtitleCandidate {
	// Bubble sort for stable ordering (go stdlib not used to preserve order)
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

// SelectBest returns the first (best) candidate if available
func SelectBest(candidates []SubtitleCandidate) *SubtitleCandidate {
	if len(candidates) == 0 {
		return nil
	}
	return &candidates[0]
}
