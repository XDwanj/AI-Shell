package main

import (
	"fmt"
	"os"

	"AI-Shell/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行错误: %v\n", err)
		os.Exit(1)
	}
}
