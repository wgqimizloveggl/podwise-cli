package episode

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// ParseSeq extracts the integer episode seq from a podwise episode URL.
// Expected format: https://podwise.ai/dashboard/episodes/<seq>
// Also accepts https://app.podwise.ai/dashboard/episodes/<seq> and beta.podwise.ai.
func ParseSeq(input string) (int, error) {
	const hint = "(expected https://podwise.ai/dashboard/episodes/<id>)"

	u, err := url.Parse(input)
	if err != nil || u.Scheme != "https" || !isPodwiseHost(u.Host) {
		return 0, fmt.Errorf("%q is not a valid podwise episode URL %s", input, hint)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "dashboard" || parts[1] != "episodes" || parts[2] == "" {
		return 0, fmt.Errorf("%q is not a valid podwise episode URL %s", input, hint)
	}

	seq, err := strconv.Atoi(parts[2])
	if err != nil || seq <= 0 {
		return 0, fmt.Errorf("episode ID %q is not a positive integer %s", parts[2], hint)
	}
	return seq, nil
}

// BuildEpisodeURL builds a podwise episode URL from a sequence number.
func BuildEpisodeURL(seq int) string {
	return fmt.Sprintf("https://podwise.ai/dashboard/episodes/%d", seq)
}

// IsYouTubeURL reports whether rawURL points to a YouTube video.
// Recognised hosts: youtube.com (and www.), youtu.be.
func IsYouTubeURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" {
		return false
	}
	switch u.Hostname() {
	case "youtube.com", "www.youtube.com":
		return u.Query().Get("v") != ""
	case "youtu.be":
		return len(u.Path) > 1
	}
	return false
}

// IsLocalMediaFile reports whether path refers to an existing regular file
// (not a directory or URL). Extension validation is left to episode.Upload
// so that the error message lists all supported formats.
func IsLocalMediaFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// IsXiaoyuzhouURL reports whether rawURL points to a Xiaoyuzhou episode.
func IsXiaoyuzhouURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" {
		return false
	}
	return u.Hostname() == "www.xiaoyuzhoufm.com" && strings.HasPrefix(u.Path, "/episode/")
}

// isPodwiseHost reports whether host is a recognised Podwise web host.
// Accepted: podwise.ai, app.podwise.ai, beta.podwise.ai.
func isPodwiseHost(host string) bool {
	switch host {
	case "podwise.ai", "app.podwise.ai", "beta.podwise.ai":
		return true
	}
	return false
}

// trimTime removes a leading "00:" hour prefix from a time string only when
// the remainder is still a valid mm:ss[-based] string (contains at least one
// more colon). Examples:
//
//	"00:01:02" → "01:02"
//	"00:00:05" → "00:05"
//	"01:02:03" → "01:02:03"  (hour is non-zero, unchanged)
//	"00:05"    → "00:05"     (already mm:ss, unchanged)
func trimTime(s string) string {
	if strings.HasPrefix(s, "00:") && strings.ContainsRune(s[3:], ':') {
		return s[3:]
	}
	return s
}
