package main

import (
	"bytes"
	"github.com/bmatcuk/doublestar"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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

	workingDir = filepath.Join(workingDir, "..", "iac")
	templateDir := filepath.Join(workingDir, "templates")

	globPattern := workingDir + "/**/*.tpl"
	files, err := doublestar.Glob(globPattern)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		log.Println("Processing: " + file)
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal("File reading error: "+file, err)
		}
		spec := parseSpec(string(data))
		result := evaluateTemplate(spec, templateDir)
		writeToFile(file, result)
	}
}

func writeToFile(file string, content string) {
	newFile := file + ".yaml"
	err := ioutil.WriteFile(newFile, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func evaluateTemplate(spec Spec, templateDir string) string {
	tmpfile := createTempFile(spec.toString())
	defer os.Remove(tmpfile.Name())

	templatePath := filepath.Join(templateDir, spec.Template)
	if !pathExist(templatePath) {
		log.Fatalln("Template not exist: ", templatePath)
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

func parseSpec(data string) Spec {
	var spec Spec
	err := yaml.UnmarshalStrict([]byte(data), &spec)
	if err != nil {
		log.Fatal(err)
	}
	return spec
}

type Spec struct {
	Template string
	Values   yaml.MapSlice
}

func (s Spec) toString() string {
	out, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	res := "#@data/values\n#@overlay/match-child-defaults missing_ok=True\n---\n" + string(out)
	return res
}

func pathExist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}
