package main

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"os"
	"strings"
)

func isDir(file *os.File) bool {

	// Get file info
	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

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

func calcHash(file *os.File) string {

	// Create a new hash
	hash := md5.New()

	// Copy the file contents to the hash
	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	// Return the hash as a string
	return fmt.Sprintf("%x", hash.Sum(nil))

}

func calcImageHash(img image.Image) string {
	// TODO: Implement a hash function for images

	// Using bounds of the image as a hash for now
	return fmt.Sprintf("%d%d", img.Bounds().Dx(), img.Bounds().Dy())
}

func findDuplicates(list []string) ([][]string, error) {
	duplicates := make(map[string][]string)
	imageDuplicates := make(map[string][]string)
	var result [][]string

	// Iterate over the list of files
	for i := range list {

		file, err := os.Open(list[i])
		if err != nil {
			fmt.Println("Error opening file:", err)
			return nil, err
		}
		defer file.Close()

		// If file is a directory, skip it
		if isDir(file) {
			continue
		}

		// Get file extension
		extension := strings.ToLower(strings.Split(list[i], ".")[len(strings.Split(list[i], "."))-1])

		// Calculate hash of the text file
		if extension == "txt" {
			hash := calcHash(file)
			duplicates[hash] = append(duplicates[hash], list[i])
		}

		// Check if file is an image
		if extension == "jpg" || extension == "jpeg" || extension == "png" {

			fmt.Println("Calculating image hash for", list[i])
			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Println("Error decoding image:", err)
				continue
			}

			hash := calcImageHash(img)
			imageDuplicates[hash] = append(imageDuplicates[hash], list[i])
		}

	}

	// Append duplicates to result
	for _, files := range duplicates {
		if len(files) > 1 {
			result = append(result, files)
		}
	}

	for _, files := range imageDuplicates {
		if len(files) > 1 {
			result = append(result, files)
		}
	}

	return result, nil
}
