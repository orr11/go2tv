package subtitles

import (
	"context"
	"errors"
	"testing"
	"time"
)

// mockProvider implements Provider for testing
type mockProvider struct {
	searchResults []SubtitleCandidate
	searchErr     error
	downloadPath  string
	downloadErr   error
}

func (m *mockProvider) Search(ctx context.Context, req SearchRequest) ([]SubtitleCandidate, error) {
	if m.searchErr != nil {
		return nil, m.searchErr
	}
	return m.searchResults, nil
}

func (m *mockProvider) Download(ctx context.Context, candidate SubtitleCandidate, destDir string) (string, error) {
	if m.downloadErr != nil {
		return "", m.downloadErr
	}
	return m.downloadPath, nil
}

func TestServiceTryAutoDownload(t *testing.T) {
	t.Run("disabled mode", func(t *testing.T) {
		svc := NewService(&mockProvider{}, Config{Mode: DownloadDisabled}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := svc.TryAutoDownload(ctx, "movie.mkv")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
	})

	t.Run("auto download success", func(t *testing.T) {
		mock := &mockProvider{
			searchResults: []SubtitleCandidate{
				{
					ID:          "1",
					FileName:    "movie.srt",
					Language:    "en",
					Type:        "srt",
					Score:       0.95,
					DownloadURL: "http://example.com/sub.srt",
				},
			},
			downloadPath: "/tmp/movie.srt",
		}

		svc := NewService(mock, Config{Mode: DownloadAuto, Language: "en"}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		path, err := svc.TryAutoDownload(ctx, "movie.mkv")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if path != "/tmp/movie.srt" {
			t.Errorf("got path %s, want /tmp/movie.srt", path)
		}
	})

	t.Run("search error propagation", func(t *testing.T) {
		mock := &mockProvider{
			searchErr: errors.New("API error"),
		}

		svc := NewService(mock, Config{Mode: DownloadAuto}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := svc.TryAutoDownload(ctx, "movie.mkv")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("no results found", func(t *testing.T) {
		mock := &mockProvider{
			searchResults: []SubtitleCandidate{},
		}

		svc := NewService(mock, Config{Mode: DownloadAuto}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := svc.TryAutoDownload(ctx, "movie.mkv")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("respects language filter", func(t *testing.T) {
		mock := &mockProvider{
			searchResults: []SubtitleCandidate{
				{ID: "1", Language: "en", Score: 0.95},
				{ID: "2", Language: "es", Score: 0.90},
			},
			downloadPath: "/tmp/sub.srt",
		}

		svc := NewService(mock, Config{Mode: DownloadAuto, Language: "en"}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := svc.TryAutoDownload(ctx, "movie.mkv")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestServiceSearchAndSelect(t *testing.T) {
	t.Run("returns sorted candidates", func(t *testing.T) {
		mock := &mockProvider{
			searchResults: []SubtitleCandidate{
				{ID: "1", Language: "en", Score: 0.8},
				{ID: "2", Language: "en", Score: 0.95},
				{ID: "3", Language: "en", Score: 0.7},
			},
		}

		svc := NewService(mock, Config{Mode: DownloadEnabled, Language: "en"}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		result, err := svc.SearchAndSelect(ctx, "movie.mkv")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Candidates) != 3 {
			t.Errorf("got %d candidates, want 3", len(result.Candidates))
		}

		if result.Best == nil || result.Best.ID != "2" {
			t.Errorf("best candidate ID is %v, want 2", result.Best)
		}
	})

	t.Run("disabled mode returns error", func(t *testing.T) {
		svc := NewService(&mockProvider{}, Config{Mode: DownloadDisabled}, "/tmp")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		_, err := svc.SearchAndSelect(ctx, "movie.mkv")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestServiceIsEnabled(t *testing.T) {
	tt := []struct {
		mode     DownloadMode
		expected bool
	}{
		{DownloadDisabled, false},
		{DownloadEnabled, true},
		{DownloadAuto, true},
	}

	for _, tc := range tt {
		t.Run(tc.mode.String(), func(t *testing.T) {
			svc := NewService(&mockProvider{}, Config{Mode: tc.mode}, "/tmp")
			if svc.IsEnabled() != tc.expected {
				t.Errorf("IsEnabled() = %v, want %v", svc.IsEnabled(), tc.expected)
			}
		})
	}
}

func TestServiceCanAutoDownload(t *testing.T) {
	tt := []struct {
		mode     DownloadMode
		expected bool
	}{
		{DownloadDisabled, false},
		{DownloadEnabled, false},
		{DownloadAuto, true},
	}

	for _, tc := range tt {
		t.Run(tc.mode.String(), func(t *testing.T) {
			svc := NewService(&mockProvider{}, Config{Mode: tc.mode}, "/tmp")
			if svc.CanAutoDownload() != tc.expected {
				t.Errorf("CanAutoDownload() = %v, want %v", svc.CanAutoDownload(), tc.expected)
			}
		})
	}
}

func TestOpenSubtitlesDetectType(t *testing.T) {
	tt := []struct {
		encoding string
		expected string
	}{
		{"utf-8", "srt"},
		{"webvtt", "vtt"},
		{"vtt", "vtt"},
		{"ass", "ass"},
		{"ssa", "ass"},
		{"subrip", "srt"},
		{"sub", "sub"},
	}

	for _, tc := range tt {
		t.Run(tc.encoding, func(t *testing.T) {
			got := detectType(tc.encoding)
			if got != tc.expected {
				t.Errorf("got %s, want %s", got, tc.expected)
			}
		})
	}
}

func TestErrorNoMatches(t *testing.T) {
	t.Run("detects no-results error", func(t *testing.T) {
		err := &ErrorNoResults{Query: "test"}
		if !ErrNoMatches(err) {
			t.Error("ErrNoMatches should return true for ErrorNoResults")
		}
	})

	t.Run("ignores other errors", func(t *testing.T) {
		err := errors.New("other error")
		if ErrNoMatches(err) {
			t.Error("ErrNoMatches should return false for other errors")
		}
	})
}
