// pkg/ignore.go
package pkg

import (
	"bufio"
	"io/ioutil"
	"os"
  "fmt"
	"strings"
)

const (
	ReaderIgnoreFile   = "readerout/.readerignore"
	DefaultPermissions = 0644
)

var defaultIgnorePatterns = []string{
	// Vread build
	"vread",

	// Version Control
	".git",

	// REPLIT
	".local", ".config", ".cache",

	// Node
	"node_modules",

	// Logs
	"*.log",

	// IDEs and Editors
	".vscode", ".idea", "*.iml", "*.ipr", "*.iws", "*~", "*.swp",

	// Operating System
	".DS_Store", "Thumbs.db",

	// Reader
	".readerignore", "files_structure.txt",

	// Additional Patterns
	".project-rc", "__pycache__/", "*.py[cod]", "*$py.class", "*.so",
	".Python", "build/", "develop-eggs/", "dist/", "downloads/", "eggs/", ".eggs/", "lib/", "lib64/",
	"parts/", "sdist/", "var/", "wheels/", "*.egg-info/", ".installed.cfg", "*.egg", "MANIFEST",
	"*.manifest", "*.spec", "pip-log.txt", "pip-delete-this-directory.txt", "htmlcov/", ".tox/",
	".coverage", ".coverage.*", ".cache", "nosetests.xml", "coverage.xml", "*.cover", ".hypothesis/",
	".pytest_cache/", "core.*", "*.mo", "*.pot", "*.log", "local_settings.py", "db.sqlite3", "instance/",
	".webassets-cache", ".scrapy", "docs/_build/", "target/", ".ipynb_checkpoints", ".python-version",
	"celerybeat-schedule", "*.sage.py", "/site", ".mypy_cache/",

	// Media files
	"*.mp4", "*.jpg", "*.jpeg", "*.png", "*.gif", "*.bmp", "*.tiff", "*.ico",
}

// EnsureIgnoreFileExists checks if the .readerignore file exists, and if not, creates it with default patterns.
func EnsureIgnoreFileExists() error {
	if _, err := os.Stat(ReaderIgnoreFile); os.IsNotExist(err) {
		content := getDefaultIgnorePatterns()
		return ioutil.WriteFile(ReaderIgnoreFile, []byte(content), DefaultPermissions)
	}
	return nil
}

// getDefaultIgnorePatterns returns a string containing all default ignore patterns, joined by newline characters.
func getDefaultIgnorePatterns() string {
	return strings.Join(defaultIgnorePatterns, "\n")
}

// ReadIgnorePatterns reads ignore patterns from the .readerignore file.
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
