package schemas

import "github.com/getkin/kin-openapi/openapi3"

var SuccessResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault("Success"),
	"details": openapi3.NewStringSchema().WithDefault("Operation completed successfully"),
	"data": openapi3.NewOneOfSchema(
		UserSchema.Value,
		PointOfPresenceSchema.Value,
		UserArraySchema.Value,
		PointsOfPresenceSchema.Value,
		RechargeTypeCountsSchema.Value,
		ReportRechargeTypeCountsSchema.Value,
		ReportCustomerSchema.Value,
		ReportCustomersSchema.Value,
		ReportExpiringCustomerSchema.Value,
		ReportExpiringCustomersSchema.Value,
		ReportRechargeSchema.Value,
		ReportRechargesSchema.Value,
		ReportRechargeSummarySchema.Value,
		ReportRechargeSummariesSchema.Value,
		ReportSummarySchema.Value,
		ReportSummariesSchema.Value,
	),
	"pages": openapi3.NewIntegerSchema().WithDefault(1),
}).NewRef()

var ErrorResponseSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault("Bad Request"),
	"details": openapi3.NewStringSchema().WithDefault("The request could not be understood or was missing required parameters."),
}).NewRef()
