package squire

import (
	"os"
	"testing"
)

func assertEqual(t *testing.T, got, want interface{}) {
	t.Helper()

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
}

func assertChoice(t *testing.T, choice Choice, text, chapterId string) {
	t.Helper()

	assertEqual(t, choice.Text, text)
	assertEqual(t, choice.ChapterId, chapterId)
}

func TestParseStory(t *testing.T) {
	t.Run("Parsing a valid story", func(t *testing.T) {
		contents, err := os.ReadFile("testdata/example.md")

		if err != nil {
			t.Fatal(err)
		}

		story, err := ParseStory(string(contents))

		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, story.Title, "You're probably going to die.")
		assertEqual(t, story.Author, "Jeffrey Chupp")
		assertEqual(t, story.StartChapterId, "intro")
		assertEqual(t, len(story.Chapters), 6)

		intro := story.Chapters["intro"]
		assertEqual(t, intro.Title, "Something isn't right here.")
		assertEqual(t, intro.Id, "intro")
		assertEqual(t, intro.Text, "You hear a phone ringing.")

		if len(intro.Choices) != 3 {
			t.Fatalf("got %d choices, want 3", len(intro.Choices))
		}

		assertChoice(t, intro.Choices[0], "pick up phone", "phone")
		assertChoice(t, intro.Choices[1], "do not answer", "ignore-phone")
		assertChoice(t, intro.Choices[2], "jump in a nearby lion's mouth", "lion")

		backpack := story.Chapters["backpack"]
		assertEqual(t, backpack.Title, "Going to school")
		assertEqual(t, backpack.Id, "backpack")
		assertEqual(t, backpack.Text, "You're on your way to school when a meteor lands on you, killing you instantly.")

		if len(backpack.Choices) != 1 {
			t.Fatalf("got %d choices, want 1", len(backpack.Choices))
		}

		assertChoice(t, backpack.Choices[0], "start over", "intro")
	})

	t.Run("Parsing an invalid story", func(t *testing.T) {
		contents, err := os.ReadFile("testdata/invalid.md")

		if err != nil {
			t.Fatal(err)
		}

		_, err = ParseStory(string(contents))

		if err == nil {
			t.Fatal(err)
		}

		assertEqual(t, err.Error(), `1: Missing title
1: Missing author
1: Invalid chapter title
5: Invalid choice
6: Invalid choice
9: Invalid chapter title
10: Missing chapter text
11: Invalid chapter id for choice
13: Unreachable chapter
16: Dead end
17: Unreachable chapter
22: Invalid chapter id for choice
24: Unreachable chapter
28: Invalid chapter id for choice
30: Unreachable chapter
31: Missing chapter text
31: Dead end
36: Invalid chapter id for choice`)
	})
}
