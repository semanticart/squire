// Description: Entrypoint for the squire command line tool
package main

import (
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"

	"github.com/semanticart/squire/pkg/converter"
	"github.com/semanticart/squire/pkg/parser"
)

func convert(story parser.Story, filename, format, outputFileName string) {
	rootDir := filepath.Dir(filename)

	var err error
	var bytes []byte

	switch format {
	case "epub":
		bytes, err = converter.ConvertToEPUB(rootDir, story)
	case "html":
		bytes, err = converter.ConvertToHTML(rootDir, story)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if outputFileName == "" {
		fmt.Println(string(bytes))
	} else {
		err = os.WriteFile(outputFileName, bytes, 0644)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

	format := flag.StringP("format", "f", "", "output format (html or epub)")
	outputFileName := flag.StringP("output", "o", "", "output filename")
	flag.Parse()

	filename := os.Args[1]

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	story, err := parser.Parse(string(content))
	if err != nil {
		for _, e := range err.(parser.CombinedStoryError).Errors {
			fmt.Printf("%s:%s\n", filename, e)
		}

		os.Exit(1)
	}

	if *format != "" {
		convert(story, filename, *format, *outputFileName)
	}
}
