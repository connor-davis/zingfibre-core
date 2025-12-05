package schemas

import "github.com/getkin/kin-openapi/openapi3"

var DynamicQuerySchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"ID":    openapi3.NewUUIDSchema(),
	"Name":  openapi3.NewStringSchema(),
	"Query": openapi3.NewStringSchema(),
}).NewRef()

var DynamicQueryArraySchema = openapi3.NewArraySchema().WithItems(DynamicQuerySchema.Value).NewRef()

var CreateDynamicQuerySchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Name":   openapi3.NewStringSchema(),
	"Prompt": openapi3.NewStringSchema(),
}).NewRef()

var UpdateDynamicQuerySchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Name": openapi3.NewStringSchema(),
}).NewRef()

var DynamicQueryResultsSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Data": openapi3.NewArraySchema().WithItems(
		openapi3.NewObjectSchema().
			WithAdditionalProperties(openapi3.NewAnyOfSchema(
				openapi3.NewStringSchema(),
				openapi3.NewIntegerSchema(),
				openapi3.NewBoolSchema(),
				openapi3.NewDateTimeSchema(),
			)),
	),
	"Columns": openapi3.NewArraySchema().WithItems(
		openapi3.NewObjectSchema().WithProperties(map[string]*openapi3.Schema{
			"name":  openapi3.NewStringSchema(),
			"type":  openapi3.NewStringSchema(),
			"label": openapi3.NewStringSchema(),
		}).WithRequired([]string{
			"name",
			"type",
			"label",
		}),
	),
}).WithRequired([]string{
	"Data",
	"Columns",
}).NewRef()
