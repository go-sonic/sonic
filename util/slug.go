package util

import (
	"regexp"
	"strconv"
	"time"
)

var (

	// r1 Match non-English and non-Chinese characters
	r1 = regexp.MustCompile(`[^(a-zA-Z0-9\\u4e00-\\u9fa5\.\-)]`)
	// r2 Match special symbol
	r2 = regexp.MustCompile(`[\?\\/:|<>\*\[\]\(\)\$%\{\}@~\.]`)
	// r3 Match whitespace characters
	r3 = regexp.MustCompile(`\s`)
)

func Slug(slug string) string {
	slug = r1.ReplaceAllString(slug, "")
	slug = r2.ReplaceAllString(slug, "")
	slug = r3.ReplaceAllString(slug, "")
	if slug == "" {
		slug = strconv.FormatInt(time.Now().UnixMilli(), 10)
	}
	return slug
}
