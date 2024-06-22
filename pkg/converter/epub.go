package converter

import (
	"bytes"
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

func ConvertToEPUB(rootDir string, story parser.Story) ([]byte, error) {
	md := newMarkdownConverter(rootDir)
	book, err := epub.NewEpub(story.Title)

	if err != nil {
		return []byte{}, err
	}

	book.SetAuthor(story.Author)

	cssPath, err := book.AddCSS(cssContent(), "")

	if err != nil {
		return []byte{}, err
	}

	for _, chapter := range story.Chapters {
		body := title(chapter)

		var mdBody string
		mdBody, err = markdownToHTML(md, chapter.Body)

		if err != nil {
			return []byte{}, err
		}

		body += mdBody
		body += choices(chapter, "", ".xhtml")

		_, err = book.AddSection(body, chapter.Title, chapter.ID, cssPath)

		if err != nil {
			return []byte{}, err
		}
	}

	var buf bytes.Buffer

	_, err = book.WriteTo(&buf)

	return buf.Bytes(), err
}
