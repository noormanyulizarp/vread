package pkg

import (
	"fmt"
	"os"
)

func CreateOutputFile() (*os.File, error) {
	return os.Create(OutputFileName)
}

func CloseFile(file *os.File) {
	_ = file.Close()
}

func HandleError(message string, err error) {
	fmt.Printf("%s: %v\n", message, err)
	os.Exit(1)
}
