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

// PrintFileDetails prints details of a file or directory to the outputFile.
func PrintFileDetails(rootPath, path string, outputFile *os.File) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting info for %s: %w", path, err)
	}

	printFileInfo(outputFile, rootPath, path, info)
	return nil
}

// printFileInfo prints detailed information of a file or directory.
func printFileInfo(outputFile *os.File, rootPath, path string, info os.FileInfo) {
	relPath, _ := filepath.Rel(rootPath, path)
	sizeStr := formatFileSize(info.Size())
	separator := strings.Repeat("-", SeparatorLength)

	fmt.Fprintf(outputFile, "%s\n// The size of (%s): %s\n", separator, relPath, sizeStr)
	fmt.Fprintf(outputFile, "%s\n// The file location of (%s):\n", separator, relPath)
	printPath(outputFile, relPath)
	fmt.Fprintf(outputFile, "%s\n", separator)

	if !info.IsDir() {
		content, err := ioutil.ReadFile(path)
		if err == nil {
			fmt.Fprintf(outputFile, "\n// The content of (%s):\n", relPath)
			fmt.Fprintln(outputFile, string(content))
			fmt.Fprintf(outputFile, "%s\n", separator)
		}
	}
}

// printPath prints the file path in a formatted manner to the outputFile.
func printPath(outputFile *os.File, relPath string) {
	dirs := strings.Split(relPath, string(os.PathSeparator))
	var indentBuilder strings.Builder

	for i, dir := range dirs {
		if dir == "" {
			continue // Skip empty directory names
		}

		if i > 0 {
			indentBuilder.WriteString("│   ")
		}

		if i == len(dirs)-1 {
			fmt.Fprintf(outputFile, "%s└── %s (<-)\n", indentBuilder.String(), dir)
		} else {
			fmt.Fprintf(outputFile, "%s├── %s\n", indentBuilder.String(), dir)
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
