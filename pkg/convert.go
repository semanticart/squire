package squire

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-shiori/go-epub"
)

//go:embed assets/default-epub-styles.css
var cssFile []byte

func cssContent() string {
	return "data:text/plain;base64," + base64.StdEncoding.EncodeToString([]byte(cssFile))
}

func header(text string) string {
	return "<h1>" + text + "</h1>\n"
}

func paragraph(text string) string {
	return "<p>" + text + "</p>\n"
}

func choiceListStart() string {
	return "<ul class=\"choices\">\n"
}

func choiceItem(text string, chapterId string, prefix string, suffix string) string {
	return fmt.Sprintf("<li><a href=\"%s%s%s\">%s</a></li>\n", prefix, chapterId, suffix, text)
}

func choiceListEnd() string {
	return "</ul>\n"
}

func ConvertToEpub(story Story) error {
	book, err := epub.NewEpub(story.Title)

	if err != nil {
		return err
	}

	book.SetAuthor(story.Author)

	cssPath, err := book.AddCSS(cssContent(), "")

	if err != nil {
		return err
	}

	for c := range len(story.Chapters) {
		for _, chapter := range story.Chapters {
			if chapter.OriginalOrder == c {
				content := header(chapter.Title)

				for _, line := range strings.Split(chapter.Text, "\n") {
					content += paragraph(line)
				}

				content += choiceListStart()

				for _, choice := range chapter.Choices {
					content += choiceItem(choice.Text, choice.ChapterId, "", ".xhtml")
				}

				content += choiceListEnd()

				_, err = book.AddSection(content, chapter.Title, chapter.Id, cssPath)

				if err != nil {
					return err
				}
			}
		}
	}

	return book.Write(story.Title + ".epub")
}
