package http_server

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/shiroyaavish/go-common/common"
	"github.com/shiroyaavish/go-common/errors/api_errors"
	"github.com/shiroyaavish/go-common/utils"
	"log"
	"os"
	"reflect"
	"time"
)

// RequestData represents the data structure used for handling HTTP requests.
// It contains the parameters, query, and body data from the request, as well as other contextual information.
// This type is generic and can be used with any parameter, query, and body types.
//
// The types of Params, Query and Body are set as such:
//   - P: Params
//   - Q: Query
//   - B: Body
type RequestData[P any, Q any, B any] struct {
	Params *P `json:"params"`
	Query  *Q `json:"query"`
	Body   *B `json:"body"`

	Context *fiber.Ctx `json:"-"`

	Pagination *Pagination `json:"pagination"`

	ProjectID         common.Project
	AccessId          *uuid.UUID
	AccessPermissions []string
}

// Pagination represents the data structure for pagination in API responses.
// It contains the limit and offset values for pagination.
// Limit specifies the maximum number of items per page, while offset specifies the number of items to skip.
// The JSON tags are used for serialization and deserialization.
type Pagination struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

// ResponseData represents the data structure used for handling HTTP responses.
// It contains the response data, status of whether there is more data available, message, duration, and hostname.
type ResponseData struct {
	Data     any    `json:"data,omitempty"`
	HasNext  bool   `json:"has_next,omitempty"`
	Message  string `json:"message,omitempty"`
	Duration int64  `json:"duration,omitempty"`
	Hostname string `json:"hostname,omitempty"`
}

// ResponseDataWithCustomStatus is a data structure representing a response with a custom status code.
// It contains the status code and the data to be returned in the response.
type ResponseDataWithCustomStatus struct {
	StatusCode int
	Data       any
}

// ResponseWithStatus returns a ResponseDataWithCustomStatus struct with the specified status code and data.
func ResponseWithStatus(statusCode int, data any) ResponseDataWithCustomStatus {
	return ResponseDataWithCustomStatus{
		StatusCode: statusCode,
		Data:       data,
	}
}

// RequestHandler is a type that represents a function used for handling HTTP requests.
// It takes a RequestData object as input and returns any value and an error.
// The RequestData object contains the parameters, query, and body data from the request,
// as well as other contextual information.
// This type is generic and can be used with any parameter, query, and body types.
//
// The types of Params, Query and Body are set as such:
//   - P: Params
//   - Q: Query
//   - B: Body
type RequestHandler[P any, Q any, B any] func(data RequestData[P, Q, B]) (any, error)

// Validator validates the data types provided
type Validator[T any] func(data *T) []utils.ErrorResponse

// RequestHandlerBuilder builds the request type and allow functions written on top of it that lets us add the handler,
// set the body type etc.
//
// The types of Params, Query and Body are set as such:
//   - P: Params
//   - Q: Query
//   - B: Body
type RequestHandlerBuilder[P any, Q any, B any] struct {
	Handler RequestHandler[P, Q, B]

	HasParams bool
	HasQuery  bool
	HasBody   bool

	ParamValidator Validator[P]
	QueryValidator Validator[Q]
	BodyValidator  Validator[B]

	AccessLevel common.AccessLevel

	ProjectIdRequired bool

	IsPaginated bool
	hostname    string
}

// NewHandler creates a new route
func NewHandler[P any, Q any, B any]() *RequestHandlerBuilder[P, Q, B] {
	return &RequestHandlerBuilder[P, Q, B]{
		AccessLevel: common.AccessLevelPublic,
	}
}

// SetHandler sets the RequestHandler with generic Types P Q and B
func (r *RequestHandlerBuilder[P, Q, B]) SetHandler(handler RequestHandler[P, Q, B]) *RequestHandlerBuilder[P, Q, B] {
	r.Handler = handler
	return r
}

// ParseParam parses the parameters (path parameters) as set in RequestData.Params
func (r *RequestHandlerBuilder[P, Q, B]) ParseParam() *RequestHandlerBuilder[P, Q, B] {
	r.HasParams = true
	return r
}

// ParseQuery parses the query parameters as set in RequestData.Query
func (r *RequestHandlerBuilder[P, Q, B]) ParseQuery() *RequestHandlerBuilder[P, Q, B] {
	r.HasQuery = true
	return r
}

// ParseBody parses the body (JSON Body) as set in RequestData.Body
func (r *RequestHandlerBuilder[P, Q, B]) ParseBody() *RequestHandlerBuilder[P, Q, B] {
	r.HasBody = true
	return r
}

// SetParamValidator sets a validation function as seen in utils.ValidateStruct
func (r *RequestHandlerBuilder[P, Q, B]) SetParamValidator(validator Validator[P]) *RequestHandlerBuilder[P, Q, B] {
	r.ParamValidator = validator
	return r
}

// SetQueryValidator sets the QueryValidator with generic Type Q for the RequestHandlerBuilder.
func (r *RequestHandlerBuilder[P, Q, B]) SetQueryValidator(validator Validator[Q]) *RequestHandlerBuilder[P, Q, B] {
	r.QueryValidator = validator
	return r
}

// SetBodyValidator sets the Validator for the request body of type B
func (r *RequestHandlerBuilder[P, Q, B]) SetBodyValidator(validator Validator[B]) *RequestHandlerBuilder[P, Q, B] {
	r.BodyValidator = validator
	return r
}

// SetAccessLevel sets the access level for the RequestHandler
func (r *RequestHandlerBuilder[P, Q, B]) SetAccessLevel(accessLevel common.AccessLevel) *RequestHandlerBuilder[P, Q, B] {
	r.AccessLevel = accessLevel
	return r
}

// SetPagination sets the IsPaginated flag for the RequestHandlerBuilder. When this flag is set to true, it indicates that the request handler should handle paginated data.
func (r *RequestHandlerBuilder[P, Q, B]) SetPagination() *RequestHandlerBuilder[P, Q, B] {
	r.IsPaginated = true
	return r
}

// SetProjectIdRequired sets the ProjectIdRequired field of the RequestHandlerBuilder to true and returns the updated builder.
func (r *RequestHandlerBuilder[P, Q, B]) SetProjectIdRequired() *RequestHandlerBuilder[P, Q, B] {
	r.ProjectIdRequired = true
	return r
}

// Build is a method of RequestHandlerBuilder[P, Q, B] that returns a fiber.Handler.
//
// It builds a handler function for processing HTTP requests.
// The handler function performs the following steps:
// - Sets the hostname of the request if available.
// - Creates a RequestData[P, Q, B] struct to hold the request data.
// - Sets the request data's Context field to the fiber.Ctx passed to the handler function.
// - Parses and validates the project ID from the request header if required.
// - Checks the access level of the request and returns an error response if access is denied.
// - Parses and validates the query string parameters if the handler builder has query parsing enabled.
// - Parses and validates the path parameters if the handler builder has path parameter parsing enabled.
// - Parses and validates the request body if the handler builder has request body parsing enabled.
// - Validates any validation errors encountered during parameter parsing and returns an error response if necessary.
// - Handles pagination if enabled, setting the pagination options based on the request query parameters.
// - Executes the user-defined handler function with the request data as input.
// - Handles any returned errors, returning an appropriate error response based on the error type.
// - Handles custom status codes for the response if the returned data implements the ResponseDataWithCustomStatus interface.
// - Processes paginated responses, slicing the result data if necessary.
// - Calculates the duration of the request in milliseconds.
// - Returns a JSON response with the response data, message, duration, and hostname, using the appropriate status code.
func (r *RequestHandlerBuilder[P, Q, B]) Build() fiber.Handler {
	if hostname, _ := os.Hostname(); len(hostname) > 0 {
		r.hostname = hostname
	}

	return func(c *fiber.Ctx) error {
		start := time.Now().UnixMilli()
		data := RequestData[P, Q, B]{
			Context:   c,
			ProjectID: common.ParseProjectFromString(c.Get("Project-ID")),
		}

		response := ResponseData{
			Message:  "ok",
			Hostname: c.Hostname(),
		}

		if r.ProjectIdRequired && data.ProjectID == common.UnknownProject {
			response.Message = api_errors.ErrInvalidParams.Message
			response.Duration = time.Now().UnixMilli() - start
			return c.Status(api_errors.ErrInvalidParams.StatusCode).JSON(response)
		}

		var accessErr *api_errors.Error
		data.AccessId, data.AccessPermissions, accessErr = r.AccessLevel.CheckAccess(c)
		if accessErr != nil {
			response.Message = accessErr.Message
			response.Duration = time.Now().UnixMilli() - start
			return c.Status(accessErr.StatusCode).JSON(response)
		}

		validationErr := make([]utils.ErrorResponse, 0)
		if r.HasQuery {
			data.Query = new(Q)
			if err := c.QueryParser(data.Query); err != nil {
				log.Println("Error parsing query", err)
				response.Message = api_errors.ErrInvalidParams.Message
				response.Duration = time.Now().UnixMilli() - start
				return c.Status(api_errors.ErrInvalidParams.StatusCode).JSON(response)
			}

			if r.QueryValidator != nil {
				validationErr = append(validationErr, r.QueryValidator(data.Query)...)
			} else {
				validationErr = append(validationErr, utils.ValidateStruct(data.Query)...)
			}
		}

		if r.HasParams {
			data.Params = new(P)
			if err := c.ParamsParser(data.Params); err != nil {
				log.Println("Error parsing params", err)
				response.Message = api_errors.ErrInvalidParams.Message
				response.Duration = time.Now().UnixMilli() - start
				return c.Status(api_errors.ErrInvalidParams.StatusCode).JSON(response)
			}
			if r.ParamValidator != nil {
				validationErr = append(validationErr, r.ParamValidator(data.Params)...)
			} else {
				validationErr = append(validationErr, utils.ValidateStruct(data.Params)...)
			}
		}

		if r.HasBody {
			data.Body = new(B)
			if err := c.BodyParser(data.Body); err != nil {
				log.Println("Error parsing body", err)
				response.Message = api_errors.ErrInvalidParams.Message
				response.Duration = time.Now().UnixMilli() - start
				return c.Status(api_errors.ErrInvalidParams.StatusCode).JSON(response)
			}
			if r.BodyValidator != nil {
				validationErr = append(validationErr, r.BodyValidator(data.Body)...)
			} else {
				validationErr = append(validationErr, utils.ValidateStruct(data.Body)...)
			}
		}

		if len(validationErr) > 0 {
			response.Message = api_errors.ErrInvalidParams.Message
			response.Duration = time.Now().UnixMilli() - start
			return c.Status(api_errors.ErrInvalidParams.StatusCode).JSON(response)
		}

		if r.IsPaginated {
			page := c.QueryInt("page", 1)
			perPage := c.QueryInt("per_page", 10)
			if perPage > 50 {
				perPage = 50
			}
			data.Pagination = &Pagination{
				Limit:  perPage + 1,
				Offset: (page - 1) * perPage,
			}

		}

		rData, err := r.Handler(data)

		var apiError api_errors.Error
		switch {
		case errors.As(err, &apiError):
			response.Message = apiError.Message
			response.Duration = time.Now().UnixMilli() - start
			return c.Status(apiError.StatusCode).JSON(response)
		case err != nil:
			log.Println("Error executing handler", err)
			response.Message = api_errors.ErrSomethingWentWrong.Message
			response.Duration = time.Now().UnixMilli() - start
			return c.Status(api_errors.ErrSomethingWentWrong.StatusCode).JSON(response)
		}

		statusCode := fiber.StatusOK

		if result, ok := rData.(ResponseDataWithCustomStatus); ok {
			response.Data = result.Data
			statusCode = result.StatusCode
		} else {
			response.Data = rData
		}

		if r.IsPaginated {
			rType := reflect.TypeOf(response.Data).Kind()
			if rType != reflect.Slice && rType != reflect.Array {
				panic("Paginated response must be an array or slice")
			}

			result := reflect.ValueOf(response.Data)
			response.HasNext = result.Len() == data.Pagination.Limit
			if response.HasNext {
				response.Data = result.Slice(0, result.Len()-1).Interface()
			} else {
				response.Data = result.Interface()
			}
		}

		response.Duration = time.Now().UnixMilli() - start

		return c.Status(statusCode).JSON(response)
	}
}
