package fetcher

import "strings"

var filenameReplacer = strings.NewReplacer(
	"/", "__",
	"\"", "__",
	"<", "__",
	">", "__",
	":", "__",
	"*", "__",
	"?", "__",
	"|", "__",
)

// parseFilename - Given a URL, parses the filename base for that URL
func parseFilename(u string) string {
	// Strip protocol, as it shouldn't exist in the filename
	fn := stripProtocol(u)

	// Remove the last slash, as pages ending in a slash should be treated as the same page
	if l := len(fn) - 1; l > 0 && fn[l:] == "/" {
		fn = fn[:l]
	}

	return filenameReplacer.Replace(fn)
}

// stripProtocol - Given a URL, strips the protocol off for recording
func stripProtocol(u string) string {
	parts := strings.SplitN(u, "://", 2)
	switch len(parts) {
	case 0:
		return "invalid-path"
	}
	return strings.TrimSpace(parts[len(parts)-1])
}
