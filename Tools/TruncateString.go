package tools

// Truncate string.
// Found here: https://dev.to/takakd/go-safe-truncate-string-9h0
func TruncateString(str string, length int) string {
	length = length - 3
	if length <= 0 {
		return ""
	}

	didTruncate := false
	if length < len(str) {
		didTruncate = true
	}

	// This code cannot support Japanese
	// orgLen := len(str)
	// if orgLen <= length {
	//     return str
	// }
	// return str[:length]

	// Support Japanese
	// Ref: Range loops https://blog.golang.org/strings
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	if didTruncate {
		truncated = truncated + "..."
	}
	return truncated
}
