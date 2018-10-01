package tempulate

import (
	"log"

	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
  Use:   "tempulate",
  Short: "Tempulate takes input parameters (YAML or JSON) and templates out a file using Golang templates and the input params",
  Long: `A useful tool (devops, dev automation)
         to generate files from templates and parameter inputs`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}