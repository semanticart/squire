package converter

import (
	"github.com/semanticart/squire/pkg/parser"

	_ "embed"
)

//go:embed assets/default-html-styles.css
var htmlCSS []byte

//go:embed assets/default-html-interactivity.js
var htmlJS []byte

func ConvertToHTML(rootDir string, story parser.Story) ([]byte, error) {
	md := newMarkdownConverter(rootDir)
	html := "<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>" + story.Title + "</title>\n</head>\n<body>\n"

	html += "<div id=\"story\">\n"
	html += "<style>\n" + string(htmlCSS) + "\n</style>\n"

	for _, chapter := range story.Chapters {
		chapterHTML := title(chapter)

		mdChapterText, err := markdownToHTML(md, chapter.Body)

		if err != nil {
			return []byte{}, err
		}

		chapterHTML += mdChapterText
		chapterHTML += choices(chapter, "#", "")

		html += chapterHTML
	}

	html += "<script>\nwindow.initialChapterID = \"" +
		story.Chapters[0].ID + "\";\n" + string(htmlJS) + "</script>\n"
	html += "</div>\n"

	html += "</body>\n</html>"

	return []byte(html), nil
}
