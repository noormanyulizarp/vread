// pkg/printDirectory.go
package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PrintDirectoryTree prints the directory structure to the outputFile.
func PrintDirectoryTree(outputFile *os.File, rootPath string, paths []string) {
	fmt.Fprintln(outputFile, "Directory Structure:")
	fmt.Fprintln(outputFile, ".")

	// Sort paths to ensure correct order
	sort.Strings(paths)

	for i, path := range paths {
		if path == rootPath {
			continue // Skip the root directory itself
		}

		relPath, _ := filepath.Rel(rootPath, path)
		parts := strings.Split(relPath, string(filepath.Separator))

		for level, part := range parts {
			prefix := ""

			// Create prefix for the current level
			for i := 0; i < level; i++ {
				prefix += "│   "
			}

			// Print the part
			if level == len(parts)-1 {
				if isDir(path) {
					fmt.Fprintf(outputFile, "%s├── %s [Folder]\n", prefix, part)
				} else {
					fmt.Fprintf(outputFile, "%s└── %s\n", prefix, part)
				}

				// Determine if this is the last item in the folder
				if i < len(paths)-1 {
					nextPath := paths[i+1]
					nextRelPath, _ := filepath.Rel(rootPath, nextPath)
					nextParts := strings.Split(nextRelPath, string(filepath.Separator))

					if len(nextParts) <= level {
						// Print extra line with vertical bar if it's the last item
						fmt.Fprintln(outputFile, prefix)
					}
				} else {
					// Always print the extra line for the very last item
					fmt.Fprintln(outputFile, prefix)
				}
			}
		}
	}
}

// isDir checks if a given path is a directory
func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
