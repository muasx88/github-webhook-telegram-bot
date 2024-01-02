package utils

import "strings"

func EscapeHTML(text string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, "&", "&amp;"), "<", "&lt;"), ">", "&gt;")
}
