package schemas

import "github.com/getkin/kin-openapi/openapi3"

var LoginRequestSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Email":    openapi3.NewStringSchema().WithFormat("email"),
	"Password": openapi3.NewStringSchema().WithMinLength(6).WithMaxLength(100),
}).NewRef()
