// Package apigen provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package apigen

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oapi-codegen/runtime"
)

// TaskModel defines model for TaskModel.
type TaskModel struct {
	Attempts        int32     `json:"attempts"`
	CreatedAt       time.Time `json:"created_at"`
	CreatedBy       string    `json:"created_by"`
	ID              string    `json:"id"`
	Payload         *string   `json:"payload"`
	Status          string    `json:"status"`
	StatusUpdatedAt time.Time `json:"status_updated_at"`
	Type            string    `json:"type"`
}

// TokenModel defines model for TokenModel.
type TokenModel struct {
	Token string `json:"token"`
}

// SingInResponse defines model for SingInResponse.
type SingInResponse struct {
	Data TokenModel `json:"data"`
}

// TaskResponse defines model for TaskResponse.
type TaskResponse struct {
	Task TaskModel `json:"task"`
}

// CreateTaskRequest defines model for CreateTaskRequest.
type CreateTaskRequest struct {
	Payload *string `json:"payload,omitempty"`
	Type    string  `json:"type"`
}

// SignInRequest defines model for SignInRequest.
type SignInRequest struct {
	Captcha  string `json:"captcha"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpRequest defines model for SignUpRequest.
type SignUpRequest struct {
	Captcha  string `json:"captcha"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignInJSONBody defines parameters for SignIn.
type SignInJSONBody struct {
	Captcha  string `json:"captcha"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpJSONBody defines parameters for SignUp.
type SignUpJSONBody struct {
	Captcha  string `json:"captcha"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateTaskJSONBody defines parameters for CreateTask.
type CreateTaskJSONBody struct {
	Payload *string `json:"payload,omitempty"`
	Type    string  `json:"type"`
}

// SignInJSONRequestBody defines body for SignIn for application/json ContentType.
type SignInJSONRequestBody SignInJSONBody

// SignUpJSONRequestBody defines body for SignUp for application/json ContentType.
type SignUpJSONRequestBody SignUpJSONBody

// CreateTaskJSONRequestBody defines body for CreateTask for application/json ContentType.
type CreateTaskJSONRequestBody CreateTaskJSONBody

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Ready
	// (GET /ready)
	Ready(w http.ResponseWriter, r *http.Request)
	// Sign In User
	// (POST /v1/signin)
	SignIn(w http.ResponseWriter, r *http.Request)
	// Sign Up New User
	// (POST /v1/signup)
	SignUp(w http.ResponseWriter, r *http.Request)
	// Create Task
	// (POST /v1/tasks)
	CreateTask(w http.ResponseWriter, r *http.Request)
	// Get task
	// (GET /v1/tasks/{id})
	GetTask(w http.ResponseWriter, r *http.Request, id string)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Ready
// (GET /ready)
func (_ Unimplemented) Ready(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Sign In User
// (POST /v1/signin)
func (_ Unimplemented) SignIn(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Sign Up New User
// (POST /v1/signup)
func (_ Unimplemented) SignUp(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Create Task
// (POST /v1/tasks)
func (_ Unimplemented) CreateTask(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// Get task
// (GET /v1/tasks/{id})
func (_ Unimplemented) GetTask(w http.ResponseWriter, r *http.Request, id string) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Ready operation middleware
func (siw *ServerInterfaceWrapper) Ready(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Ready(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// SignIn operation middleware
func (siw *ServerInterfaceWrapper) SignIn(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SignIn(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// SignUp operation middleware
func (siw *ServerInterfaceWrapper) SignUp(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.SignUp(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateTask operation middleware
func (siw *ServerInterfaceWrapper) CreateTask(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateTask(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// GetTask operation middleware
func (siw *ServerInterfaceWrapper) GetTask(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", chi.URLParam(r, "id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.GetTask(w, r, id)
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
		r.Get(options.BaseURL+"/ready", wrapper.Ready)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/signin", wrapper.SignIn)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/signup", wrapper.SignUp)
	})
	r.Group(func(r chi.Router) {
		r.Post(options.BaseURL+"/v1/tasks", wrapper.CreateTask)
	})
	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/v1/tasks/{id}", wrapper.GetTask)
	})

	return r
}
