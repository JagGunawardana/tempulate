package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"os"

	"github.com/JagGunawardana/tempulate/src/munge"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tempulate",
	Short: "Tempulate takes input parameters (YAML or JSON) and templates out a file using Golang templates and the input params",
	Long: `A useful tool (devops, dev automation)
         to generate files from templates and parameter inputs`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		files := strings.Split(params, ",")
		paramFiles := strings.Split(params, ",")
		files = append(files, templateFile)
		if err := checkFiles(files); err != nil {
			log.Fatal(err)
		}
		var f *os.File
		var err error
		if outputFile == "" {
			outputFile = "STDOUT"
			f = os.Stdout
		} else {
			f, err = os.Create(outputFile)
			defer f.Close()
			if err != nil {
				log.Fatal("Failed to open output file")
			}
		}
		if !quiet {
			fmt.Printf("Template file: %s\nParameter file(s): %s\nOutput file: %s\n", templateFile, params, outputFile)
		}
		contents, err := ioutil.ReadFile(templateFile)
		if err != nil {
			log.Fatal(err)
		}
		output, err := munge.MungeFile(string(contents), paramFiles)
		if err != nil {
			log.Fatal(err)
		}
		n, err := f.Write([]byte(output))
		if err != nil {
			log.Fatal(err)
		}
		if !quiet {
			fmt.Printf("SUCCESS: wrote %d bytes", n)
		}
	},
}

// checkFiles examines a list of files to check that they all exist
func checkFiles(toCheck []string) error {
	issues := []string{}
	for _, f := range toCheck {
		if _, err := os.Stat(f); os.IsNotExist(err) {
			issues = append(issues, f)
		}
	}
	if len(issues) != 0 {
		return fmt.Errorf("Failed to find files: %s", strings.Join(issues, ","))
	}
	return nil
}

var params string
var templateFile string
var outputFile string
var quiet bool

func init() {
	rootCmd.Flags().StringVarP(&params, "params", "p", "", "params files (YAML/JSON) provided as a list")
	rootCmd.MarkFlagRequired("params")
	rootCmd.Flags().StringVarP(&templateFile, "template", "t", "", "File to template out")
	rootCmd.MarkFlagRequired("template")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Templated output file")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Turn on quiet mode (just output file - no other info)")
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
