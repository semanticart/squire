package converter

import (
	"bytes"

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
	return "<h1 id=\"" + chapter.ID + "\">" + chapter.Title + "</h1>"
}

func choices(chapter parser.Chapter, prefix string, suffix string) string {
	choiceHTML := ""
	if len(chapter.Choices) > 0 {
		choiceHTML += "<ul class=\"choices\">"

		for _, choice := range chapter.Choices {
			choiceHTML += "<li><a href=\"" + prefix + choice.ChapterID + suffix + "\">" + choice.Text + "</a></li>"
		}

		choiceHTML += "</ul>"
	}

	return choiceHTML
}
