package main

import (
	"fmt"

	"github.com/cafebazaar/go-boilerplate/pkg/postview"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version info",
	Long:  `All softwares have versions. This is aggregator`,
	Run: func(cmd *cobra.Command, args []string) {
		printVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion() {
	fmt.Printf("%-18s %-18s Commit:%s                  (%s)\n", postview.Title, postview.Version,
		postview.Commit, postview.BuildTime)
}
