package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/mariomac/gostream/stream"
)

type StringMatcher struct {
	regex *regexp.Regexp
}

func NewStringMatcher(
	words []string,
	withWordBoundary bool,
) *StringMatcher {
	sm := &StringMatcher{}

	boundedWord := stream.OfSlice(words).
		Map(func(s string) string {
			if withWordBoundary {
				return fmt.Sprintf("\\b%s\\b", s)
			}
			return s
		}).
		ToSlice()

	boundedWordRegexString := strings.Join(boundedWord, "|")

	if r, err := regexp.Compile(boundedWordRegexString); err != nil {
		return nil
	} else {
		sm.regex = r
	}

	return sm
}

func (sm *StringMatcher) HasInDocument(doc string) bool {
	return sm.regex.Match([]byte(doc))
}
