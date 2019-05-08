package main

import (
	"github.com/dmolesUC3/mrt-bits/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}