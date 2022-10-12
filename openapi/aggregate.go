package openapi

import (
	"fmt"
	"reflect"
)

func AggregateSpec(base Spec, sources map[string]Spec) (Spec, error) {
	for prefix, spec := range sources {
		base.Info.Description += fmt.Sprintf("\n\n# %s (`%s`)\n\n%s\n", spec.Info.Title, prefix, spec.Info.Description)
		// for _, spec := range sources {
		for path, endpoint := range spec.Paths {
			if existing, exists := base.Paths[path]; exists && !reflect.DeepEqual(endpoint, existing) {
				return base, fmt.Errorf("path %s already exists, cannot overwrite", path)
			}
			base.Paths[path] = endpoint
		}
		for name, schema := range spec.Components.Schemas {
			if existing, exists := base.Components.Schemas[name]; exists && !reflect.DeepEqual(schema, existing) {
				return base, fmt.Errorf("schema %s already exists, cannot overwrite", name)
			}
			base.Components.Schemas[name] = schema
		}
	}
	return base, nil
}
