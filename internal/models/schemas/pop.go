package schemas

import "github.com/getkin/kin-openapi/openapi3"

var PointOfPresenceSchema = openapi3.NewStringSchema().NewRef()

var PointsOfPresenceSchema = openapi3.NewArraySchema().WithItems(PointOfPresenceSchema.Value).NewRef()
