package openapi

import (
	"fmt"
	"reflect"
)

// AggregateSpec takes a base [Doc] and expands this base with the content of all soucre [Doc]s.
func AggregateOpenAPIDoc(base Doc, sources []Doc) (Doc, error) {
	tags := []Tag{}
	for _, spec := range sources {
		tags = append(tags, Tag{Name: spec.Info.Title, Description: spec.Info.Description})
		for path, pathItem := range spec.Paths {
			if existing, exists := base.Paths[path]; exists && !reflect.DeepEqual(pathItem, existing) {
				return base, fmt.Errorf("path %s already exists, cannot overwrite", path)
			}
			base.Paths[path] = appendTags(pathItem, spec.Info.Title)
		}
		for name, schema := range spec.Components.Schemas {
			if existing, exists := base.Components.Schemas[name]; exists && !reflect.DeepEqual(schema, existing) {
				return base, fmt.Errorf("schema %s already exists, cannot overwrite", name)
			}
			base.Components.Schemas[name] = schema
		}
		for name, scheme := range spec.Components.SecuritySchemes {
			if existing, exists := base.Components.SecuritySchemes[name]; exists && !reflect.DeepEqual(scheme, existing) {
				return base, fmt.Errorf("security scheme %s already exists, cannot overwrite", name)
			}
			base.Components.SecuritySchemes[name] = scheme
		}
	}
	base.Tags = tags
	return base, nil
}

func appendTags(path PathItem, tag ...string) PathItem {
	if path.Get != nil {
		path.Get.Tags = append(path.Get.Tags, tag...)
	}
	if path.Put != nil {
		path.Put.Tags = append(path.Put.Tags, tag...)
	}
	if path.Post != nil {
		path.Post.Tags = append(path.Post.Tags, tag...)
	}
	if path.Delete != nil {
		path.Delete.Tags = append(path.Delete.Tags, tag...)
	}
	if path.Options != nil {
		path.Options.Tags = append(path.Options.Tags, tag...)
	}
	if path.Head != nil {
		path.Head.Tags = append(path.Head.Tags, tag...)
	}
	if path.Patch != nil {
		path.Patch.Tags = append(path.Patch.Tags, tag...)
	}
	if path.Trace != nil {
		path.Trace.Tags = append(path.Trace.Tags, tag...)
	}
	return path
}
