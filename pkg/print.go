// pkg/print.go
package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	SeparatorLength = 55
)

// PrintDirectoryTree prints the directory structure to the outputFile.
func PrintDirectoryTree(outputFile *os.File, rootPath string, paths []string) {
	fmt.Fprintln(outputFile, "Directory Structure:")
	for _, path := range paths {
		relPath, _ := filepath.Rel(rootPath, path)
		dirTree := generateDirTree(relPath)
		fmt.Fprintln(outputFile, dirTree)
	}
	fmt.Fprintln(outputFile)
}

// generateDirTree generates a string representation of the directory tree for a given path.
func generateDirTree(path string) string {
	parts := strings.Split(path, string(os.PathSeparator))
	tree := ""
	for i, part := range parts {
		if i == len(parts)-1 {
			tree += "└── " + part
		} else {
			tree += "├── " + part + "\n" + strings.Repeat("│   ", i+1)
		}
	}
	return tree
}

// PrintFileDetails prints details of a file or directory to the outputFile.
func PrintFileDetails(rootPath, path string, outputFile *os.File) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting info for %s: %w", path, err)
	}

	sizeStr := formatFileSize(info.Size())
	relPath, _ := filepath.Rel(rootPath, path)

	separator := strings.Repeat("-", SeparatorLength)
	fmt.Fprintf(outputFile, "%s\n// The size of (%s): %s\n", separator, relPath, sizeStr)
	fmt.Fprintf(outputFile, "%s\n// The file location of (%s):\n", separator, relPath)
	printPath(outputFile, relPath)
	fmt.Fprintf(outputFile, "%s\n", separator)

	if !info.IsDir() {
		content, _ := ioutil.ReadFile(path)
		fmt.Fprintf(outputFile, "\n// The content of (%s):\n", relPath)
		fmt.Fprintln(outputFile, string(content))
		fmt.Fprintf(outputFile, "%s\n", separator)
	}
	return nil
}

// printPath prints the file path in a formatted manner to the outputFile.
func printPath(outputFile *os.File, relPath string) {
	dirs := strings.Split(relPath, string(os.PathSeparator))
	indent := ""
	for i, dir := range dirs {
		if i == len(dirs)-1 {
			fmt.Fprintf(outputFile, "%s└── %s (<-)\n", indent, dir)
		} else {
			fmt.Fprintf(outputFile, "%s├── %s\n", indent, dir)
			indent += "│   "
		}
	}
}

// formatFileSize formats the size of the file for printing.
func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	return fmt.Sprintf("%dKB", size/1024)
}
