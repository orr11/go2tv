package subtitles

import (
	"testing"
)

func TestBuildQueryFromPath(t *testing.T) {
	tt := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "quality suffix removal",
			path:     "movie.2024.1080p.mkv",
			expected: "movie 2024",
		},
		{
			name:     "codec removal",
			path:     "The.Matrix.1999.x264.mkv",
			expected: "The Matrix 1999",
		},
		{
			name:     "dot and dash normalization",
			path:     "The-Matrix-1999.mkv",
			expected: "The Matrix 1999",
		},
		{
			name:     "underscore normalization",
			path:     "The_Matrix_Reloaded_2003.mkv",
			expected: "The Matrix Reloaded 2003",
		},
		{
			name:     "hdtv codec",
			path:     "Breaking.Bad.S01E01.HDTV.x264.mkv",
			expected: "Breaking Bad S01E01",
		},
		{
			name:     "multiple quality tags",
			path:     "Movie.2024.1080p.BluRay.x264.mkv",
			expected: "Movie 2024",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildQueryFromPath(tc.path)
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestNormalizeQuery(t *testing.T) {
	tt := []struct {
		name     string
		query    string
		expected string
	}{
		{
			name:     "leading/trailing space",
			query:    "  The Matrix  ",
			expected: "the matrix",
		},
		{
			name:     "multiple internal spaces",
			query:    "The    Matrix    1999",
			expected: "the matrix 1999",
		},
		{
			name:     "mixed case",
			query:    "ThE MaTrIx",
			expected: "the matrix",
		},
		{
			name:     "empty string",
			query:    "",
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := NormalizeQuery(tc.query)
			if got != tc.expected {
				t.Errorf("got %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestFilterByLanguage(t *testing.T) {
	candidates := []SubtitleCandidate{
		{ID: "1", Language: "en"},
		{ID: "2", Language: "es"},
		{ID: "3", Language: "en"},
		{ID: "4", Language: "fr"},
	}

	t.Run("filter by language", func(t *testing.T) {
		got := FilterByLanguage(candidates, "en")
		if len(got) != 2 {
			t.Fatalf("got %d candidates, want 2", len(got))
		}
		if got[0].ID != "1" || got[1].ID != "3" {
			t.Errorf("got IDs %s,%s, want 1,3", got[0].ID, got[1].ID)
		}
	})

	t.Run("all languages", func(t *testing.T) {
		got := FilterByLanguage(candidates, "all")
		if len(got) != len(candidates) {
			t.Errorf("got %d candidates, want %d", len(got), len(candidates))
		}
	})

	t.Run("empty language", func(t *testing.T) {
		got := FilterByLanguage(candidates, "")
		if len(got) != len(candidates) {
			t.Errorf("got %d candidates, want %d", len(got), len(candidates))
		}
	})

	t.Run("no match", func(t *testing.T) {
		got := FilterByLanguage(candidates, "de")
		if len(got) != 0 {
			t.Errorf("got %d candidates, want 0", len(got))
		}
	})
}

func TestFilterByType(t *testing.T) {
	candidates := []SubtitleCandidate{
		{ID: "1", Type: "srt"},
		{ID: "2", Type: "vtt"},
		{ID: "3", Type: "srt"},
		{ID: "4", Type: "ass"},
	}

	t.Run("filter by type", func(t *testing.T) {
		got := FilterByType(candidates, "srt")
		if len(got) != 2 {
			t.Fatalf("got %d candidates, want 2", len(got))
		}
	})

	t.Run("all types", func(t *testing.T) {
		got := FilterByType(candidates, "all")
		if len(got) != len(candidates) {
			t.Errorf("got %d candidates, want %d", len(got), len(candidates))
		}
	})
}

func TestSortByScore(t *testing.T) {
	candidates := []SubtitleCandidate{
		{ID: "1", Score: 0.8, DownloadCount: 100},
		{ID: "2", Score: 0.95, DownloadCount: 200},
		{ID: "3", Score: 0.8, DownloadCount: 500},
		{ID: "4", Score: 0.7, DownloadCount: 1000},
	}

	sorted := SortByScore(candidates)

	if len(sorted) != len(candidates) {
		t.Fatalf("got %d candidates, want %d", len(sorted), len(candidates))
	}

	if sorted[0].ID != "2" {
		t.Errorf("first (highest score) is %s, want 2", sorted[0].ID)
	}

	// Check that ID 3 (0.8 score, 500 downloads) comes before ID 1 (0.8 score, 100 downloads)
	id1Idx, id3Idx := -1, -1
	for i, c := range sorted {
		if c.ID == "1" {
			id1Idx = i
		}
		if c.ID == "3" {
			id3Idx = i
		}
	}
	if id1Idx >= 0 && id3Idx >= 0 && id3Idx > id1Idx {
		t.Errorf("ID 3 (500 downloads) should come before ID 1 (100 downloads)")
	}
}

func TestSelectBest(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		got := SelectBest([]SubtitleCandidate{})
		if got != nil {
			t.Errorf("got %v, want nil", got)
		}
	})

	t.Run("single candidate", func(t *testing.T) {
		candidates := []SubtitleCandidate{
			{ID: "1", Score: 0.9},
		}
		got := SelectBest(candidates)
		if got == nil || got.ID != "1" {
			t.Errorf("got %v, want ID 1", got)
		}
	})

	t.Run("multiple candidates returns first", func(t *testing.T) {
		candidates := []SubtitleCandidate{
			{ID: "1", Score: 0.9},
			{ID: "2", Score: 0.8},
		}
		got := SelectBest(candidates)
		if got == nil || got.ID != "1" {
			t.Errorf("got %v, want ID 1", got)
		}
	})
}
