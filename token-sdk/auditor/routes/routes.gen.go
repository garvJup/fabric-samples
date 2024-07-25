// Package routes provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.3 DO NOT EDIT.
package routes

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
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	strictecho "github.com/oapi-codegen/runtime/strictmiddleware/echo"
)

// Account Information about an account and its balance
type Account struct {
	// Balance balance in base units for each currency
	Balance []Amount `json:"balance"`

	// Id account id as registered at the Certificate Authority
	Id string `json:"id"`
}

// Amount The amount to issue, transfer or redeem.
type Amount struct {
	// Code the code of the token
	Code string `json:"code"`

	// Value value in base units (usually cents)
	Value int64 `json:"value"`
}

// Error defines model for Error.
type Error struct {
	// Message High level error message
	Message string `json:"message"`

	// Payload Details about the error
	Payload string `json:"payload"`
}

// TransactionRecord A transaction
type TransactionRecord struct {
	// Amount The amount to issue, transfer or redeem.
	Amount Amount `json:"amount"`

	// Id transaction id
	Id string `json:"id"`

	// Message user provided message
	Message string `json:"message"`

	// Recipient the recipient of the transaction
	Recipient string `json:"recipient"`

	// Sender the sender of the transaction
	Sender string `json:"sender"`

	// Status Unknown | Pending | Confirmed | Deleted
	Status string `json:"status"`

	// Timestamp timestamp in the format: "2018-03-20T09:12:28Z"
	Timestamp time.Time `json:"timestamp"`
}

// Code The token code to filter on
type Code = string

// Id account id as registered at the Certificate Authority
type Id = string

// AccountSuccess defines model for AccountSuccess.
type AccountSuccess struct {
	Message string `json:"message"`

	// Payload Information about an account and its balance
	Payload Account `json:"payload"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse = Error

// HealthSuccess defines model for HealthSuccess.
type HealthSuccess struct {
	// Message ok
	Message string `json:"message"`
}

// TransactionsSuccess defines model for TransactionsSuccess.
type TransactionsSuccess struct {
	Message string              `json:"message"`
	Payload []TransactionRecord `json:"payload"`
}

// AuditorAccountParams defines parameters for AuditorAccount.
type AuditorAccountParams struct {
	Code *Code `form:"code,omitempty" json:"code,omitempty"`
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get an account and their balance of a certain type
	// (GET /auditor/accounts/{id})
	AuditorAccount(ctx echo.Context, id Id, params AuditorAccountParams) error
	// Get all transactions for an account
	// (GET /auditor/accounts/{id}/transactions)
	AuditorTransactions(ctx echo.Context, id Id) error
	// Returns 200 if the service is healthy
	// (GET /healthz)
	Healthz(ctx echo.Context) error
	// Returns 200 if the service is ready to accept calls
	// (GET /readyz)
	Readyz(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// AuditorAccount converts echo context to params.
func (w *ServerInterfaceWrapper) AuditorAccount(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params AuditorAccountParams
	// ------------- Optional query parameter "code" -------------

	err = runtime.BindQueryParameter("form", true, false, "code", ctx.QueryParams(), &params.Code)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter code: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuditorAccount(ctx, id, params)
	return err
}

// AuditorTransactions converts echo context to params.
func (w *ServerInterfaceWrapper) AuditorTransactions(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id Id

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.AuditorTransactions(ctx, id)
	return err
}

// Healthz converts echo context to params.
func (w *ServerInterfaceWrapper) Healthz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Healthz(ctx)
	return err
}

// Readyz converts echo context to params.
func (w *ServerInterfaceWrapper) Readyz(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.Readyz(ctx)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/auditor/accounts/:id", wrapper.AuditorAccount)
	router.GET(baseURL+"/auditor/accounts/:id/transactions", wrapper.AuditorTransactions)
	router.GET(baseURL+"/healthz", wrapper.Healthz)
	router.GET(baseURL+"/readyz", wrapper.Readyz)

}

type AccountSuccessJSONResponse struct {
	Message string `json:"message"`

	// Payload Information about an account and its balance
	Payload Account `json:"payload"`
}

type ErrorResponseJSONResponse Error

type HealthSuccessJSONResponse struct {
	// Message ok
	Message string `json:"message"`
}

type TransactionsSuccessJSONResponse struct {
	Message string              `json:"message"`
	Payload []TransactionRecord `json:"payload"`
}

type AuditorAccountRequestObject struct {
	Id     Id `json:"id"`
	Params AuditorAccountParams
}

type AuditorAccountResponseObject interface {
	VisitAuditorAccountResponse(w http.ResponseWriter) error
}

type AuditorAccount200JSONResponse struct{ AccountSuccessJSONResponse }

func (response AuditorAccount200JSONResponse) VisitAuditorAccountResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AuditorAccountdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response AuditorAccountdefaultJSONResponse) VisitAuditorAccountResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type AuditorTransactionsRequestObject struct {
	Id Id `json:"id"`
}

type AuditorTransactionsResponseObject interface {
	VisitAuditorTransactionsResponse(w http.ResponseWriter) error
}

type AuditorTransactions200JSONResponse struct {
	TransactionsSuccessJSONResponse
}

func (response AuditorTransactions200JSONResponse) VisitAuditorTransactionsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type AuditorTransactionsdefaultJSONResponse struct {
	Body       Error
	StatusCode int
}

func (response AuditorTransactionsdefaultJSONResponse) VisitAuditorTransactionsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)

	return json.NewEncoder(w).Encode(response.Body)
}

type HealthzRequestObject struct {
}

type HealthzResponseObject interface {
	VisitHealthzResponse(w http.ResponseWriter) error
}

type Healthz200JSONResponse struct{ HealthSuccessJSONResponse }

func (response Healthz200JSONResponse) VisitHealthzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Healthz503JSONResponse struct{ ErrorResponseJSONResponse }

func (response Healthz503JSONResponse) VisitHealthzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

type ReadyzRequestObject struct {
}

type ReadyzResponseObject interface {
	VisitReadyzResponse(w http.ResponseWriter) error
}

type Readyz200JSONResponse struct{ HealthSuccessJSONResponse }

func (response Readyz200JSONResponse) VisitReadyzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Readyz503JSONResponse struct{ ErrorResponseJSONResponse }

func (response Readyz503JSONResponse) VisitReadyzResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(503)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Get an account and their balance of a certain type
	// (GET /auditor/accounts/{id})
	AuditorAccount(ctx context.Context, request AuditorAccountRequestObject) (AuditorAccountResponseObject, error)
	// Get all transactions for an account
	// (GET /auditor/accounts/{id}/transactions)
	AuditorTransactions(ctx context.Context, request AuditorTransactionsRequestObject) (AuditorTransactionsResponseObject, error)
	// Returns 200 if the service is healthy
	// (GET /healthz)
	Healthz(ctx context.Context, request HealthzRequestObject) (HealthzResponseObject, error)
	// Returns 200 if the service is ready to accept calls
	// (GET /readyz)
	Readyz(ctx context.Context, request ReadyzRequestObject) (ReadyzResponseObject, error)
}

type StrictHandlerFunc = strictecho.StrictEchoHandlerFunc
type StrictMiddlewareFunc = strictecho.StrictEchoMiddlewareFunc

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
}

// AuditorAccount operation middleware
func (sh *strictHandler) AuditorAccount(ctx echo.Context, id Id, params AuditorAccountParams) error {
	var request AuditorAccountRequestObject

	request.Id = id
	request.Params = params

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AuditorAccount(ctx.Request().Context(), request.(AuditorAccountRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AuditorAccount")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(AuditorAccountResponseObject); ok {
		return validResponse.VisitAuditorAccountResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// AuditorTransactions operation middleware
func (sh *strictHandler) AuditorTransactions(ctx echo.Context, id Id) error {
	var request AuditorTransactionsRequestObject

	request.Id = id

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.AuditorTransactions(ctx.Request().Context(), request.(AuditorTransactionsRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "AuditorTransactions")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(AuditorTransactionsResponseObject); ok {
		return validResponse.VisitAuditorTransactionsResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Healthz operation middleware
func (sh *strictHandler) Healthz(ctx echo.Context) error {
	var request HealthzRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Healthz(ctx.Request().Context(), request.(HealthzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Healthz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(HealthzResponseObject); ok {
		return validResponse.VisitHealthzResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Readyz operation middleware
func (sh *strictHandler) Readyz(ctx echo.Context) error {
	var request ReadyzRequestObject

	handler := func(ctx echo.Context, request interface{}) (interface{}, error) {
		return sh.ssi.Readyz(ctx.Request().Context(), request.(ReadyzRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Readyz")
	}

	response, err := handler(ctx, request)

	if err != nil {
		return err
	} else if validResponse, ok := response.(ReadyzResponseObject); ok {
		return validResponse.VisitReadyzResponse(ctx.Response())
	} else if response != nil {
		return fmt.Errorf("unexpected response type: %T", response)
	}
	return nil
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"BgAA///nhAFuMBsAAA==",
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
