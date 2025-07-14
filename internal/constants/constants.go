package constants

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
)
