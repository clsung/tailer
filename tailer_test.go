package tailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexNotWatch(t *testing.T) {
	testStrings := []string{
		"tailer.root.log.WARNING.20150211-023041",
		"root.log.WARNING.20150211-023041.1",
		"gobzip-root.log.WARNING.20150211-023041",
		".root.log.WARNING.20150211-023041.swp",
		"root.log.WARNING.20150211.gz",
	}
	for _, str := range testStrings {
		assert.Equal(t, RegexNotWatch.MatchString(str), true)
	}
	assert.Equal(t, RegexNotWatch.MatchString("root.12345.log"), false)
}
