package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/semanticart/squire/pkg/parser"
)

func showValidationErrors(err error, fileName string) {
	for _, e := range err.(parser.CombinedStoryError).Errors {
		fmt.Printf("%s:%s\n", fileName, e)
	}
}

// validateCmd represents the validate command
var validateCmd = &cobra.Command{
	Use:   "validate [filename.md]",
	Short: "Validate a squire story",
	Long:  `Validate a squire story and return detailed information about any errors`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileName := args[0]

		content, err := os.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		_, err = parser.Parse(string(content))
		if err != nil {
			showValidationErrors(err, fileName)
			os.Exit(1)
		}

		fmt.Println("Story is valid")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
