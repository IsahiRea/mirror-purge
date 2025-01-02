package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"strings"

	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
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
	// Resize image to 8x8
	const size = 8
	resized := resizeImage(img, size, size)

	// Convert image to grayscale
	gary := image.NewGray(resized.Bounds())
	draw.Draw(gary, gary.Bounds(), resized, image.Point{}, draw.Src)

	// Create a new hash
	hash := md5.New()
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			// Get the color of the pixel
			color := gary.GrayAt(x, y)
			// Add the color to the hash
			hash.Write([]byte{color.Y})
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func resizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dst
}

func calcVideoHash(file *os.File) string {
	cmd := exec.Command("ffmpeg", "-i", file.Name(), "-vf", "fps=1", "-q:v", "2", "-f", "mjpeg", "pipe:1")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println("Error Extracting Frames:", err)
		return ""
	}

	hash := sha256.New()

	//TODO: Implement video hashing
	reader := bytes.NewReader(out.Bytes())
	for {
		img, err := jpeg.Decode(reader)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				break
			}

			fmt.Println("Error decoding image:", err)
			continue
		}

		const size = 8
		resized := resize.Resize(size, size, img, resize.Lanczos3)

		gray := image.NewGray(resized.Bounds())
		draw.Draw(gray, gray.Bounds(), resized, image.Point{}, draw.Src)

		for y := 0; y < size; y++ {
			for x := 0; x < size; x++ {
				color := gray.GrayAt(x, y)
				hash.Write([]byte{color.Y})
			}
		}
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func findDuplicates(list []string) ([][]string, error) {
	duplicates := make(map[string][]string)
	imageDuplicates := make(map[string][]string)
	videoDuplicates := make(map[string][]string)
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

			img, _, err := image.Decode(file)
			if err != nil {
				fmt.Println("Error decoding image:", err)
				continue
			}

			hash := calcImageHash(img)
			imageDuplicates[hash] = append(imageDuplicates[hash], list[i])
		}

		if extension == "mp4" || extension == "avi" || extension == "mkv" {
			hash := calcVideoHash(file)
			videoDuplicates[hash] = append(videoDuplicates[hash], list[i])
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
