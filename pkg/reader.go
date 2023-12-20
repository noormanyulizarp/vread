package pkg

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const (
	OutputFolder   = "readerout"
	OutputFileName = OutputFolder + "/files_structure.txt"
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

// ProcessDirectoryStructure processes each path in the directory structure and prints details.
func ProcessDirectoryStructure(outputFile *os.File, rootPath string, paths []string) {
	for _, path := range paths {
		if err := PrintFileDetails(rootPath, path, outputFile); err != nil {
			fmt.Printf("Error printing file details for %s: %v\n", path, err)
			continue
		}
	}
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
