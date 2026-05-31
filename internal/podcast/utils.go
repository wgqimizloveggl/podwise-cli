package podcast

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// isPodwiseHost reports whether host is a recognised Podwise web host.
// Accepted: podwise.ai, app.podwise.ai, beta.podwise.ai.
func isPodwiseHost(host string) bool {
	switch host {
	case "podwise.ai", "app.podwise.ai", "beta.podwise.ai":
		return true
	}
	return false
}

// ParseSeq extracts the integer podcast seq from a podwise podcast URL.
// Expected format: https://podwise.ai/dashboard/podcasts/<seq>
// Also accepts https://app.podwise.ai/dashboard/podcasts/<seq> and beta.podwise.ai.
func ParseSeq(input string) (int, error) {
	const hint = "(expected https://podwise.ai/dashboard/podcasts/<id>)"

	u, err := url.Parse(input)
	if err != nil || u.Scheme != "https" || !isPodwiseHost(u.Host) {
		return 0, fmt.Errorf("%q is not a valid podwise podcast URL %s", input, hint)
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) != 3 || parts[0] != "dashboard" || parts[1] != "podcasts" || parts[2] == "" {
		return 0, fmt.Errorf("%q is not a valid podwise podcast URL %s", input, hint)
	}

	seq, err := strconv.Atoi(parts[2])
	if err != nil || seq <= 0 {
		return 0, fmt.Errorf("podcast ID %q is not a positive integer %s", parts[2], hint)
	}
	return seq, nil
}

// BuildPodcastURL builds a podwise podcast URL from a sequence number.
func BuildPodcastURL(seq int) string {
	return fmt.Sprintf("https://podwise.ai/dashboard/podcasts/%d", seq)
}
