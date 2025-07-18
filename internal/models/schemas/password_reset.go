package schemas

import "github.com/getkin/kin-openapi/openapi3"

var PasswordResetSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"Code":        openapi3.NewStringSchema(),
	"NewPassword": openapi3.NewStringSchema(),
}).NewRef()
