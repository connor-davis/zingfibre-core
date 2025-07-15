package constants

import "github.com/getkin/kin-openapi/openapi3"

const (
	InternalServerError        string = "Internal server error"
	InternalServerErrorDetails string = "An unexpected error occurred. Please try again later or contact support."
	UnauthorizedError          string = "Unauthorized"
	UnauthorizedErrorDetails   string = "You are not authorized to access this resource. Please log in or contact support."
	NotFoundError              string = "Not Found"
	NotFoundErrorDetails       string = "The requested resource could not be found. Please check the URL or contact support."
	BadRequestError            string = "Bad Request"
	BadRequestErrorDetails     string = "The request could not be understood or was missing required parameters."
	ConflictError              string = "Conflict"
	ConflictErrorDetails       string = "The request could not be completed due to a conflict with the current state of the resource."
	ForbiddenError             string = "Forbidden"
	ForbiddenErrorDetails      string = "You do not have permission to access this resource. Please check your permissions or contact support."
	Created                    string = "Created"
	CreatedDetails             string = "The resource has been successfully created."
	Success                    string = "Success"
	SuccessDetails             string = "The request was successful."
)

var CreatedSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault(Created),
	"details": openapi3.NewStringSchema().WithDefault(CreatedDetails),
})

var CreatedResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			CreatedSchema,
		).
		WithDescription(CreatedDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"message": Created,
					"details": CreatedDetails,
				},
			},
		}),
}

var SuccessObjectSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault(Success),
	"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
	"data":    openapi3.NewObjectSchema().WithDefault(map[string]any{}),
})

var SuccessObjectResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			SuccessObjectSchema,
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
					"data":    map[string]any{},
				},
				Schema: SuccessObjectSchema.NewRef(),
			},
		}),
}

var SuccessArraySchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault(Success),
	"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
	"data":    openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
})

var SuccessArrayResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			SuccessArraySchema,
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
					"data":    []any{},
				},
				Schema: SuccessArraySchema.NewRef(),
			},
		}),
}

var SuccessPagingSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault(Success),
	"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
	"pages":   openapi3.NewIntegerSchema().WithDefault(1),
	"data":    openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
})

var SuccessPagingResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			SuccessPagingSchema,
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
					"pages":   1,
					"data":    []any{},
				},
				Schema: SuccessPagingSchema.NewRef(),
			},
		}),
}

var SuccessSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"message": openapi3.NewStringSchema().WithDefault(Success),
	"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
})

var SuccessResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			SuccessSchema,
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
				},
				Schema: SuccessSchema.NewRef(),
			},
		}),
}

var BadRequestSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(BadRequestError),
	"details": openapi3.NewStringSchema().WithDefault(BadRequestErrorDetails),
})

var BadRequestResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			BadRequestSchema,
		).
		WithDescription(BadRequestErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   BadRequestError,
					"details": BadRequestErrorDetails,
				},
				Schema: BadRequestSchema.NewRef(),
			},
		}),
}

var UnauthorizedSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(UnauthorizedError),
	"details": openapi3.NewStringSchema().WithDefault(UnauthorizedErrorDetails),
})

var UnauthorizedResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			UnauthorizedSchema,
		).
		WithDescription(UnauthorizedErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   UnauthorizedError,
					"details": UnauthorizedErrorDetails,
				},
				Schema: UnauthorizedSchema.NewRef(),
			},
		}),
}

var ConflictSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(ConflictError),
	"details": openapi3.NewStringSchema().WithDefault(ConflictErrorDetails),
})

var ConflictResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			ConflictSchema,
		).
		WithDescription(ConflictErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   ConflictError,
					"details": ConflictErrorDetails,
				},
				Schema: ConflictSchema.NewRef(),
			},
		}),
}

var ForbiddenSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(ForbiddenError),
	"details": openapi3.NewStringSchema().WithDefault(ForbiddenErrorDetails),
})

var ForbiddenResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			ForbiddenSchema,
		).
		WithDescription(ForbiddenErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   ForbiddenError,
					"details": ForbiddenErrorDetails,
				},
				Schema: ForbiddenSchema.NewRef(),
			},
		}),
}

var InternalServerErrorSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(InternalServerError),
	"details": openapi3.NewStringSchema().WithDefault(InternalServerErrorDetails),
})

var InternalServerErrorResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			InternalServerErrorSchema,
		).
		WithDescription(InternalServerErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   InternalServerError,
					"details": InternalServerErrorDetails,
				},
				Schema: InternalServerErrorSchema.NewRef(),
			},
		}),
}

var NotFoundSchema = openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
	"error":   openapi3.NewStringSchema().WithDefault(NotFoundError),
	"details": openapi3.NewStringSchema().WithDefault(NotFoundErrorDetails),
})

var NotFoundResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			NotFoundSchema,
		).
		WithDescription(NotFoundErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   NotFoundError,
					"details": NotFoundErrorDetails,
				},
				Schema: NotFoundSchema.NewRef(),
			},
		}),
}
