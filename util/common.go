package util

import (
	"errors"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
	"github.com/go-sonic/sonic/handler/trans"
	"github.com/go-sonic/sonic/util/xerr"
)

func IfElse(condition bool, a, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func CompositeURL(urls ...string) string {
	builder := strings.Builder{}

	for i, url := range urls {
		if url == "" {
			continue
		}
		url = strings.TrimSuffix(url, "/")
		url = strings.TrimPrefix(url, "/")
		if i != 0 {
			builder.WriteString("/")
		}
		builder.WriteString(url)
	}
	return builder.String()
}

var htmlRegexp = regexp.MustCompile(`(<[^<]*?>)|(<[\s]*?/[^<]*?>)|(<[^<]*?/[\s]*?>)`)

func CleanHTMLTag(htmlContent string) string {
	if htmlContent == "" {
		return ""
	}
	return htmlRegexp.ReplaceAllString(htmlContent, "")
}

var blankRegexp = regexp.MustCompile(`\s`)

func HTMLFormatWordCount(html string) int64 {
	text := CleanHTMLTag(html)
	return int64(utf8.RuneCountInString(text) - len(blankRegexp.FindSubmatchIndex(StringToBytes(text))))
}

func WrapJSONBindErr(err error) error {
	e := validator.ValidationErrors{}
	if errors.As(err, &e) {
		return xerr.WithStatus(e, xerr.StatusBadRequest).WithMsg(trans.Translate(e))
	}
	return xerr.WithStatus(err, xerr.StatusBadRequest)
}
