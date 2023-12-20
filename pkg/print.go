// pkg/print.go
package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	SeparatorLength = 55
)

// PrintDirectoryTree prints the directory structure to the outputFile.
func PrintDirectoryTree(outputFile *os.File, rootPath string, paths []string) {
	fmt.Fprintln(outputFile, "Directory Structure:")
	sort.Strings(paths) // Sort paths to ensure correct order
	var lastDir string
	for _, path := range paths {
		dir, file := filepath.Split(path)
		// Print directory if it's different from the last one
		if dir != lastDir && dir != "" {
			fmt.Fprintln(outputFile, formatDir(dir, rootPath))
			lastDir = dir
		}
		// Print file
		if file != "" {
			fmt.Fprintln(outputFile, formatFile(dir, file, rootPath))
		}
	}
}

// formatDir formats directory string.
func formatDir(dir, rootPath string) string {
	relDir, _ := filepath.Rel(rootPath, dir)
	parts := strings.Split(relDir, string(os.PathSeparator))
	indent := strings.Repeat("│   ", len(parts)-1)
	return fmt.Sprintf("%s└── %s [Folder]", indent, parts[len(parts)-1])
}

// formatFile formats file string.
func formatFile(dir, file, rootPath string) string {
	relDir, _ := filepath.Rel(rootPath, dir)
	parts := strings.Split(relDir, string(os.PathSeparator))
	indent := strings.Repeat("│   ", len(parts))
	return fmt.Sprintf("%s└── %s", indent, file)
}

// PrintFileDetails prints details of a file or directory to the outputFile.
func PrintFileDetails(rootPath, path string, outputFile *os.File) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("error getting info for %s: %w", path, err)
	}

	printFileInfo(outputFile, rootPath, path, info)
	return nil
}

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
		if i == len(dirs)-1 {
			fmt.Fprintf(outputFile, "%s└── %s (<-)\n", indentBuilder.String(), dir)
		} else {
			fmt.Fprintf(outputFile, "%s├── %s\n", indentBuilder.String(), dir)
			indentBuilder.WriteString("│   ")
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
