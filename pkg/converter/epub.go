package converter

import (
	"bytes"
	"encoding/base64"

	epub "github.com/go-shiori/go-epub"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/semanticart/squire/pkg/parser"

	_ "embed"
)

//go:embed assets/default-epub-styles.css
var epubCSS []byte

func cssContent() string {
	return "data:text/plain;base64," + base64.StdEncoding.EncodeToString(epubCSS)
}

func newMarkdownConverter() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
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

func ConvertToEPUB(story parser.Story) error {
	md := newMarkdownConverter()
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

		mdBody, err := markdownToHTML(md, chapter.Body)

		if err != nil {
			return err
		}

		body += mdBody

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
