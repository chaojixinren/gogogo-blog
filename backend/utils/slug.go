package utils

import (
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`[^a-z0-9]+`)

func Slugify(value string) string {
	slug := strings.ToLower(value)
	slug = slugPattern.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")
	if slug == "" {
		return "item"
	}
	return slug
}
