package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"hash"
	_ "image/gif"
	_ "image/png"
	"io"
	"log"
	"os"
	"path/filepath"
)

func scanDir(dir string) []string {

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return files
}

func calcHash(filePath string) (string, error) {

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new hash
	var hasher hash.Hash
	switch useHash {
	case "sha256":
		hasher = sha256.New()
	default:
		hasher = md5.New()
	}

	// Copy the file contents to the hash
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	// Return the hash as a string
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil

}

func findDuplicates(files []string) (map[string][]string, error) {

	// Create a map to store the hashes
	hashes := make(map[string][]string)

	// Calculate the hash for each file
	for _, file := range files {
		hash, err := calcHash(file)
		if err != nil {
			return nil, err
		}
		hashes[hash] = append(hashes[hash], file)
	}

	// Create a map to store the duplicates
	duplicates := make(map[string][]string)

	// Find files with the same hash
	for hash, files := range hashes {
		if len(files) > 1 {
			duplicates[hash] = files
		}
	}

	return duplicates, nil
}
