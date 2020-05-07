package main

import (
	"gopkg.in/yaml.v2"
	"log"
)

type Manifest struct {
	ApiVersion string `yaml:"apiVersion"`
	Kind       string
	Metadata   yaml.MapSlice
	Spec       yaml.MapSlice
}

func parseManifest(data string) Manifest {
	var manifest Manifest
	err := yaml.Unmarshal([]byte(data), &manifest)
	if err != nil {
		log.Fatal(err)
	}
	return manifest
}

func (m Manifest) patch(patch Patch) Manifest {
	for i, item := range m.Metadata {
		if item.Key == "Namespace" {
			item.Value = patch.Namespace
			m.Metadata[i] = item
			return m
		}
	}
	m.Metadata = append(m.Metadata, yaml.MapItem{
		Key:   "Namespace",
		Value: patch.Namespace,
	})
	return m
}

func (m Manifest) evaluate() string {
	out, err := yaml.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
