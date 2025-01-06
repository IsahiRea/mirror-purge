package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var dir string
	var useHash bool
	var outputFile string
	var deleteDuplicates bool

	// Define command line arguments
	flag.BoolVar(&useHash, "hash", false, "Use hash comparisons")
	flag.BoolVar(&useHash, "h", false, "Use hash comparisons (shorthand)")
	flag.StringVar(&outputFile, "output", "", "Specify an output file for the results")
	flag.StringVar(&outputFile, "o", "", "Specify an output file for the results (shorthand)")
	flag.BoolVar(&deleteDuplicates, "delete", false, "Prompt to delete duplicates")
	flag.BoolVar(&deleteDuplicates, "d", false, "Prompt to delete duplicates (shorthand)")

	// Parse command line arguments
	flag.Parse()

	// Get the directory from the remaining arguments
	args := flag.Args()
	if len(args) < 1 {
		dir = "."
	} else {
		dir = args[0]
	}

	//TODO: Implement different logic for hash comparisons
	// Default for now is to use hash comparisons

	// Scan the directory
	dirScan := scanDir(dir)
	fmt.Println(dirScan)

	copies, err := findDuplicates(dirScan)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(copies)

	// Handle output file if specified
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(fmt.Sprintf("%v", copies)), 0644)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
		}
	}

	// Handle delete duplicates if specified
	if deleteDuplicates {
		//TODO: Implement delete logic here
		fmt.Println("Delete duplicates option selected")

		for _, duplicates := range copies {
			for _, file := range duplicates[1:] {
				err := os.Remove(file)
				if err != nil {
					fmt.Printf("Error deleting file %s: %v\n", file, err)
				} else {
					fmt.Printf("Deleted duplicate file %s\n", file)
				}
			}
		}
	}

}
