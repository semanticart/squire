package internal_test

import (
	"testing"

	"github.com/semanticart/squire/internal"
	"github.com/stretchr/testify/assert"
)

func TestParseChapterTitleAndID(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		title   string
		id      string
		deadEnd bool
	}{
		{
			name:    "it parses a valid chapter title",
			line:    "# Chapter Title {#chapter-title}",
			title:   "Chapter Title",
			id:      "chapter-title",
			deadEnd: false,
		},

		{
			name:    "it parses a valid chapter title with an intentional dead end",
			line:    "# Chapter Title {#chapter-title} !!",
			title:   "Chapter Title",
			id:      "chapter-title",
			deadEnd: true,
		},

		{
			name:    "it parses an invalid chapter title with no text",
			line:    "# {#chapter-title}",
			title:   "",
			id:      "chapter-title",
			deadEnd: false,
		},

		{
			name:    "it parses an invalid chapter title with no text and an intentional dead end",
			line:    "# {#chapter-title} !!",
			title:   "",
			id:      "chapter-title",
			deadEnd: true,
		},

		{
			name:    "it parses an invalid chapter title with no id",
			line:    "# Chapter Title",
			title:   "Chapter Title",
			id:      "",
			deadEnd: false,
		},

		{
			name:    "it parses an invalid chapter title with no id and an intentional dead end",
			line:    "# Chapter Title !!",
			title:   "Chapter Title",
			id:      "",
			deadEnd: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, id, intentionalDeadEnd := internal.ParseChapterTitleAndID(tt.line)

			assert.Equal(t, tt.title, title)
			assert.Equal(t, tt.id, id)
			assert.Equal(t, tt.deadEnd, intentionalDeadEnd)
		})
	}
}
