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

func (s Spec) toString() string {
	out, err := yaml.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	res := "#@data/values\n#@overlay/match-child-defaults missing_ok=True\n---\n" + string(out)
	return res
}

func parseSpec(data string) Spec {
	var spec Spec
	err := yaml.UnmarshalStrict([]byte(data), &spec)
	if err != nil {
		log.Fatal(err)
	}
	return spec
}
