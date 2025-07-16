package schemas

import "github.com/getkin/kin-openapi/openapi3"

var ReportCustomerSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"FirstName":      openapi3.NewStringSchema(),
	"Surname":        openapi3.NewStringSchema(),
	"Email":          openapi3.NewStringSchema().WithFormat("email"),
	"PhoneNumber":    openapi3.NewStringSchema(),
	"RadiusUsername": openapi3.NewStringSchema(),
}).NewRef()

var ReportCustomersSchema = openapi3.NewArraySchema().WithItems(ReportCustomerSchema.Value).NewRef()
