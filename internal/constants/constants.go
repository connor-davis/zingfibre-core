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

var CreatedResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"message": openapi3.NewStringSchema().WithDefault(Created),
				"details": openapi3.NewStringSchema().WithDefault(CreatedDetails),
			}),
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

var SuccessObjectResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"message": openapi3.NewStringSchema().WithDefault(Success),
				"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
				"data":    openapi3.NewObjectSchema().WithDefault(map[string]any{}),
			}),
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
					"data":    map[string]any{},
				},
				Schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
					"message": openapi3.NewStringSchema().WithDefault(Success),
					"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
					"data":    openapi3.NewObjectSchema(),
				}).NewRef(),
			},
		}),
}

var SuccessArrayResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"message": openapi3.NewStringSchema().WithDefault(Success),
				"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
				"data":    openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
			}),
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
					"data":    []any{},
				},
				Schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
					"message": openapi3.NewStringSchema().WithDefault(Success),
					"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
					"data":    openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
				}).NewRef(),
			},
		}),
}

var SuccessPagingResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"message": openapi3.NewStringSchema().WithDefault(Success),
				"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
				"pages":   openapi3.NewIntegerSchema().WithDefault(1),
				"data":    openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
			}),
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
				Schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
					"message": openapi3.NewStringSchema().WithDefault(Success),
					"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
					"pages":   openapi3.NewIntegerSchema().WithDefault(1),
					"data":    openapi3.NewArraySchema().WithItems(openapi3.NewObjectSchema()),
				}).NewRef(),
			},
		}),
}

var SuccessResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"message": openapi3.NewStringSchema().WithDefault(Success),
				"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
			}),
		).
		WithDescription(SuccessDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]any{
					"message": Success,
					"details": SuccessDetails,
				},
				Schema: openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
					"message": openapi3.NewStringSchema().WithDefault(Success),
					"details": openapi3.NewStringSchema().WithDefault(SuccessDetails),
				}).NewRef(),
			},
		}),
}

var BadRequestResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"error":   openapi3.NewStringSchema().WithDefault(BadRequestError),
				"details": openapi3.NewStringSchema().WithDefault(BadRequestErrorDetails),
			}),
		).
		WithDescription(BadRequestErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   BadRequestError,
					"details": BadRequestErrorDetails,
				},
			},
		}),
}

var UnauthorizedResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"error":   openapi3.NewStringSchema().WithDefault(UnauthorizedError),
				"details": openapi3.NewStringSchema().WithDefault(UnauthorizedErrorDetails),
			}),
		).
		WithDescription(UnauthorizedErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   UnauthorizedError,
					"details": UnauthorizedErrorDetails,
				},
			},
		}),
}

var ConflictResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"error":   openapi3.NewStringSchema().WithDefault(ConflictError),
				"details": openapi3.NewStringSchema().WithDefault(ConflictErrorDetails),
			}),
		).
		WithDescription(ConflictErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   ConflictError,
					"details": ConflictErrorDetails,
				},
			},
		}),
}

var ForbiddenResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"error":   openapi3.NewStringSchema().WithDefault(ForbiddenError),
				"details": openapi3.NewStringSchema().WithDefault(ForbiddenErrorDetails),
			}),
		).
		WithDescription(ForbiddenErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   ForbiddenError,
					"details": ForbiddenErrorDetails,
				},
			},
		}),
}

var InternalServerErrorResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"error":   openapi3.NewStringSchema().WithDefault(InternalServerError),
				"details": openapi3.NewStringSchema().WithDefault(InternalServerErrorDetails),
			}),
		).
		WithDescription(InternalServerErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   InternalServerError,
					"details": InternalServerErrorDetails,
				},
			},
		}),
}

var NotFoundResponse = openapi3.ResponseRef{
	Value: openapi3.NewResponse().
		WithJSONSchema(
			openapi3.NewSchema().WithProperties(map[string]*openapi3.Schema{
				"error":   openapi3.NewStringSchema().WithDefault(NotFoundError),
				"details": openapi3.NewStringSchema().WithDefault(NotFoundErrorDetails),
			}),
		).
		WithDescription(NotFoundErrorDetails).
		WithContent(openapi3.Content{
			"application/json": &openapi3.MediaType{
				Example: map[string]string{
					"error":   NotFoundError,
					"details": NotFoundErrorDetails,
				},
			},
		}),
}
