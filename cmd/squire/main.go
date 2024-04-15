package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/semanticart/squire/pkg"
)

func convert(story squire.Story, format string) {
	fmt.Println("Converting to", format)

	// TODO: make this configurable
	const fileNameBase = "output"

	var err error
	var buf bytes.Buffer
	var extension = "html"

	switch {
	case format == "html":
		buf, err = squire.ConvertToHtml(story, true)
	case format == "html-inner":
		buf, err = squire.ConvertToHtml(story, false)
	case format == "epub":
		extension = "epub"
		buf, err = squire.ConvertToEpub(story)
	default:
		fmt.Println("Invalid conversion format. Allowed options are 'html', 'html-inner', or 'epub'")
	}

	if err == nil {
		os.WriteFile(fileNameBase+"."+extension, buf.Bytes(), 0644)
	} else {
		log.Fatalf("error converting to %s: %v", format, err)
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <story.md>", os.Args[0])
	}

	convertFlag := flag.String("convert", "", "Conversion format (html, html-inner, or epub)")
	flag.Parse()

	fileName := os.Args[1]

	contents, err := os.ReadFile(fileName)

	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	story, err := squire.ParseStory(string(contents))

	if err != nil {
		for _, e := range strings.Split(err.Error(), "\n") {
			fmt.Println(fmt.Sprintf("%s:%s", fileName, e))
		}
		os.Exit(1)
	}

	convertFormat := *convertFlag

	if convertFormat != "" {
		convert(story, convertFormat)
	}
}
