package schemas

import "github.com/getkin/kin-openapi/openapi3"

var MonthlyStatisticsSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Revenue":                 openapi3.NewInt64Schema(),
	"RevenueGrowth":           openapi3.NewInt64Schema(),
	"RevenueGrowthPercentage": openapi3.NewFloat64Schema(),
	"UniquePurchasers":        openapi3.NewInt64Schema(),
}).NewRef()
