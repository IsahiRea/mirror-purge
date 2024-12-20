package main

import (
	"fmt"
	"log"
	"os"
)

func main() {

	var dir string

	//TODO: Figure out Command Line Arguments

	// ifHash := os.Args[0]
	// ifOutput := os.Args[1]
	// ifDelete := os.Args[2]

	argsOptions := os.Args[1:]

	if len(argsOptions) < 2 {
		dir = "."
	} else {
		dir = os.Args[2]
	}

	dirScan := scanDir(dir)
	fmt.Println(dirScan)

	copies, err := findDuplicates(dirScan)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(copies)
}
