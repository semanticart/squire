package converter

import (
	"bytes"
	"fmt"

	img64 "github.com/tenkoh/goldmark-img64"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/semanticart/squire/pkg/parser"
)

func newMarkdownConverter(rootDir string) goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM,
			img64.Img64,
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			img64.WithParentPath(rootDir),
		),
	)
}

func markdownToHTML(md goldmark.Markdown, markdown string) (string, error) {
	var buf bytes.Buffer

	err := md.Convert([]byte(markdown), &buf)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func title(chapter parser.Chapter) string {
	return fmt.Sprintf("<h1 id=\"%s\">%s</h1>", chapter.ID, chapter.Title)
}

func choices(chapter parser.Chapter, prefix string, suffix string) string {
	choiceHTML := ""
	if len(chapter.Choices) > 0 {
		choiceHTML += "<ul class=\"choices\">"

		for _, choice := range chapter.Choices {
			choiceHTML += fmt.Sprintf("<li><a href=\"%s%s%s\">%s</a></li>", prefix, choice.ChapterID, suffix, choice.Text)
		}

		choiceHTML += "</ul>"
	}

	return choiceHTML
}
