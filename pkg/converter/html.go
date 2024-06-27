package converter

import (
	_ "embed"
	"fmt"

	"github.com/semanticart/squire/pkg/parser"
)

//go:embed assets/default-html-styles.css
var htmlCSS []byte

//go:embed assets/default-html-interactivity.js
var htmlJS []byte

func ConvertToHTML(rootDir string, story parser.Story, inline bool) ([]byte, error) {
	md := newMarkdownConverter(rootDir)
	html := ""

	if !inline {
		html = fmt.Sprintf("<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>%s</title>\n</head>\n<body>\n", story.Title)
	}

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

	html += fmt.Sprintf("<script>\nwindow.initialChapterID = \"%s\";\n%s</script>\n", story.Chapters[0].ID, string(htmlJS))

	html += "</div>\n"

	if !inline {
		html += "</body>\n</html>"
	}

	return []byte(html), nil
}
