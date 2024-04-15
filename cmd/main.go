package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	flag "github.com/spf13/pflag"

	"github.com/semanticart/squire/pkg"
)

func convert(story squire.Story, format string) {
	fmt.Println("Converting to", format)

	var err error

	switch {
	case format == "html":
		err = squire.ConvertToHtml(story, true)
	case format == "html-inner":
		err = squire.ConvertToHtml(story, false)
	case format == "epub":
		err = squire.ConvertToEpub(story)
	default:
		fmt.Println("Invalid conversion format. Allowed options are 'html', 'html-inner', or 'epub'")
	}

	if err != nil {
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
