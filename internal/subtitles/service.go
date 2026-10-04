package subtitles

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

// Service coordinates subtitle search, download, and selection
type Service struct {
	provider Provider
	config   Config
	cacheDir string
}

// NewService creates a subtitle service
func NewService(provider Provider, config Config, cacheDir string) *Service {
	return &Service{
		provider: provider,
		config:   config,
		cacheDir: cacheDir,
	}
}

// SearchResult contains search results and best candidate
type SearchResult struct {
	Candidates []SubtitleCandidate
	Best       *SubtitleCandidate
	Query      string
}

// TryAutoDownload attempts to search and download subtitle without user interaction
// Returns local file path or error
func (s *Service) TryAutoDownload(ctx context.Context, mediaPath string) (string, error) {
	if s.config.Mode == DownloadDisabled {
		return "", fmt.Errorf("download: disabled")
	}

	// Build query from media filename or use override
	query := s.config.SearchText
	if query == "" {
		query = BuildQueryFromPath(mediaPath)
	}

	if query == "" {
		return "", fmt.Errorf("download: empty query from %s", mediaPath)
	}

	// Search
	req := SearchRequest{
		Query:    query,
		Language: s.config.Language,
		Type:     s.config.Type,
	}

	candidates, err := s.provider.Search(ctx, req)
	if err != nil {
		return "", fmt.Errorf("download: search failed: %w", err)
	}

	if len(candidates) == 0 {
		return "", fmt.Errorf("download: no candidates found")
	}

	// Sort and select best
	sorted := SortByScore(candidates)
	best := SelectBest(sorted)
	if best == nil {
		return "", fmt.Errorf("download: no valid candidate selected")
	}

	// Download
	return s.provider.Download(ctx, *best, s.cacheDir)
}

// SearchAndSelect performs full search with result inspection
func (s *Service) SearchAndSelect(ctx context.Context, mediaPath string) (*SearchResult, error) {
	if s.config.Mode == DownloadDisabled {
		return nil, fmt.Errorf("search: download disabled")
	}

	// Build query
	query := s.config.SearchText
	if query == "" {
		query = BuildQueryFromPath(mediaPath)
	}

	if query == "" {
		return nil, fmt.Errorf("search: empty query from %s", mediaPath)
	}

	// Search
	req := SearchRequest{
		Query:    query,
		Language: s.config.Language,
		Type:     s.config.Type,
	}

	candidates, err := s.provider.Search(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	// Sort
	sorted := SortByScore(candidates)

	return &SearchResult{
		Candidates: sorted,
		Best:       SelectBest(sorted),
		Query:      query,
	}, nil
}

// Download retrieves a specific subtitle candidate
func (s *Service) Download(ctx context.Context, candidate SubtitleCandidate) (string, error) {
	if s.cacheDir == "" {
		tmpDir := os.TempDir()
		s.cacheDir = filepath.Join(tmpDir, "go2tv-subtitles")
	}

	return s.provider.Download(ctx, candidate, s.cacheDir)
}

// CanAutoDownload checks if auto-download is enabled
func (s *Service) CanAutoDownload() bool {
	return s.config.Mode == DownloadAuto
}

// IsEnabled checks if any download mode is active
func (s *Service) IsEnabled() bool {
	return s.config.Mode != DownloadDisabled
}

// SetConfig updates service configuration
func (s *Service) SetConfig(cfg Config) {
	s.config = cfg
}
