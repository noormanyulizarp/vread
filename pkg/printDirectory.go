// pkg/printDirectory.go
package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func PrintDirectoryTree(outputFile *os.File, rootPath string, paths []string) {
	fmt.Fprintln(outputFile, "Directory Structure:")
	fmt.Fprintln(outputFile, ".")

	sort.Strings(paths)

	for i, path := range paths {
		if path == rootPath {
			continue
		}

		relPath, _ := filepath.Rel(rootPath, path)
		parts := strings.Split(relPath, string(filepath.Separator))

		for level, part := range parts {
			prefix := strings.Repeat("│   ", level)

			branchSymbol := determineBranchSymbol(level, i, parts, paths, rootPath)

			// Print the part
			if level == len(parts)-1 {
				if isDir(path) {
					fmt.Fprintf(outputFile, "%s%s%s [Folder]\n", prefix, branchSymbol, part)
				} else {
					fmt.Fprintf(outputFile, "%s%s%s\n", prefix, branchSymbol, part)
				}

				// Handle vertical line for sub-folders
				if branchSymbol == "└── " && level != 0 {
					verticalLinePrefix := strings.Repeat("│   ", level-1) + "│      "
					fmt.Fprintln(outputFile, verticalLinePrefix)
				}
			}
		}
	}

	// No additional vertical line after the last item of the root directory
}

func determineBranchSymbol(level, index int, parts []string, paths []string, rootPath string) string {
	if index == len(paths)-1 {
		return "└── "
	}

	nextPath := paths[index+1]
	nextRelPath, _ := filepath.Rel(rootPath, nextPath)
	nextParts := strings.Split(nextRelPath, string(filepath.Separator))

	if len(nextParts) <= level {
		return "└── "
	}

	return "├── "
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
