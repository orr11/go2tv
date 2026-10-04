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
