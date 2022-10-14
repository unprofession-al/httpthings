package openapi

import (
	"fmt"
	"reflect"
)

func AggregateSpec(base Spec, sources []Spec) (Spec, error) {
	tags := []tagObject{}
	for _, spec := range sources {
		tags = append(tags, tagObject{Name: spec.Info.Title, Description: spec.Info.Description})
		for path, eps := range spec.Paths {
			if existing, exists := base.Paths[path]; exists && !reflect.DeepEqual(eps, existing) {
				return base, fmt.Errorf("path %s already exists, cannot overwrite", path)
			}
			tagged := endpoints{}
			for verb, ep := range eps {
				ep.Tags = append(ep.Tags, spec.Info.Title)
				tagged[verb] = ep
			}
			base.Paths[path] = tagged
		}
		for name, schema := range spec.Components.Schemas {
			if existing, exists := base.Components.Schemas[name]; exists && !reflect.DeepEqual(schema, existing) {
				return base, fmt.Errorf("schema %s already exists, cannot overwrite", name)
			}
			base.Components.Schemas[name] = schema
		}
	}
	base.Tags = tags
	return base, nil
}
