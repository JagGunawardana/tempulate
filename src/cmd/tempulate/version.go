package main

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print version number",
  Long:  `Standard semver of this binary`,
  Run: func(cmd *cobra.Command, args []string) {
	  fmt.Printf("tempulate version: %s\n", version)
  },
}