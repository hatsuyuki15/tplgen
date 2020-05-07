package main

import (
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	var data = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: artifactory
  labels:
    app: artifactory
spec:
  replicas: 1
  strategy:
    type: Recreate
  template:
    metadata:
      labels:
        app: artifactory
    spec:
      containers:
        - name: artifactory
          image: docker.bintray.io/jfrog/artifactory-oss
          ports:
            - containerPort: 8081
              name: rest
            - containerPort: 8082
              name: ui
  selector:
    matchLabels:
      app: artifactory
`
	manifest := parseManifest(data)
	if strings.Compare(manifest.ApiVersion, "apps/v1") != 0 {
		t.Fatal("apiVersion mismatch")
	}
}

func TestPatchWithExistingNameSpace(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: yaml.MapSlice{
			{
				Key:   "Namespace",
				Value: "default",
			},
		},
		Spec: yaml.MapSlice{},
	}
	patch := Patch{
		Namespace: "dev",
		Resources: nil,
	}
	res := manifest.patch(patch)
	if !(len(res.Metadata) > 0 && res.Metadata[0].Value == "dev") {
		t.Fatal("Expected `dev` namespace but got:", res)
	}
}

func TestPatchWithNoNamespace(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata:   yaml.MapSlice{},
		Spec:       yaml.MapSlice{},
	}
	patch := Patch{
		Namespace: "dev",
		Resources: nil,
	}
	res := manifest.patch(patch)
	if !(len(res.Metadata) > 0 && res.Metadata[0].Value == "dev") {
		t.Fatal("Expected `dev` namespace but got:", res)
	}
}

func TestPatch_ShouldNotPatchEmpty(t *testing.T) {
	manifest := Manifest{
		ApiVersion: "apps/v1",
		Kind:       "Deployment",
		Metadata: yaml.MapSlice{
			{
				Key:   "Namespace",
				Value: "default",
			},
		},
		Spec: yaml.MapSlice{},
	}
	patch := Patch{
		Namespace: "",
		Resources: nil,
	}
	res := manifest.patch(patch)
	if manifest.Metadata[0].Value != "default" {
		t.Fatal("Expected `default` namespace but got:", res)
	}
}
