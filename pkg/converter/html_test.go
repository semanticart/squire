package converter_test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/semanticart/squire/pkg/converter"
	"github.com/semanticart/squire/pkg/parser"
)

func TestConvertStory(t *testing.T) {
	t.Run("Converting a story to html", func(t *testing.T) {
		contents, err := os.ReadFile("../parser/testdata/valid.md")
		if err != nil {
			t.Fatal(err)
		}

		story, err := parser.Parse(string(contents))
		if err != nil {
			t.Fatal(err)
		}

		buf, err := converter.ConvertToHTML("testdata/", story, false)
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile("testdata/example.converted.html")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, strings.TrimSpace(string(buf)), strings.TrimSpace(string(expected)))
	})

	t.Run("Converting a story to inner-html", func(t *testing.T) {
		contents, err := os.ReadFile("../parser/testdata/valid.md")
		if err != nil {
			t.Fatal(err)
		}

		story, err := parser.Parse(string(contents))
		if err != nil {
			t.Fatal(err)
		}

		buf, err := converter.ConvertToHTML("testdata/", story, true)
		if err != nil {
			t.Fatal(err)
		}

		expected, err := os.ReadFile("testdata/example.converted.inner.html")
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, strings.TrimSpace(string(buf)), strings.TrimSpace(string(expected)))
	})
}
