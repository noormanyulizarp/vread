// pkg/printDirectory.go
package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PrintDirectoryTree prints the tree structure of directories and files.
func PrintDirectoryTree(outputFile *os.File, rootPath string, paths []string) {
	fmt.Fprintln(outputFile, "Directory Structure:")
	fmt.Fprintln(outputFile, ".")

	sort.Strings(paths)

	for i, path := range paths {
		if path == rootPath {
			continue
		}

		relPath, err := filepath.Rel(rootPath, path)
		if err != nil {
			fmt.Fprintf(outputFile, "Error processing path %s: %v\n", path, err)
			continue
		}

		parts := strings.Split(relPath, string(filepath.Separator))
		printDirectoryPath(outputFile, parts, i, paths, path, rootPath)
	}
}

// printDirectoryPath handles the printing of each individual path.
func printDirectoryPath(outputFile *os.File, parts []string, index int, paths []string, currentPath, rootPath string) {
	for level, part := range parts {
		prefix := getPrefix(level)
		branchSymbol := getBranchSymbol(level, index, paths, rootPath)

		if level == len(parts)-1 {
			fmt.Fprintf(outputFile, "%s%s%s\n", prefix, branchSymbol, formatPath(part, currentPath))

			handleVerticalLineFormatting(outputFile, branchSymbol, level, index, len(paths))
		}
	}
}

// getPrefix generates the indentation prefix for each level of the path.
func getPrefix(level int) string {
	return strings.Repeat("│   ", level)
}

// getBranchSymbol determines the correct branch symbol for the tree structure.
func getBranchSymbol(level, index int, paths []string, rootPath string) string {
	if isLastItemInDirectory(level, index, paths, rootPath) {
		return "└── "
	}
	return "├── "
}

// formatPath formats the path string, adding [Folder] if it's a directory.
func formatPath(part, path string) string {
	if isDir(path) {
		return fmt.Sprintf("%s [Folder]", part)
	}
	return part
}

// handleVerticalLineFormatting manages the vertical line formatting in the tree.
func handleVerticalLineFormatting(outputFile *os.File, branchSymbol string, level, index, totalPaths int) {
	if level != 0 && branchSymbol == "└── " {
		verticalLinePrefix := strings.Repeat("│   ", level-1) + "    "
		fmt.Fprintln(outputFile, verticalLinePrefix)
	} else if level == 0 && branchSymbol == "└── " && index < totalPaths-1 {
		fmt.Fprintln(outputFile, "│   ")
	}
}

// isLastItemInDirectory determines if the current path is the last in its directory level.
func isLastItemInDirectory(level, index int, paths []string, rootPath string) bool {
	if index == len(paths)-1 {
		return true
	}

	nextPathRel, _ := filepath.Rel(rootPath, paths[index+1])
	nextParts := strings.Split(nextPathRel, string(filepath.Separator))
	return len(nextParts) <= level
}

// isDir checks if the given path is a directory.
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
