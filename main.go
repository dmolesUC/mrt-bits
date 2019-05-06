package main

import (
	"fmt"
	"github.com/dmolesUC3/mrt-bits/bits"
	"os"
)

func main() {
	if err := bits.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}