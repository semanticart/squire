package parser_test

import (
	"os"
	"testing"

	"github.com/semanticart/squire/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func assertChoice(t *testing.T, choice parser.Choice, text, chapterID string) {
	t.Helper()

	assert.Equal(t, choice.Text, text)
	assert.Equal(t, choice.ChapterID, chapterID)
}

func TestParse(t *testing.T) {
	t.Run("it parses a valid story", func(t *testing.T) {
		content, err := os.ReadFile("testdata/valid.md")
		if err != nil {
			t.Fatal(err)
		}

		story, err := parser.Parse(string(content))
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, "You're probably going to die.", story.Title)
		assert.Equal(t, "Jeffrey Chupp", story.Author)

		if len(story.Chapters) != 6 {
			t.Fatalf("expected 6 chapters, got %d", len(story.Chapters))
		}

		chapter1 := story.Chapters[0]
		assert.Equal(t, "Something isn't right here.", chapter1.Title)
		assert.Equal(t, "intro", chapter1.ID)
		assert.Equal(t, "You hear a phone ringing. _Something_ makes you suspicious of it.", chapter1.Body)

		if len(chapter1.Choices) != 3 {
			t.Fatalf("expected 3 choices, got %d", len(chapter1.Choices))
		}

		assertChoice(t, chapter1.Choices[0], "pick up phone", "phone")
		assertChoice(t, chapter1.Choices[1], "do not answer", "ignore-phone")
		assertChoice(t, chapter1.Choices[2], "jump in a nearby lion's mouth", "lion")

		lastChapter := story.Chapters[5]
		assert.Equal(t, "Going to school", lastChapter.Title)
		assert.Equal(t, "backpack", lastChapter.ID)
		assert.Equal(t, "You're on your way to school when a meteor lands on you. You gain super powers and institute world peace.\n\nYou win.", lastChapter.Body)

		if len(lastChapter.Choices) != 0 {
			t.Fatalf("expected 0 choices, got %d", len(lastChapter.Choices))
		}
	})

	t.Run("it returns an error when the story is invalid", func(t *testing.T) {
		content, err := os.ReadFile("testdata/invalid.md")
		if err != nil {
			t.Fatal(err)
		}

		_, err = parser.Parse(string(content))

		if err == nil {
			t.Fatal("expected an error, got nil")
		}

		expected := `1: Missing title
1: Missing author
1: Missing chapter id
5: Missing choice id
6: Missing choice text
9: Missing chapter title
9: Unreachable chapter
10: Missing chapter text
11: Invalid chapter id for choice
16: Dead end
22: Invalid chapter id for choice
24: Unreachable chapter
28: Invalid chapter id for choice
30: Unreachable chapter
31: Missing chapter text
31: Dead end`

		assert.Equal(t, expected, err.Error())
	})
}
