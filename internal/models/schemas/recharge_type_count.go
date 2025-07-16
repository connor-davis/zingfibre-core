package schemas

import "github.com/getkin/kin-openapi/openapi3"

var RechargeTypeCountsSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Items": openapi3.NewArraySchema().WithItems(
		openapi3.NewObjectSchema().
			WithProperties(map[string]*openapi3.Schema{
				"Period": openapi3.NewStringSchema(),
			}).
			WithAdditionalProperties(openapi3.NewIntegerSchema()),
	),
	"Types": openapi3.NewArraySchema().WithItems(
		openapi3.NewStringSchema(),
	),
}).NewRef()
