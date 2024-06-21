package converter

import (
	epub "github.com/go-shiori/go-epub"

	"github.com/semanticart/squire/pkg/parser"
)

func ConvertToEPUB(story parser.Story) error {
	book, err := epub.NewEpub(story.Title)

	if err != nil {
		return err
	}

	book.SetAuthor(story.Author)

	for _, chapter := range story.Chapters {
		body := "<h1>" + chapter.Title + "</h1>"

		body += chapter.Body

		for _, choice := range chapter.Choices {
			body += "<p><a href=\"" + choice.ChapterID + ".xhtml\">" + choice.Text + "</a></p>"
		}

		_, err = book.AddSection(body, chapter.Title, chapter.ID, "")

		if err != nil {
			return err
		}
	}

	return book.Write("output.epub")
}
