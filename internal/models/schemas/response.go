package schemas

import "github.com/getkin/kin-openapi/openapi3"

var SuccessResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault("Success"),
	"details": openapi3.NewStringSchema().WithDefault("Operation completed successfully"),
	"data": openapi3.NewOneOfSchema(
		UserSchema.Value,
		PointOfInterestSchema.Value,
		RechargeTypeCountSchema.Value,
		UserArraySchema.Value,
		PointOfInterestArraySchema.Value,
		RechargeTypeCountArraySchema.Value,
	),
	"pages": openapi3.NewIntegerSchema().WithDefault(1),
}).NewRef()

var ErrorResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault("Bad Request"),
	"details": openapi3.NewStringSchema().WithDefault("The request could not be understood or was missing required parameters."),
}).NewRef()
