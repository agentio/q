package main

import (
	"log"
	"os"
	"sort"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/genproto/googleapis/api/serviceconfig"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/testing/protocmp"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("%s", err)
	}
}

func run() error {
	b1, err := os.ReadFile("serviceconfig.json")
	if err != nil {
		return err
	}
	b2, err := os.ReadFile("current.json")
	if err != nil {
		return err
	}
	var c1 serviceconfig.Service
	err = protojson.Unmarshal(b1, &c1)
	if err != nil {
		return err
	}

	var c2 serviceconfig.Service
	err = protojson.Unmarshal(b2, &c2)
	if err != nil {
		return err
	}

	sort.Slice(c1.Types, func(i, j int) bool {
		return c1.Types[i].Name < c1.Types[j].Name
	})
	sort.Slice(c2.Types, func(i, j int) bool {
		return c2.Types[i].Name < c2.Types[j].Name
	})

	sort.Slice(c1.Documentation.Rules, func(i, j int) bool {
		return c1.Documentation.Rules[i].Selector < c1.Documentation.Rules[j].Selector
	})
	sort.Slice(c2.Documentation.Rules, func(i, j int) bool {
		return c2.Documentation.Rules[i].Selector < c2.Documentation.Rules[j].Selector
	})

	for _, t := range c1.Types {
		for _, field := range t.Fields {
			sort.Slice(field.Options, func(i, j int) bool {
				return field.Options[i].Name < field.Options[j].Name
			})
		}
	}
	for _, t := range c2.Types {
		for _, field := range t.Fields {
			sort.Slice(field.Options, func(i, j int) bool {
				return field.Options[i].Name < field.Options[j].Name
			})
		}
	}

	for _, a := range c1.Apis {
		sort.Slice(a.Methods, func(i, j int) bool {
			return a.Methods[i].Name < a.Methods[j].Name
		})
		for _, method := range a.Methods {
			sort.Slice(method.Options, func(i, j int) bool {
				return method.Options[i].Name < method.Options[j].Name
			})
		}
	}
	for _, a := range c2.Apis {
		sort.Slice(a.Methods, func(i, j int) bool {
			return a.Methods[i].Name < a.Methods[j].Name
		})
		for _, method := range a.Methods {
			sort.Slice(method.Options, func(i, j int) bool {
				return method.Options[i].Name < method.Options[j].Name
			})
		}
	}

	b1, err = protojson.Marshal(&c1)
	if err != nil {
		return err
	}
	err = os.WriteFile("c1.json", b1, 0644)
	if err != nil {
		return err
	}
	b2, err = protojson.Marshal(&c2)
	if err != nil {
		return err
	}
	err = os.WriteFile("c2.json", b2, 0644)
	if err != nil {
		return err
	}

	opts := cmp.Options{
		protocmp.Transform(),
	}
	if !cmp.Equal(&c2, &c1, opts) {
		log.Printf("GetDiff returned unexpected diff (-want +got):\n%s", cmp.Diff(&c2, &c1, opts))
	}
	return nil
}
