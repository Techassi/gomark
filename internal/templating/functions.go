package templating

import (
	"strings"
)

var alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func FormatNoImage(u string) string {
	return string(u[0])
}

func FormatColor(c, url string) string {
	var (
		color     string = ""
		charIndex int    = strings.Index(alphabet, strings.ToUpper(c))
	)

	if url != "" {
		return " white"
	}

	switch {
	case charIndex < 7:
		color = ""
	case charIndex < 14:
		color = " blue"
	case charIndex < 21:
		color = " orange"
	default:
		color = " green"
	}
	return color
}
