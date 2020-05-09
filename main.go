package main

import (
	"bytes"
	"fmt"
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

	files := listFiles(workingDir, "**/tplgen.yaml")
	for _, file := range files {
		data := readFile(file)
		patch := parsePatch(data)
		processPatch(patch, file, templateDir)
	}
}

func processPatch(patch Patch, patchFile string, templateDir string) {
	for _, resourcePath := range patch.Resources {
		resourceRootPath := filepath.Dir(patchFile)
		resources := listFiles(resourceRootPath, resourcePath)
		for _, resource := range resources {
			if filepath.Base(resource) != "tplgen.yaml" {
				processResource(resource, patch, templateDir)
			}
		}
	}
}

func processResource(resource string, patch Patch, templateDir string) {
	data := readFile(resource)
	ext := filepath.Ext(resource)
	if ext == ".tpl" {
		spec := parseSpec(data)
		spec = spec.patch(patch)
		evaluatedResult := evaluate(spec, templateDir)
		writeToStdout(resource, evaluatedResult)
	} else {
		manifest := parseManifest(data)
		manifest = manifest.patch(patch)
		writeToStdout(resource, manifest.evaluate())
	}
}

func writeToStdout(file string, content string) {
	fmt.Println("---")
	fmt.Println("#", file)
	fmt.Println(content)
}

func evaluate(spec Spec, templateDir string) string {
	templatePath := filepath.Join(templateDir, spec.Template)
	if !pathExist(templatePath) {
		log.Fatal("Template not exist: ", templatePath)
	}

	var tmpfile *os.File
	if isHelmTemplate(templatePath) {
		tmpfile = createTempFile(spec.toHelm())
	} else {
		tmpfile = createTempFile(spec.toYtt())
	}
	defer os.Remove(tmpfile.Name())

	var cmd *exec.Cmd
	if isHelmTemplate(templatePath) {
		cmd = exec.Command("helm", "template", "-n", spec.Namespace, spec.Name, templatePath, "-f", tmpfile.Name())
	} else {
		cmd = exec.Command("ytt", "-f", templatePath, "-f", tmpfile.Name())
	}

	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err := cmd.Run()
	if err != nil {
		log.Fatal(errOut.String())
	}

	return out.String()
}

func isHelmTemplate(path string) bool {
	return pathExist(filepath.Join(path, "Chart.yaml"))
}
