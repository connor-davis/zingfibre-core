package schemas

import "github.com/getkin/kin-openapi/openapi3"

var ResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(""),
	"details": openapi3.NewStringSchema().WithDefault(""),
	"data": openapi3.NewAnyOfSchema([]*openapi3.Schema{
		UserSchema.Value,
		openapi3.NewArraySchema().WithItems(UserSchema.Value),
		PointOfInterestSchema.Value,
		openapi3.NewArraySchema().WithItems(PointOfInterestSchema.Value),
	}...),
	"pages": openapi3.NewIntegerSchema().WithDefault(1),
}).NewRef()
