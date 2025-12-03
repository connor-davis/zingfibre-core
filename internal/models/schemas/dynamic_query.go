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
