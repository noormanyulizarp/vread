// main.go

package main

import (
	"main/cmd"
	"main/pkg"
)

func main() {
	rootCmd := cmd.NewRootCommand()

	if err := pkg.HandleError("Error executing command", rootCmd.Execute()); err != nil {
		pkg.HandleError("Error executing command", err)
	}
}
