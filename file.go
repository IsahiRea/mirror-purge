package main

import (
	"crypto/md5"
	"fmt"
	_ "image/gif"
	_ "image/png"
	"io"
	"io/fs"
	"os"
)

//TODO: Implement different logic for hash comparisons
// Default for now is to use hash comparisons

func scanDir(dir string) (list []string) {

	// Open the directory
	fileSys := os.DirFS(dir)
	dirSys, err := fs.ReadDir(fileSys, dir)
	if err != nil {
		return nil
	}

	for _, entry := range dirSys {
		path := dir + "/" + entry.Name()
		if entry.IsDir() {
			list = append(list, scanDir(path)...)
		} else {
			list = append(list, path)
		}
	}

	return list
}

func calcHash(filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new hash
	hash := md5.New()

	// Copy the file contents to the hash
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	// Return the hash as a string
	return fmt.Sprintf("%x", hash.Sum(nil)), nil

}

func findDuplicates(files []string) (map[string][]string, error) {

	hashes := make(map[string][]string)
	for _, file := range files {
		hash, err := calcHash(file)
		if err != nil {
			return nil, err
		}
		hashes[hash] = append(hashes[hash], file)
	}

	duplicates := make(map[string][]string)
	for hash, files := range hashes {
		if len(files) > 1 {
			duplicates[hash] = files
		}
	}

	return duplicates, nil
}
