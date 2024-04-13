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

	switch {
	case format == "html":
		// TODO: Implement HTML Conversion
		fmt.Println("HTML conversion not implemented yet")
	case format == "epub":
		err := squire.ConvertToEpub(story)

		if err != nil {
			log.Fatalf("error converting to epub: %v", err)
		}
	default:
		fmt.Println("Invalid conversion format. Allowed options are 'html' or 'epub'")
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <story.md>", os.Args[0])
	}

	convertFlag := flag.String("convert", "", "Conversion format (html or epub)")
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
