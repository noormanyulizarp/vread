// pkg/file.go
package pkg

import (
	"fmt"
	"os"
)

func CreateOutputFile() (*os.File, error) {
	return os.Create(OutputFileName)
}

func CloseFile(file *os.File) error {
	return file.Close()
}

func HandleError(message string, err error) error {
	if err != nil {
		fmt.Printf("%s: %v\n", message, err)
	}
	return err
}
