package tailer

import (
	"testing"

	"gopkg.in/fsnotify.v1"

	"github.com/stretchr/testify/assert"
)

func TestIsUnwantEvent(t *testing.T) {
	tr := &Tailer{
		visitedFile: map[string]bool{
			"abc.123.log": true,
			"ghi.456.log": true,
		},
	}
	assert.Equal(t, tr.visitedFile["def.log.gz"], false)
	testEvents := []fsnotify.Event{
		fsnotify.Event{Name: "abc.123.log"},
		fsnotify.Event{Name: "def.log.gz"},
		fsnotify.Event{Name: "ghi.456.log"},
		fsnotify.Event{Name: "tailer.root.log.WARNING.20150211-023041"},
	}
	for _, ev := range testEvents {
		assert.Equal(t, tr.isUnwantEvent(ev), true)
	}
	// after isUnwantEvent, ths visitedFile should add def.log.gz
	assert.Equal(t, tr.visitedFile["def.log.gz"], true)
	assert.Equal(t, tr.isUnwantEvent(fsnotify.Event{Name: "root.12345.log"}), false)
}
