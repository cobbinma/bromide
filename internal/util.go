package internal

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
)

func Diff(old, new string) string {
	edits := myers.ComputeEdits("current", old, new)
	diff := fmt.Sprint(gotextdiff.ToUnified("current", "incoming", old, edits))

	out := ""
	for _, line := range strings.Split(diff, "\n") {
		green := "\x1b[32m"
		red := "\x1b[31m"
		reset := "\x1b[0m"

		addition := strings.HasPrefix(line, "+")
		subtraction := strings.HasPrefix(line, "-")

		l := ""
		switch {
		case addition:
			l = green + line
		case subtraction:
			l = red + line
		default:
			l = reset + line
		}

		out = out + l + "\n"
	}

	return out
}

func Sanitize(input string) string {
	// Define a regular expression pattern to match invalid characters
	pattern := "[\\/\\\\\\x00:\\*\\?\"<>\\|&\\#]"

	// Compile the regular expression pattern
	regex := regexp.MustCompile(pattern)

	// Replace invalid characters with an empty string
	return regex.ReplaceAllString(input, "")
}
