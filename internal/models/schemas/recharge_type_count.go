package schemas

import "github.com/getkin/kin-openapi/openapi3"

var RechargeTypeCountSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Count":  openapi3.NewIntegerSchema(),
	"Type":   openapi3.NewStringSchema(),
	"Period": openapi3.NewStringSchema(),
}).NewRef()

var RechargeTypeCountArraySchema = openapi3.NewArraySchema().WithProperties(map[string]*openapi3.Schema{
	"Count":  openapi3.NewIntegerSchema(),
	"Type":   openapi3.NewStringSchema(),
	"Period": openapi3.NewStringSchema(),
}).NewRef()
