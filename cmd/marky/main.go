package main

import (
	"fmt"
	"log"
	"os"

	"github.com/flaviodelgrosso/marky"
	"github.com/spf13/cobra"
)

func main() {
	var output string

	cmd := &cobra.Command{
		Use:   "marky <inputfile> [--output <outputfile>]",
		Short: "Convert files to markdown",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			input := args[0]

			// Check if input file exists
			if _, err := os.Stat(input); os.IsNotExist(err) {
				return fmt.Errorf("input file does not exist: %s", input)
			}

			md := marky.New()
			result, err := md.Convert(input)
			if err != nil {
				return fmt.Errorf("failed to convert file: %w", err)
			}

			if output == "console" {
				log.Println(result)
				return nil
			}

			if err := os.WriteFile(output, []byte(result), 0o644); err != nil {
				return fmt.Errorf("failed to write to output file: %w", err)
			}
			log.Printf("Content written to %s\n", output)
			return nil
		},
	}

	cmd.Flags().StringVarP(&output, "output", "o", "console", "Specify the output file path")

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
