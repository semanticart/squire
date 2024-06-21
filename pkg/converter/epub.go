package converter

import (
	"encoding/base64"

	epub "github.com/go-shiori/go-epub"

	"github.com/semanticart/squire/pkg/parser"

	_ "embed"
)

//go:embed assets/default-epub-styles.css
var epubCSS []byte

func cssContent() string {
	return "data:text/plain;base64," + base64.StdEncoding.EncodeToString(epubCSS)
}

func ConvertToEPUB(story parser.Story) error {
	book, err := epub.NewEpub(story.Title)

	if err != nil {
		return err
	}

	book.SetAuthor(story.Author)

	cssPath, err := book.AddCSS(cssContent(), "")

	if err != nil {
		return err
	}

	for _, chapter := range story.Chapters {
		body := "<h1>" + chapter.Title + "</h1>"

		body += chapter.Body

		if len(chapter.Choices) > 0 {
			body += "<ul class=\"choices\">"

			for _, choice := range chapter.Choices {
				body += "<li><a href=\"" + choice.ChapterID + ".xhtml\">" + choice.Text + "</a></li>"
			}

			body += "</ul>"
		}

		_, err = book.AddSection(body, chapter.Title, chapter.ID, cssPath)

		if err != nil {
			return err
		}
	}

	return book.Write("output.epub")
}
