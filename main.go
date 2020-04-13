package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templateDir := filepath.Join(workingDir, "templates")
	if !pathExist(templateDir) {
		log.Fatal("Template dir doesn't exist: " + templateDir)
	}

	files := listFiles(workingDir)
	for _, file := range files {
		data := readFile(file)
		spec := parseSpec(data)
		evaluatedResult := evaluate(spec, templateDir)
		writeFile(file+".yaml", evaluatedResult)
	}
}

func evaluate(spec Spec, templateDir string) string {
	tmpfile := createTempFile(spec.toString())
	defer os.Remove(tmpfile.Name())

	templatePath := filepath.Join(templateDir, spec.Template)
	if !pathExist(templatePath) {
		log.Fatal("Template not exist: ", templatePath)
	}

	cmd := exec.Command("ytt", "-f", templatePath, "-f", tmpfile.Name())

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		log.Fatal(errOut.String())
	}

	return out.String()
}
