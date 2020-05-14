package main

import (
	"gopkg.in/yaml.v2"
	"log"
)

type Spec struct {
	Name      string
	Namespace string
	Template  string
	Values    yaml.MapSlice
}

func (s Spec) toYtt() string {
	s.validate()

	out, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	res := "#@data/values\n#@overlay/match-child-defaults missing_ok=True\n---\n" + string(out)
	return res
}

func (s Spec) toHelm() string {
	s.validate()

	out, err := yaml.Marshal(s.Values)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
}

func (s Spec) validate() {
	if s.Namespace == "" {
		log.Fatal("Namespace must not empty: ", s)
	}
	if s.Name == "" {
		log.Fatal("Name must not be empty: ", s)
	}
}

func parseSpec(data string) Spec {
	var spec Spec
	err := yaml.UnmarshalStrict([]byte(data), &spec)
	if err != nil {
		log.Fatal(err)
	}
	return spec
}

func (s Spec) patch(patch Patch) Spec {
	if patch.Namespace != "" {
		s.Namespace = patch.Namespace
	}
	return s
}
