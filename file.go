package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"os"
)

func scanDir(dir string) (list []string) {

	fileSys := os.DirFS(dir)
	dirSys, err := fs.ReadDir(fileSys, dir)
	if err != nil {
		return nil
	}

	for i := range dirSys {
		fileName := dirSys[i].Name()
		list = append(list, fileName)
	}

	return list
}

func calcHash(file *os.File) string {
	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))

}

func findDuplicates(list []string) ([]string, error) {

	var copies []string

	for i := range list {

		file, err := os.Open(list[i])
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if !isDir(file) {
			copies = append(copies, calcHash(file))
		}

	}

	//TODO: Return Duplicates instead of Hashes of the Files

	return copies, nil
}

func isDir(file *os.File) bool {

	fileInfo, err := file.Stat()
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}
