package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/semanticart/squire/pkg"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("usage: %s <story.md>", os.Args[0])
	}

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

	fmt.Println(story)
}
