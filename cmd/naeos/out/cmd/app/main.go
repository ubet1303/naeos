package main

import (
	"fmt"
	"os"

	"github.com/example/sample-specification/internal/core/config"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("sample-specification started on port %d\n", cfg.Port)
}
