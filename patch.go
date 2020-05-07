package main

import (
	"gopkg.in/yaml.v2"
	"log"
)

type Patch struct {
	Namespace string
	Resources []string
}

func parsePatch(data string) Patch {
	var patch Patch
	err := yaml.Unmarshal([]byte(data), &patch)
	if err != nil {
		log.Fatal(err)
	}
	return patch
}
