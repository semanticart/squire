// Description: Entrypoint for the squire command line tool
package main

import (
	"fmt"
	"os"

	"github.com/semanticart/squire/pkg/converter"
	"github.com/semanticart/squire/pkg/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a filename")
		os.Exit(1)
	}

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

	fmt.Println(story.Title + " is valid")

	converter.ConvertToEPUB()
}
