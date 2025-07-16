package schemas

import "github.com/getkin/kin-openapi/openapi3"

var ReportRechargeTypeCountSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"RechargeName":    openapi3.NewStringSchema(),
	"RechargeCount":   openapi3.NewIntegerSchema(),
	"RechargePeriod":  openapi3.NewStringSchema(),
	"RechargeMaxDate": openapi3.NewStringSchema().WithFormat("date-time"),
}).NewRef()

var ReportRechargeTypeCountsSchema = openapi3.NewArraySchema().WithItems(ReportRechargeTypeCountSchema.Value).NewRef()

var ReportCustomerSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"FirstName":      openapi3.NewStringSchema(),
	"Surname":        openapi3.NewStringSchema(),
	"Email":          openapi3.NewStringSchema().WithFormat("email"),
	"PhoneNumber":    openapi3.NewStringSchema(),
	"RadiusUsername": openapi3.NewStringSchema(),
}).NewRef()

var ReportCustomersSchema = openapi3.NewArraySchema().WithItems(ReportCustomerSchema.Value).NewRef()

var ReportExpiringCustomerSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"FirstName":            openapi3.NewStringSchema(),
	"Surname":              openapi3.NewStringSchema(),
	"Email":                openapi3.NewStringSchema().WithFormat("email"),
	"PhoneNumber":          openapi3.NewStringSchema(),
	"RadiusUsername":       openapi3.NewStringSchema(),
	"LastPurchaseDuration": openapi3.NewStringSchema(),
	"LastPurchaseSpeed":    openapi3.NewStringSchema(),
	"Expiration":           openapi3.NewStringSchema(),
	"Address":              openapi3.NewStringSchema(),
}).NewRef()

var ReportExpiringCustomersSchema = openapi3.NewArraySchema().WithItems(ReportExpiringCustomerSchema.Value).NewRef()

var ReportRechargeSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"DateCreated": openapi3.NewStringSchema().WithFormat("date-time"),
	"Email":       openapi3.NewStringSchema().WithFormat("email"),
	"FirstName":   openapi3.NewStringSchema(),
	"Surname":     openapi3.NewStringSchema(),
	"ItemName":    openapi3.NewStringSchema(),
	"Amount":      openapi3.NewFloat64Schema(),
	"Successful":  openapi3.NewBoolSchema(),
	"ServiceId":   openapi3.NewInt64Schema(),
	"BuildName":   openapi3.NewStringSchema(),
	"BuildType":   openapi3.NewStringSchema(),
}).NewRef()

var ReportRechargesSchema = openapi3.NewArraySchema().WithItems(ReportRechargeSchema.Value).NewRef()

// Schema for GetReportsRechargesSummary (same as ReportRechargeSchema, but for summary)
var ReportRechargeSummarySchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"DateCreated": openapi3.NewStringSchema().WithFormat("date-time"),
	"Email":       openapi3.NewStringSchema().WithFormat("email"),
	"FirstName":   openapi3.NewStringSchema(),
	"Surname":     openapi3.NewStringSchema(),
	"ItemName":    openapi3.NewStringSchema(),
	"Amount":      openapi3.NewFloat64Schema(),
	"Successful":  openapi3.NewBoolSchema(),
	"ServiceId":   openapi3.NewInt64Schema(),
	"BuildName":   openapi3.NewStringSchema(),
	"BuildType":   openapi3.NewStringSchema(),
}).NewRef()

var ReportRechargeSummariesSchema = openapi3.NewArraySchema().WithItems(ReportRechargeSummarySchema.Value).NewRef()

var ReportSummarySchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"DateCreated":    openapi3.NewStringSchema().WithFormat("date-time"),
	"ItemName":       openapi3.NewStringSchema(),
	"RadiusUsername": openapi3.NewStringSchema(),
	"AmountGross":    openapi3.NewStringSchema(),
	"AmountFee":      openapi3.NewStringSchema(),
	"AmountNet":      openapi3.NewStringSchema(),
	"CashCode":       openapi3.NewStringSchema(),
	"CashAmount":     openapi3.NewStringSchema(),
	"ServiceId":      openapi3.NewInt64Schema(),
	"BuildName":      openapi3.NewStringSchema(),
	"BuildType":      openapi3.NewStringSchema(),
}).NewRef()

var ReportSummariesSchema = openapi3.NewArraySchema().WithItems(ReportSummarySchema.Value).NewRef()
