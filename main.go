package main

import (
	"fmt"
	"goadmin/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
