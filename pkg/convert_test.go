package squire

import (
	"os"
	"strings"
	"testing"
)

func TestConvertStory(t *testing.T) {
	t.Cleanup(func() {
		os.Remove("You're probably going to die..html")
	})

	t.Run("Converting a story to html", func(t *testing.T) {
		contents, err := os.ReadFile("testdata/example.md")

		if err != nil {
			t.Fatal(err)
		}

		story, err := ParseStory(string(contents))

		if err != nil {
			t.Fatal(err)
		}

		buf, err := ConvertToHtml(story, true)

		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile("testdata/example.converted.html")

		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, strings.TrimSpace(buf.String()), strings.TrimSpace(string(expected)))
	})

	t.Run("Converting a story to inner-html", func(t *testing.T) {
		contents, err := os.ReadFile("testdata/example.md")

		if err != nil {
			t.Fatal(err)
		}

		story, err := ParseStory(string(contents))

		if err != nil {
			t.Fatal(err)
		}

		buf, err := ConvertToHtml(story, false)

		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile("testdata/example.converted.inner.html")

		if err != nil {
			t.Fatal(err)
		}

		assertEqual(t, strings.TrimSpace(buf.String()), strings.TrimSpace(string(expected)))
	})
}
