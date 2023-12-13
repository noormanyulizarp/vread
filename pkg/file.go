package pkg

import (
	"fmt"
	"os"
)

func CreateOutputFile() (*os.File, error) {
	return os.Create(OutputFileName)
}

func CloseFile(file *os.File) {
	if err := file.Close(); err != nil {
		fmt.Printf("Error closing file: %v\n", err)
	}
}
