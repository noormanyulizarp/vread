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
	SeparatorLength    = 55
	OutputFolder       = "readerout"
	ReaderIgnoreFile   = OutputFolder + "/.readerignore"
	OutputFileName     = OutputFolder + "/files_structure.txt"
	DefaultPermissions = 0644
)

var defaultIgnorePatterns = []string{
	"# Vread build", "vread", "",
	"# Version Control", ".git", "",
	"# REPLIT", ".local", ".config", ".cache", "",
	"# Node", "node_modules", "",
	"# Logs", "*.log", "",
	"# IDEs and Editors", ".vscode", ".idea", "*.iml", "*.ipr", "*.iws", "*~", "*.swp", "",
	"# Operating System", ".DS_Store", "Thumbs.db", "",
	"# Reader", ".readerignore", "files_structure.txt", "",
	"# Additional Patterns", ".project-rc", "__pycache__/", "*.py[cod]", "*$py.class", "*.so",
	".Python", "build/", "develop-eggs/", "dist/", "downloads/", "eggs/", ".eggs/", "lib/", "lib64/",
	"parts/", "sdist/", "var/", "wheels/", "*.egg-info/", ".installed.cfg", "*.egg", "MANIFEST",
	"*.manifest", "*.spec", "pip-log.txt", "pip-delete-this-directory.txt", "htmlcov/", ".tox/",
	".coverage", ".coverage.*", ".cache", "nosetests.xml", "coverage.xml", "*.cover", ".hypothesis/",
	".pytest_cache/", "core.*", "*.mo", "*.pot", "*.log", "local_settings.py", "db.sqlite3", "instance/",
	".webassets-cache", ".scrapy", "docs/_build/", "target/", ".ipynb_checkpoints", ".python-version",
	"celerybeat-schedule", "*.sage.py", "/site", ".mypy_cache/", "",
	"# Media files", "*.mp4", "*.jpg", "*.jpeg", "*.png", "*.gif", "*.bmp", "*.tiff", "*.ico",
}

func EnsureIgnoreFileExists() error {
	if _, err := os.Stat(ReaderIgnoreFile); os.IsNotExist(err) {
		content := getDefaultIgnorePatterns()
		return ioutil.WriteFile(ReaderIgnoreFile, []byte(content), DefaultPermissions)
	}
	return nil
}

func getDefaultIgnorePatterns() string {
	return strings.Join(defaultIgnorePatterns, "\n")
}

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

func ReadIgnorePatterns() []string {
	var patterns []string
	file, err := os.Open(ReaderIgnoreFile)
	if err != nil {
		HandleError(fmt.Sprintf("Error opening %s", ReaderIgnoreFile), err)
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

func PrintDirectoryTree(outputFile *os.File, rootPath string, paths []string) {
	fmt.Fprintln(outputFile, "Directory Structure:")
	for _, path := range paths {
		relPath, _ := filepath.Rel(rootPath, path)
		dirTree := generateDirTree(relPath)
		fmt.Fprintln(outputFile, dirTree)
	}
	fmt.Fprintln(outputFile)
}

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

func ProcessDirectoryStructure(outputFile *os.File, rootPath string, paths []string) {
	for _, path := range paths {
		if err := PrintFileDetails(rootPath, path, outputFile); err != nil {
			fmt.Printf("Error printing file details for %s: %v\n", path, err)
			continue
		}
	}
}

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

func formatFileSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	return fmt.Sprintf("%dKB", size/1024)
}

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

func HandleError(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
	os.Exit(1)
}
