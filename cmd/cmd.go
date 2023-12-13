package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"vread/pkg"
)

// Execute initializes and runs the root CLI command.
func Execute() {
	var (
		outputStructureOnly bool
		includePattern      string
	)

	rootCmd := createRootCmd(&outputStructureOnly, &includePattern)
	if err := rootCmd.Execute(); err != nil {
		// Error handling: log the error and exit
		os.Exit(1)
	}
}

// createRootCmd creates the root Cobra command with necessary flags and options.
func createRootCmd(outputStructureOnly *bool, includePattern *string) *cobra.Command {
	return &cobra.Command{
		Use:   "vread",
		Short: "VRead is a tool for analyzing directory structures",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVRead(*outputStructureOnly, *includePattern)
		},
	}
}

// initializeFlags initializes the flags for the root command.
func initializeFlags(rootCmd *cobra.Command, outputStructureOnly *bool, includePattern *string) {
	rootCmd.Flags().BoolVarP(outputStructureOnly, "structure", "s", false, "Output only the directory structure")
	rootCmd.Flags().StringVarP(includePattern, "include", "i", "", "Include a specific pattern regardless of .readerignore")
}

// runVRead handles the core logic for the vread command.
func runVRead(outputStructureOnly bool, includePattern string) error {
	rootPath := "."

	if err := pkg.SetupOutputEnvironment(outputStructureOnly, includePattern); err != nil {
		return err
	}

	paths, err := pkg.GetPathsToProcess(rootPath, includePattern)
	if err != nil {
		return err
	}

	outputFile, err := pkg.CreateOutputFile()
	if err != nil {
		return err
	}
	defer pkg.CloseFile(outputFile)

	if outputStructureOnly {
		return pkg.PrintDirectoryTree(outputFile, rootPath, paths)
	}
	return pkg.ProcessDirectoryStructure(outputFile, rootPath, paths)
}

// init function to initialize the command with flags.
func init() {
	rootCmd := createRootCmd(&outputStructureOnly, &includePattern)
	initializeFlags(rootCmd)
}
