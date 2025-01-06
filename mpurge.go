package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var (
	dir              string
	useHash          string
	outputFile       string
	deleteDuplicates bool
)

// TODO: Find duplicates using similarities in file size and hash
// TODO: Implement command line argument for traversing subdirectories
func main() {

	// Define command line arguments
	flag.StringVar(&useHash, "hash", "md5", "Hash algorithm to use (md5, sha256)")
	flag.StringVar(&useHash, "h", "md5", "Hash algorithm to use (shorthand)")
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

	// Scan the directory
	dirScan := scanDir(dir)

	// Find duplicates
	copies, err := findDuplicates(dirScan)
	if err != nil {
		log.Fatal(err)
	}

	// Print duplicates
	fmt.Println("Duplicates found:")
	for _, duplicates := range copies {
		if len(duplicates) > 1 {
			fmt.Println(duplicates)
		}
	}

	// Handle output file if specified
	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(fmt.Sprintf("%v", copies)), 0644)
		if err != nil {
			fmt.Println("Error writing to output file:", err)
		}
	}

	// Handle delete duplicates if specified
	if deleteDuplicates {
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
