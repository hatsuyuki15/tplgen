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
	out, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	res := "#@data/values\n#@overlay/match-child-defaults missing_ok=True\n---\n" + string(out)
	return res
}

func (s Spec) toHelm() string {
	out, err := yaml.Marshal(s.Values)
	if err != nil {
		log.Fatal(err)
	}
	return string(out)
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
	s.Namespace = patch.Namespace
	return s
}
