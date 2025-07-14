package schemas

import "github.com/getkin/kin-openapi/openapi3"

var UserSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"ID":          openapi3.NewUUIDSchema(),
	"Email":       openapi3.NewStringSchema().WithFormat("email"),
	"MfaEnabled":  openapi3.NewBoolSchema(),
	"MfaVerified": openapi3.NewBoolSchema(),
	"Role":        openapi3.NewStringSchema().WithEnum([]interface{}{"admin", "staff", "user"}).WithDefault("user"),
}).NewRef()
