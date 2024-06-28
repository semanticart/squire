package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/semanticart/squire/pkg/converter"
	"github.com/semanticart/squire/pkg/parser"
)

var (
	format         string
	outputFileName string

	extension = map[string]string{
		"epub":        "epub",
		"html":        "html",
		"inline-html": "html",
	}
)

func convert(story parser.Story, filename, format, outputFileName string) {
	rootDir := filepath.Dir(filename)

	var err error
	var bytes []byte

	switch format {
	case "epub":
		bytes, err = converter.ConvertToEPUB(rootDir, story)
	case "html":
		bytes, err = converter.ConvertToHTML(rootDir, story, false)
	case "inline-html":
		bytes, err = converter.ConvertToHTML(rootDir, story, true)
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if outputFileName == "" {
		outputFileName = fmt.Sprintf("output.%s", extension[format])
	}

	err = os.WriteFile(outputFileName, bytes, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Written to", outputFileName)
}

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert a squire story to publishable format",
	Long: `Convert a squire story to publishable format. Supported formats include:
  - epub: A format that can be uploaded to the Kindle store
  - html: A format that can be hosted on a website
  - inline-html: A format that can be embedded in a website`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		if format != "epub" && format != "html" && format != "inline-html" {
			fmt.Println("Invalid format. Please provide a valid format (epub, html, or inline-html)")
			os.Exit(1)
		}

		content, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		story, err := parser.Parse(string(content))
		if err != nil {
			fmt.Println("Story is invalid:")
			showValidationErrors(err, fileName)
		}

		convert(story, fileName, format, outputFileName)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&format, "format", "f", "", "output format (inline-html, html, or epub)")
	_ = convertCmd.MarkFlagRequired("format")

	convertCmd.Flags().StringVarP(&outputFileName, "output", "o", "", "output filename (defaults to output + extension for format)")
}
