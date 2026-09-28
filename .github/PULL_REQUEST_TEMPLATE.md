# OpenSubtitles Download Integration

## Summary

Adds automatic subtitle downloading from OpenSubtitles API with three modes: Disabled, Enabled (user-driven), and Auto (automatic).

## Features

- **Three Download Modes**:
  - `Disabled`: No subtitle downloads (default)
  - `Enabled`: User searches and selects subtitles manually
  - `Auto`: Automatically downloads best-matching subtitle without prompting

- **Easy Configuration**:
  - Default search language (ISO 639-1: en, es, fr, etc.)
  - Subtitle type preference (srt, vtt, all)
  - Optional free-text query override

- **Smart Filtering & Selection**:
  - Query normalization removes quality tags and codecs from filenames
  - Results sorted by match confidence and download popularity
  - Trusted uploaders marked for visibility

- **Seamless Integration**:
  - Downloaded subtitles cached locally
  - Automatically attached to stream menu when selected
  - No re-selection needed after download
  - Works with both DLNA and Chromecast

- **Zero Network Calls in Tests**:
  - All unit tests use mocked provider
  - No live API calls during `go test ./...`
  - Comprehensive test coverage for query, filtering, and service logic

## Implementation Details

### New Package: `internal/subtitles/`

**Core Files:**
- `config.go` - Download mode enum and settings
- `provider.go` - Provider interface and subtitle types
- `query.go` - Query normalization, filtering, sorting
- `opensubtitles.go` - OpenSubtitles API v1 client
- `service.go` - High-level service coordinating search/download

**Integration Files:**
- `settings_integration.go` - App settings extension
- `gui_integration.go` - GUI display formatting (desktop only)
- `media_integration.go` - Media + subtitle metadata
- `validation.go` - Language/type validation

**Tests:**
- `query_test.go` - Query normalization & filtering tests
- `subtitles_test.go` - Service & mock provider tests
- `integration_test.go` - End-to-end mode tests
- `validation_test.go` - Config validation tests

### Usage Flow

**Auto Mode:**
```go
svc := subtitles.NewService(provider, Config{Mode: DownloadAuto, Language: "en"}, cacheDir)
path, err := svc.TryAutoDownload(ctx, "/path/to/movie.mkv")
// subtitle automatically downloaded and ready to use
```

**Enabled Mode:**
```go
result, err := svc.SearchAndSelect(ctx, "/path/to/movie.mkv")
// Present result.Candidates to user, let them pick
path, err := svc.Download(ctx, result.Candidates[userChoice])
```

## Testing

```bash
# Run all subtitle tests
go test -v ./internal/subtitles

# Run full suite
go test -v ./...
```

All tests pass with mocked provider (no network calls).

## Architecture Notes

- **Provider Pattern**: Interface allows swapping providers (OpenSubtitles, Subscene, etc.) without UI changes
- **Service Layer**: Abstracts search/download/filtering logic from UI concerns
- **No Blocking UI**: All network operations use context with timeouts
- **Platform Support**: Desktop GUI only (android/ios builds skip subtitle UI)
- **Cache Management**: Downloaded subtitles stored in temp directory, cleaned on exit or when cache grows

## Configuration Storage

Settings stored in Fyne preferences:
- `SubtitleMode`: Download mode (Disabled/Enabled/Auto)
- `SubtitleLanguage`: Default search language
- `SubtitleType`: Preferred subtitle format
- `SubtitleCacheDir`: Where to store downloaded files

## Future Enhancements

- [ ] Support additional providers (Subscene, Podnapisi, etc.)
- [ ] Movie hash-based search (SHA1 of file)
- [ ] Batch download for playlists
- [ ] Subtitle format conversion on download
- [ ] Search history & favorite uploads

## Checklist

- [x] Code follows AGENTS.md style guidelines
- [x] All tests pass (no live API calls)
- [x] Query normalization tested
- [x] Filtering & sorting tested
- [x] Service layer tested (auto + enabled modes)
- [x] Config validation tested
- [x] Desktop GUI integration prepared
- [x] No breaking changes to existing code
