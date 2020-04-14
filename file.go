package main

import (
	"github.com/bmatcuk/doublestar"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func readFile(file string) string {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("File reading error: "+file, err)
	}
	return string(data)
}

func listFiles(directory string) []string {
	globPattern := directory + "/**/*.tpl"
	files, err := doublestar.Glob(globPattern)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func excludePath(files []string, path string) (res []string) {
	for _, file := range files {
		if !hasSubPath(path, file) {
			res = append(res, file)
		}
	}
	return res
}

func hasSubPath(path string, subPath string) bool {
	relativePath, err := filepath.Rel(path, subPath)
	if err != nil {
		log.Fatal(err)
	}
	return !strings.Contains(relativePath, "..")
}

func writeFile(file string, content string) {
	err := ioutil.WriteFile(file, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func createTempFile(content string) *os.File {
	tmpfile, err := ioutil.TempFile("", "tpl*.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Write([]byte(content)); err != nil {
		tmpfile.Close()
		log.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
	return tmpfile
}

func pathExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
