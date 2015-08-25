package main

import (
	"log"
	"os"
	"stash.tsrapplabs.com/ut/pbc"
)

func main() {

	root := "."
	if len(os.Args) > 1 {
		root = os.Args[1]
	}

	err := pbc.Compile(root)

	if err != nil {
		log.Fatal(err)
	}
}
