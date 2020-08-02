package utils

import (
    "fmt"
    "regexp"
    "strings"
)

var removeTag = regexp.MustCompile("<[^>]*>")
var removeBlank = regexp.MustCompile("\n+\\s*")

func CleanHtmlContent(raw string, contentLength int) string {
    for _, breakTag := range []string{"li", "p"} {
        raw = strings.ReplaceAll(raw, fmt.Sprintf("</%s>", breakTag), "\n")
    }
    raw = strings.ReplaceAll(raw, "<li>", "- ")

    raw = removeTag.ReplaceAllString(raw, "")
    raw = removeBlank.ReplaceAllString(raw, "\n")

    r := []rune(raw)

    if len(r) > contentLength {
        r = r[:contentLength]
    }

    return strings.TrimSuffix(string(r), "\n") + "......"
}
