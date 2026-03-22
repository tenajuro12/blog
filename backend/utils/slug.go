package utils

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

var nonAlphanumeric = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(title string) string {
	s := strings.ToLower(title)
	s = nonAlphanumeric.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return fmt.Sprintf("%s-%d", s, time.Now().UnixNano())
}
