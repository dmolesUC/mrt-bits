package main

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}