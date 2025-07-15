package schemas

import "github.com/getkin/kin-openapi/openapi3"

var PointOfInterestSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"ID":   openapi3.NewUUIDSchema(),
	"Name": openapi3.NewStringSchema(),
	"Key":  openapi3.NewStringSchema(),
}).NewRef()
