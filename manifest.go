package main

import (
	"gopkg.in/yaml.v2"
	"log"
	"strings"
)

type Manifest struct {
	data yaml.MapSlice
}

func parseManifests(data string) (manifests []Manifest) {
	data = strings.ReplaceAll(data, "\r\n", "\n")
	docs := strings.Split(data, "\n---\n")
	for _, doc := range docs {
		var mf Manifest
		err := yaml.Unmarshal([]byte(doc), &mf.data)
		if err != nil {
			log.Fatal(err)
		}
		manifests = append(manifests, mf)
	}
	return manifests
}

func (m Manifest) patch(patch Patch) Manifest {
	if patch.Namespace != "" {
		for i, item := range m.data {
			if item.Key == "metadata" {
				metadata := item.Value.(yaml.MapSlice)
				existNamespace := false
				for j, metadataElem := range metadata {
					if metadataElem.Key == "namespace" {
						metadataElem.Value = patch.Namespace
						metadata[j] = metadataElem
						existNamespace = true
					}
				}
				if !existNamespace {
					metadata = append(metadata, yaml.MapItem{
						Key:   "namespace",
						Value: patch.Namespace,
					})
					m.data[i].Value = metadata
				}
			}
		}
	}
	return m
}

func (m *Manifest) Namespace() string {
	metadata := get(m.data, "metadata").(yaml.MapSlice)
	namespace := get(metadata, "namespace")
	if namespace == nil {
		return ""
	} else {
		return namespace.(string)
	}
}

func get(slice yaml.MapSlice, key string) interface{} {
	for _, item := range slice {
		if item.Key == key {
			return item.Value
		}
	}
	return nil
}

func (m Manifest) evaluate() string {
	out, err := yaml.Marshal(m.data)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}
