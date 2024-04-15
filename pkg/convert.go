package squire

import (
	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/go-shiori/go-epub"
)

//go:embed assets/default-epub-styles.css
var epubCSS []byte

//go:embed assets/default-html-interactivity.js
var htmlJS []byte

//go:embed assets/default-html-styles.css
var htmlCSS []byte

func cssContent() string {
	return "data:text/plain;base64," + base64.StdEncoding.EncodeToString([]byte(epubCSS))
}

func header(id string, text string) string {
	return fmt.Sprintf("<h1 id=\"%s\">%s</h1>\n", id, text)
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

func sortedChapters(chapters map[string]Chapter) []Chapter {
	sorted := make([]Chapter, 0, len(chapters))

	for c := range len(chapters) {
		for _, chapter := range chapters {
			if chapter.OriginalOrder == c {
				sorted = append(sorted, chapter)
			}
		}
	}
	return sorted
}

func ConvertToHtml(story Story, full bool) (bytes.Buffer, error) {
	html := ""

	if full {
		html += "<!DOCTYPE html>\n<html>\n<head>\n<meta charset=\"utf-8\">\n<title>" + story.Title + "</title>\n</head>\n<body>\n"
	}

	html += "<div id='story'>\n"
	html += "<style>" + string(htmlCSS) + "</style>\n"

	for _, chapter := range sortedChapters(story.Chapters) {
		html += header(chapter.Id, chapter.Title)

		for _, line := range strings.Split(chapter.Text, "\n") {
			html += paragraph(line)
		}

		if len(chapter.Choices) > 0 {
			html += choiceListStart()

			for _, choice := range chapter.Choices {
				html += choiceItem(choice.Text, choice.ChapterId, "#", "")
			}

			html += choiceListEnd()
		}
	}

	html += "<script>\n" + string(htmlJS) + "</script>\n"

	html += "</div>"

	if full {
		html += "</body>\n</html>"
	}

	buffer := bytes.NewBufferString(html)

	return *buffer, nil
}

func ConvertToEpub(story Story) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	book, err := epub.NewEpub(story.Title)

	if err != nil {
		return buffer, err
	}

	book.SetAuthor(story.Author)

	cssPath, err := book.AddCSS(cssContent(), "")

	if err != nil {
		return buffer, err
	}

	for _, chapter := range sortedChapters(story.Chapters) {
		content := header(chapter.Id, chapter.Title)

		for _, line := range strings.Split(chapter.Text, "\n") {
			content += paragraph(line)
		}

		if len(chapter.Choices) > 0 {
			content += choiceListStart()

			for _, choice := range chapter.Choices {
				content += choiceItem(choice.Text, choice.ChapterId, "", ".xhtml")
			}

			content += choiceListEnd()
		}

		_, err = book.AddSection(content, chapter.Title, chapter.Id, cssPath)

		if err != nil {
			return buffer, err
		}
	}

	_, err = book.WriteTo(&buffer)

	if err != nil {
		return buffer, err
	}

	return buffer, nil
}
