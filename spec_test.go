package main

import (
	"testing"
)

func TestPatch(t *testing.T) {
	spec := Spec{
		Name:      "backend",
		Namespace: "default",
		Template:  "app",
		Values:    nil,
	}
	patch := Patch{
		Namespace: "dev",
		Resources: nil,
	}
	res := spec.patch(patch)
	if res.Namespace != "dev" {
		t.Fatal("Expected `dev` namespace but got:", res)
	}
}
