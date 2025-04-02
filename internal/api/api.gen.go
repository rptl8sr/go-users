// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Error defines model for Error.
type Error struct {
	// Error Error message
	Error *string `json:"error,omitempty"`
}

// Health Health response
type Health struct {
	// Status Health response
	Status *string `json:"status,omitempty"`
}

// User defines model for User.
type User struct {
	// CreatedAt User creation timestamp
	CreatedAt time.Time `json:"created_at"`

	// Email User's email address
	Email openapi_types.Email `json:"email"`

	// FirstName User's first name
	FirstName string `json:"first_name"`

	// Id Unique user identifier
	Id uint `json:"id"`

	// LastName User's last name
	LastName string `json:"last_name"`

	// UpdatedAt User last update timestamp
	UpdatedAt time.Time `json:"updated_at"`
}

// UserRequest defines model for UserRequest.
type UserRequest struct {
	// Email User's email address
	Email openapi_types.Email `json:"email"`

	// FirstName User's first name
	FirstName string `json:"first_name"`

	// LastName User's last name
	LastName string `json:"last_name"`
}

// PostUserJSONRequestBody defines body for PostUser for application/json ContentType.
type PostUserJSONRequestBody = UserRequest

// PutUserJSONRequestBody defines body for PutUser for application/json ContentType.
type PutUserJSONRequestBody = UserRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Service Health
	// (GET /health)
	Health(w http.ResponseWriter, r *http.Request)
	// Create new user
	// (POST /users)
	PostUser(w http.ResponseWriter, r *http.Request)
	// Get user by ID
	// (GET /users/{id})
	GetUser(w http.ResponseWriter, r *http.Request, id uint)
	// Update user
	// (PUT /users/{id})
	PutUser(w http.ResponseWriter, r *http.Request, id uint)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Service Health
// (GET /health)
func (_ Unimplemented) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create new user
// (POST /users)
func (_ Unimplemented) PostUser(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get user by ID
// (GET /users/{id})
func (_ Unimplemented) GetUser(w http.ResponseWriter, r *http.Request, id uint) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Update user
// (PUT /users/{id})
func (_ Unimplemented) PutUser(w http.ResponseWriter, r *http.Request, id uint) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Health operation middleware
func (siw *ServerInterfaceWrapper) Health(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Health(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PostUser operation middleware
func (siw *ServerInterfaceWrapper) PostUser(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PostUser(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetUser operation middleware
func (siw *ServerInterfaceWrapper) GetUser(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id uint

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetUser(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// PutUser operation middleware
func (siw *ServerInterfaceWrapper) PutUser(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id uint

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.PutUser(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/health", wrapper.Health)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/users", wrapper.PostUser)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/users/{id}", wrapper.GetUser)
	})
	r.Group(func(r chi.Router) {
		r.Put(options.BaseURL+"/users/{id}", wrapper.PutUser)
	})

	return r
}

type HealthRequestObject struct {
}

type HealthResponseObject interface {
	VisitHealthResponse(w http.ResponseWriter) error
}

type Health200JSONResponse Health

func (response Health200JSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Health500JSONResponse Error

func (response Health500JSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type HealthdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response HealthdefaultJSONResponse) VisitHealthResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type PostUserRequestObject struct {
	Body *PostUserJSONRequestBody
}

type PostUserResponseObject interface {
	VisitPostUserResponse(w http.ResponseWriter) error
}

type PostUser201JSONResponse User

func (response PostUser201JSONResponse) VisitPostUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)

	return json.NewEncoder(w).Encode(response)
}

type PostUser400JSONResponse Error

func (response PostUser400JSONResponse) VisitPostUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PostUser409JSONResponse Error

func (response PostUser409JSONResponse) VisitPostUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(409)

	return json.NewEncoder(w).Encode(response)
}

type PostUser500JSONResponse Error

func (response PostUser500JSONResponse) VisitPostUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type GetUserRequestObject struct {
	Id uint `json:"id"`
}

type GetUserResponseObject interface {
	VisitGetUserResponse(w http.ResponseWriter) error
}

type GetUser200JSONResponse User

func (response GetUser200JSONResponse) VisitGetUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type GetUser404JSONResponse Error

func (response GetUser404JSONResponse) VisitGetUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type GetUser500JSONResponse Error

func (response GetUser500JSONResponse) VisitGetUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

type PutUserRequestObject struct {
	Id   uint `json:"id"`
	Body *PutUserJSONRequestBody
}

type PutUserResponseObject interface {
	VisitPutUserResponse(w http.ResponseWriter) error
}

type PutUser200JSONResponse User

func (response PutUser200JSONResponse) VisitPutUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type PutUser400JSONResponse Error

func (response PutUser400JSONResponse) VisitPutUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(400)

	return json.NewEncoder(w).Encode(response)
}

type PutUser404JSONResponse Error

func (response PutUser404JSONResponse) VisitPutUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(404)

	return json.NewEncoder(w).Encode(response)
}

type PutUser500JSONResponse Error

func (response PutUser500JSONResponse) VisitPutUserResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Service Health
	// (GET /health)
	Health(ctx context.Context, request HealthRequestObject) (HealthResponseObject, error)
	// Create new user
	// (POST /users)
	PostUser(ctx context.Context, request PostUserRequestObject) (PostUserResponseObject, error)
	// Get user by ID
	// (GET /users/{id})
	GetUser(ctx context.Context, request GetUserRequestObject) (GetUserResponseObject, error)
	// Update user
	// (PUT /users/{id})
	PutUser(ctx context.Context, request PutUserRequestObject) (PutUserResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// Health operation middleware
func (sh *strictHandler) Health(w http.ResponseWriter, r *http.Request) {
	var request HealthRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Health(ctx, request.(HealthRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Health")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(HealthResponseObject); ok {
		if err := validResponse.VisitHealthResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PostUser operation middleware
func (sh *strictHandler) PostUser(w http.ResponseWriter, r *http.Request) {
	var request PostUserRequestObject

	var body PostUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PostUser(ctx, request.(PostUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PostUser")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PostUserResponseObject); ok {
		if err := validResponse.VisitPostUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// GetUser operation middleware
func (sh *strictHandler) GetUser(w http.ResponseWriter, r *http.Request, id uint) {
	var request GetUserRequestObject

	request.Id = id

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.GetUser(ctx, request.(GetUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "GetUser")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(GetUserResponseObject); ok {
		if err := validResponse.VisitGetUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// PutUser operation middleware
func (sh *strictHandler) PutUser(w http.ResponseWriter, r *http.Request, id uint) {
	var request PutUserRequestObject

	request.Id = id

	var body PutUserJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		sh.options.RequestErrorHandlerFunc(w, r, fmt.Errorf("can't decode JSON body: %w", err))
		return
	}
	request.Body = &body

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.PutUser(ctx, request.(PutUserRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "PutUser")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PutUserResponseObject); ok {
		if err := validResponse.VisitPutUserResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RYUW/bNhD+K8RtQF6UWE6zh+lp3Tpk3h4WtOvLiqBgpbPNViIZHunVMPTfhyNlW7bk",
	"ZcWSJsCAPCjU8e7j3XcfT95AaRprNGpPUGyAyiU2Mj7+7Jxx/ICfZWNrjI9pDd4SOqGNF3MTdAVtBtYZ",
	"i84rpJ7dBiqk0inrldFQJJeiQSK5QMjAry1CAeSd0gto292K+fARS89+f0FZ++XQVVoXDskaTeysB5O8",
	"9IHYz6chtu1L9jiXofbJLrs3wL9Ay3k5SlnpUHqs3kuOc5lfXp3nL84v8z+meZHz358MvZGqhgICofuh",
	"23tRmgYymCtH/r2WDQf61Sw1ZKAqKKYZ1HL/5pVhhMFW98QaZKMP7zjHscrRQBktvGqQvGwsozKuiUE4",
	"3jm/GeZnd6wxt2ck4lshq8ohUd9n2jbir5+LE06jiYgmIw44cYONWt0FFJx6oSrUXs0Vuj6eoLTfe1Pa",
	"4wIdu+sV4AQctjiJpl+t0dTH3cnqi7PfZuDwLiiHFRTv+OTZLrG9PPYPkfXJcADv9gTVX+NdQPLHIvFl",
	"bD6m8VBMnhuN/lPdjwpzT02GmWcHSs/NMPbLm5mYG5eo3EgtF9hgoq7yXJqIjsTLmxlksEJHad/0Ir/I",
	"+VzGopZWQQEv4lIGVvplrMFkudPhBY7w9Q26lSpRdHIdfbmoG7NqJ6bAR09yGp1e5nmUIKM94yw2IK2t",
	"VRn3TT4Se95eSfz0rcM5FPDNZH9nTboLa9JFiOk5hPYafXBanP3+25lQc0HoVtzqJIIVUgsXtO7q+t0D",
	"4km35wicmfbotKzFmwSkM8z2t9FjI7hGjU7W28htBhSaRrr1WB29XBDztFu4ZfMJMyxW0BoaIcNPUUZI",
	"SKHxr05ZtfBLFLQmj82AHjeGfLw6U28g+R9NtX6wTPSlqj1sQO8CtgNaTh809FgNor5TKEskmoe6XotO",
	"e5kJV1+HhytZq0oobYMXlfQyhf7+8UPHw8vaoazWAj8r8vTU7XfQBIm/O/L2uiAKaL8JJhtVtSdVMUkP",
	"bVsgXU08S31Yi9mrQRtc474LHkkm/5GPu2H+Kr/6Siw4+IZ4NgS4Rp9qtq3Tcf35ZnSyQR918N3o9BY3",
	"Kv6Xb1HIoJtx4ih2qEBZ71D3jZ3tbQY2jJLN1rLEAdmGahueldjmTyC23Wz75GL7/26zt+nT5pTGsm3c",
	"nFosOP6omEirJqspcBt0O8a6j856E7BAXVmjtKd9F6YobTb4iWJrGkfpNPWel0ssPwmpK9H9eLBzsx07",
	"b9u/AwAA///X12lCShEAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
