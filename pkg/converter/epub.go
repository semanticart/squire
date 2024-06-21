package converter

import (
	epub "github.com/go-shiori/go-epub"
)

func ConvertToEPUB() {
	book, err := epub.NewEpub("My Title")

	if err != nil {
		panic(err)
	}

	book.SetAuthor("Jeffrey Chupp")

	book.AddSection(`

	<p>Hello, <b>world</b></p>

	<a href="second-section.xhtml">Some choice goes here</a>
	`, "First Section", "first-section", "")
	book.AddSection(`<p>Another text</p>

	<a href="first-section.xhtml">Start over</a>
	`, "Second Section", "second-section", "")

	book.Write("output.epub")
}
