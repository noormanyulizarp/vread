// pkg/reader.go
package pkg

import (
	"bufio"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	SeparatorLength = 55
	OutputFolder    = "readerout"
	OutputFileName  = OutputFolder + "/files_structure.txt"
)

// GetPathsToProcess walks through the directory tree starting from rootPath, respecting ignore patterns.
func GetPathsToProcess(rootPath string, excludePatterns []string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %q: %w", path, err)
		}
		if !shouldExclude(path, excludePatterns) {
			paths = append(paths, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(paths)
	return paths, nil
}

// ReadIgnorePatterns reads ignore patterns from .readerignore file.
func ReadIgnorePatterns() []string {
	var patterns []string
	file, err := os.Open(ReaderIgnoreFile)
	if err != nil {
		HandleError(fmt.Sprintf("Error opening %s", ReaderIgnoreFile), err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}
	return patterns
}

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

// ProcessDirectoryStructure processes each path in the directory structure and prints details.
func ProcessDirectoryStructure(outputFile *os.File, rootPath string, paths []string) {
	for _, path := range paths {
		if err := PrintFileDetails(rootPath, path, outputFile); err != nil {
			fmt.Printf("Error printing file details for %s: %v\n", path, err)
			continue
		}
	}
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

// formatFileSize formats the size of the file for printing.
func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	return fmt.Sprintf("%dKB", size/1024)
}

// shouldExclude determines if a given path should be excluded based on provided patterns.
func shouldExclude(path string, patterns []string) bool {
	relPath, err := filepath.Rel(".", path)
	if err != nil {
		return false
	}

	for _, pattern := range patterns {
		matched, err := filepath.Match(pattern, relPath)
		if err != nil {
			continue
		}

		for _, part := range strings.Split(relPath, string(os.PathSeparator)) {
			dirMatched, _ := filepath.Match(pattern, part)
			if dirMatched {
				return true
			}
		}

		if matched {
			return true
		}
	}
	return false
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
