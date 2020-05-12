package main

import (
	"gopkg.in/yaml.v2"
	"testing"
)

func TestParse(t *testing.T) {
	var data = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: artifactory
  namespace: default
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
	manifests := parseManifests(data)
	manifest := manifests[0]
	if manifest.Namespace() != "default" {
		t.Fatal("Expected `default` namespace got: ", manifest)
	}
}

func TestPatchWithExistingNameSpace(t *testing.T) {
	manifest := Manifest{
		data: yaml.MapSlice{
			{
				Key: "metadata",
				Value: yaml.MapSlice{
					{
						Key:   "namespace",
						Value: "default",
					},
				},
			},
		},
	}
	patch := Patch{
		Namespace: "dev",
		Resources: nil,
	}
	res := manifest.patch(patch)
	if res.Namespace() != "dev" {
		t.Fatal("Expected `dev` namespace got: ", res)
	}
}

func TestPatchWithNoNamespace(t *testing.T) {
	manifest := Manifest{
		data: yaml.MapSlice{
			{
				Key:   "metadata",
				Value: yaml.MapSlice{},
			},
		},
	}
	patch := Patch{
		Namespace: "dev",
		Resources: nil,
	}
	res := manifest.patch(patch)
	if res.Namespace() != "dev" {
		t.Fatal("Expected `dev` namespace got: ", res)
	}
}

func TestPatch_ShouldNotPatchEmpty(t *testing.T) {
	manifest := Manifest{
		data: yaml.MapSlice{
			{
				Key: "metadata",
				Value: yaml.MapSlice{
					{
						Key:   "namespace",
						Value: "default",
					},
				},
			},
		},
	}
	patch := Patch{
		Namespace: "",
		Resources: nil,
	}
	res := manifest.patch(patch)
	if res.Namespace() != "default" {
		t.Fatal("Expected `default` namespace got: ", res)
	}
}
