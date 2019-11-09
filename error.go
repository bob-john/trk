package main

import (
	"fmt"
	"os"
)

func must(err error) {
	if err != nil {
		fmt.Printf("trk: %v\n", err)
		os.Exit(2)
	}
}
