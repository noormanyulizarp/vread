package cmd

import (
	"github.com/spf13/cobra"
	"main/pkg"
	"os"
)

func Execute() {
	var (
		outputStructureOnly bool
		includePattern      string
	)

	rootCmd := createRootCmd(&outputStructureOnly, &includePattern)
	initializeFlags(rootCmd, &outputStructureOnly, &includePattern)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func createRootCmd(outputStructureOnly *bool, includePattern *string) *cobra.Command {
	return &cobra.Command{
		Use:   "vread",
		Short: "VRead is a tool for analyzing directory structures",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVRead(*outputStructureOnly, *includePattern)
		},
	}
}

func initializeFlags(rootCmd *cobra.Command, outputStructureOnly *bool, includePattern *string) {
	rootCmd.Flags().BoolVarP(outputStructureOnly, "structure", "s", false, "Output only the directory structure")
	rootCmd.Flags().StringVarP(includePattern, "include", "i", "", "Include a specific pattern regardless of .readerignore")
}

func runVRead(outputStructureOnly bool, includePattern string) error {
	rootPath := "."

	var includePatterns []string
	if includePattern != "" {
		includePatterns = append(includePatterns, includePattern)
	}

	paths, err := pkg.GetPathsToProcess(rootPath, includePatterns)
	if err != nil {
		return pkg.HandleError("Error getting paths to process", err)
	}

	outputFile, err := pkg.CreateOutputFile()
	if err != nil {
		return pkg.HandleError("Error creating output file", err)
	}
	defer pkg.CloseFile(outputFile)

	if err := pkg.PrintDirectoryTree(outputFile, rootPath, paths); err != nil {
		return err
	}
	return pkg.ProcessDirectoryStructure(outputFile, rootPath, paths)
}
