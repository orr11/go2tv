package subtitles

import (
	"context"
	"testing"
	"time"
)

func TestIntegrationAutoDownloadFlow(t *testing.T) {
	mock := &mockProvider{
		searchResults: []SubtitleCandidate{
			{
				ID:          "1",
				FileName:    "movie.srt",
				Language:    "en",
				Type:        "srt",
				Score:       0.95,
				DownloadURL: "http://example.com/sub.srt",
				IsTrusted:   true,
			},
		},
		downloadPath: "/tmp/sub.srt",
	}

	cfg := Config{Mode: DownloadAuto, Language: "en", Type: "srt"}
	svc := NewService(mock, cfg, "/tmp")
	mediaPath := "breaking.bad.s01e01.1080p.mkv"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	subtitlePath, err := svc.TryAutoDownload(ctx, mediaPath)
	if err != nil {
		t.Fatalf("auto download failed: %v", err)
	}

	if subtitlePath != "/tmp/sub.srt" {
		t.Errorf("got subtitle path %s, want /tmp/sub.srt", subtitlePath)
	}
}

func TestIntegrationEnabledModeFlow(t *testing.T) {
	mock := &mockProvider{
		searchResults: []SubtitleCandidate{
			{ID: "1", FileName: "sub1.srt", Language: "en", Score: 0.95},
			{ID: "2", FileName: "sub2.srt", Language: "en", Score: 0.85},
			{ID: "3", FileName: "sub3.srt", Language: "es", Score: 0.90},
		},
		downloadPath: "/tmp/selected.srt",
	}

	cfg := Config{Mode: DownloadEnabled, Language: "en"}
	svc := NewService(mock, cfg, "/tmp")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := svc.SearchAndSelect(ctx, "movie.mkv")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}

	if len(result.Candidates) != 3 {
		t.Errorf("got %d candidates, want 3", len(result.Candidates))
	}

	if result.Best.ID != "1" {
		t.Errorf("best is %s, want 1", result.Best.ID)
	}

	path, err := svc.Download(ctx, *result.Best)
	if err != nil {
		t.Fatalf("download failed: %v", err)
	}

	if path != "/tmp/selected.srt" {
		t.Errorf("download path is %s, want /tmp/selected.srt", path)
	}
}

func TestIntegrationDisabledMode(t *testing.T) {
	cfg := Config{Mode: DownloadDisabled}
	svc := NewService(&mockProvider{}, cfg, "/tmp")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := svc.SearchAndSelect(ctx, "movie.mkv")
	if err == nil {
		t.Fatal("search should fail when disabled")
	}

	_, err = svc.TryAutoDownload(ctx, "movie.mkv")
	if err == nil {
		t.Fatal("auto-download should fail when disabled")
	}
}
