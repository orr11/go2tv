package subtitles

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	openSubtitlesAPIURL = "https://api.opensubtitles.com/api/v1"
	userAgent           = "go2tv/1.0"
	requestTimeout      = 10 * time.Second
)

// OpenSubtitles API client
type OpenSubtitles struct {
	apiKey string
	client *http.Client
}

// NewOpenSubtitles creates a new OpenSubtitles provider
func NewOpenSubtitles(apiKey string) *OpenSubtitles {
	return &OpenSubtitles{
		apiKey: apiKey,
		client: &http.Client{Timeout: requestTimeout},
	}
}

// searchResponse matches OpenSubtitles API response
type searchResponse struct {
	Data []struct {
		Attributes struct {
			Subtitle struct {
				URL      string `json:"url"`
				Encoding string `json:"encoding"`
			} `json:"subtitle"`
			ExternalURL string  `json:"external_url"`
			Name        string  `json:"name"`
			Language    string  `json:"language"`
			Trusted     bool    `json:"trusted"`
			Score       float64 `json:"score"`
		} `json:"attributes"`
	} `json:"data"`
}

// Search queries OpenSubtitles for matching subtitles
func (o *OpenSubtitles) Search(ctx context.Context, req SearchRequest) ([]SubtitleCandidate, error) {
	if req.Query = NormalizeQuery(req.Query); req.Query == "" {
		return nil, &ErrorNoResults{Query: req.Query}
	}

	queryParams := url.Values{}
	queryParams.Set("query", req.Query)
	if req.Language != "" && req.Language != "all" {
		queryParams.Set("languages", req.Language)
	}

	apiURL := openSubtitlesAPIURL + "/subtitles?" + queryParams.Encode()

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	httpReq.Header.Set("User-Agent", userAgent)
	httpReq.Header.Set("Accept", "application/json")
	if o.apiKey != "" {
		httpReq.Header.Set("Api-Key", o.apiKey)
	}

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search: HTTP %d: %s", resp.StatusCode, string(body))
	}

	var sr searchResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return nil, fmt.Errorf("search: %w", err)
	}

	if len(sr.Data) == 0 {
		return nil, &ErrorNoResults{Query: req.Query}
	}

	var candidates []SubtitleCandidate
	for _, item := range sr.Data {
		attrs := item.Attributes
		subType := detectType(attrs.Subtitle.Encoding)
		if req.Type != "" && req.Type != "all" && !strings.EqualFold(subType, req.Type) {
			continue
		}

		candidates = append(candidates, SubtitleCandidate{
			FileName:    attrs.Name,
			Language:    attrs.Language,
			Type:        subType,
			DownloadURL: attrs.Subtitle.URL,
			Score:       attrs.Score,
			IsTrusted:   attrs.Trusted,
		})
	}

	if len(candidates) == 0 {
		return nil, &ErrorNoResults{Query: req.Query}
	}

	return candidates, nil
}

// Download retrieves and saves a subtitle file
func (o *OpenSubtitles) Download(ctx context.Context, candidate SubtitleCandidate, destDir string) (string, error) {
	if candidate.DownloadURL == "" {
		return "", fmt.Errorf("download: empty download URL")
	}

	if err := os.MkdirAll(destDir, 0755); err != nil {
		return "", fmt.Errorf("download: %w", err)
	}

	// Sanitize filename
	fileName := filepath.Base(candidate.FileName)
	if fileName == "" || fileName == "." {
		ext := ".srt"
		if candidate.Type != "" {
			ext = "." + candidate.Type
		}
		fileName = "subtitle" + ext
	}

	destPath := filepath.Join(destDir, fileName)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, candidate.DownloadURL, nil)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}

	httpReq.Header.Set("User-Agent", userAgent)

	resp, err := o.client.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download: HTTP %d", resp.StatusCode)
	}

	file, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("download: %w", err)
	}
	defer file.Close()

	if _, err := io.Copy(file, resp.Body); err != nil {
		os.Remove(destPath)
		return "", fmt.Errorf("download: %w", err)
	}

	return destPath, nil
}

// detectType guesses subtitle type from encoding or filename
func detectType(encoding string) string {
	enc := strings.ToLower(encoding)

	switch {
	case strings.Contains(enc, "webvtt"), strings.Contains(enc, "vtt"):
		return "vtt"
	case strings.Contains(enc, "ass"), strings.Contains(enc, "ssa"):
		return "ass"
	case strings.Contains(enc, "srt"), strings.Contains(enc, "subrip"):
		return "srt"
	case strings.Contains(enc, "sub"):
		return "sub"
	default:
		return "srt"
	}
}
