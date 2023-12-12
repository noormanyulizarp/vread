package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"sort"
	"vread/pkg"
)

var structureFlag bool
var includePattern string

var rootCmd = &cobra.Command{
	Use:   "vread",
	Short: "VRead is a tool for analyzing directory structures",
	Run: func(cmd *cobra.Command, args []string) {
		runVRead()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(&structureFlag, "structure", "s", false, "Output only the directory structure")
	rootCmd.Flags().StringVarP(&includePattern, "include", "i", "", "Include a specific pattern regardless of .readerignore")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pkg.HandleError("Error executing command", err)
	}
}

func runVRead() {
	rootPath := "."

	if err := os.MkdirAll(pkg.OutputFolder, os.ModePerm); err != nil {
		pkg.HandleError("Error creating output folder", err)
	}

	if err := pkg.EnsureIgnoreFileExists(); err != nil {
		pkg.HandleError("Error ensuring ignore file exists", err)
	}

	excludePatterns := pkg.ReadIgnorePatterns()
	if includePattern != "" {
		excludePatterns = append(excludePatterns, includePattern)
	}

	paths, err := pkg.GetPathsToProcess(rootPath, excludePatterns)
	if err != nil {
		pkg.HandleError("Error processing paths", err)
	}

	outputFile, err := pkg.CreateOutputFile()
	if err != nil {
		pkg.HandleError("Error creating output file", err)
	}
	defer pkg.CloseFile(outputFile)

	if structureFlag {
		pkg.PrintDirectoryTree(outputFile, rootPath, paths)
	} else {
		pkg.PrintDirectoryTree(outputFile, rootPath, paths)
		pkg.ProcessDirectoryStructure(outputFile, rootPath, paths)
	}
}
