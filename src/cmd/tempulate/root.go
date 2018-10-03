package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tempulate",
	Short: "Tempulate takes input parameters (YAML or JSON) and templates out a file using Golang templates and the input params",
	Long: `A useful tool (devops, dev automation)
         to generate files from templates and parameter inputs`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("You used params: %s\nTemplate: %s\n", params, templateFile)
	},
}

var params string
var templateFile string

func init() {
	rootCmd.Flags().StringVarP(&params, "params", "p", "", "params files (YAML/JSON) provided as a list")
	rootCmd.MarkFlagRequired("params")
	rootCmd.Flags().StringVarP(&templateFile, "template", "t", "", "File to template out")
	rootCmd.MarkFlagRequired("template")
}
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
