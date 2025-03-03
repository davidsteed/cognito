// Package Openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by unknown module path version unknown version DO NOT EDIT.
package Openapi

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	CognitoScopes = "cognito.Scopes"
)

// Location defines model for Location.
type Location struct {
	AddressLine1   string  `json:"addressLine1"`
	DisabledAccess *string `json:"disabledAccess,omitempty"`
	Facia          string  `json:"facia"`
	Id             string  `json:"id"`
	OpenDays       struct {
		Friday    Day `json:"friday"`
		Monday    Day `json:"monday"`
		Saturday  Day `json:"saturday"`
		Sunday    Day `json:"sunday"`
		Thursday  Day `json:"thursday"`
		Tuesday   Day `json:"tuesday"`
		Wednesday Day `json:"wednesday"`
	} `json:"openDays"`
	PostCode string   `json:"postCode"`
	Products []string `json:"products"`
	Town     string   `json:"town"`
}

// Day defines model for day.
type Day struct {
	End   string `json:"end"`
	Open  bool   `json:"open"`
	Start string `json:"start"`
}

// PostLocationJSONRequestBody defines body for PostLocation for application/json ContentType.
type PostLocationJSONRequestBody = Location

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /v1/location)
	GetLocation(ctx echo.Context) error

	// (POST /v1/location)
	PostLocation(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetLocation converts echo context to params.
func (w *ServerInterfaceWrapper) GetLocation(ctx echo.Context) error {
	var err error

	ctx.Set(CognitoScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetLocation(ctx)
	return err
}

// PostLocation converts echo context to params.
func (w *ServerInterfaceWrapper) PostLocation(ctx echo.Context) error {
	var err error

	ctx.Set(CognitoScopes, []string{})

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostLocation(ctx)
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

	router.GET(baseURL+"/v1/location", wrapper.GetLocation)
	router.POST(baseURL+"/v1/location", wrapper.PostLocation)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/5RUXW/iOhD9K8i3j4HQe9/8Rtu7q2qr1arsah8qVA32AK6I7dpjWIry31d2EghN+sGb",
	"kxmPzzlzZvZMmMIajZo843vmxQoLSMc7I4CU0fFsnbHoSGGKgJQOvb9TGi/jN+0sMs48OaWXrMyYVB7m",
	"a5QTIdD73pQFCAW9ESV7fxuL+gZ2vgtn4ZSEXTxdOFwwzv7Jj6TymlEeU8qMFUZ/PtkDBXdGejijNq2C",
	"82ekBzwje4tSfz6/zJjD56AcSsYfGomOb7brtYBnjfItoQ4izLKmhWb+hIIiKms8XRuJvf21zsggKhcq",
	"wqLfNvUPcK6WxWx1T+IrRimrZaHs1MKNG1sIkw1boPro1PKemhH12/ZtBebGrBF0cg2Bo485VGlZeqAu",
	"1wUVy6EITtFuGjtcYRJmqRWZpKxmnK0QJDqWMQ1FvD0JtDJOvVTjfhTZqm8YW/pnCAW8GD0Eq5ZAuIXd",
	"EOo76GoNNkqim9x/94w/7NkXzfk0zFOrNIet5zWIoZKWX+wnv6ec3+NSGV02nxMhTNB0K0sePDprzDq/",
	"2DdHJUtWHhnX5R5j+DHGfdTyLaTvXIqaKb1I6kj0wilbbT32Ez0N/t+gpsFViKYhRetY5hokaoGDGywM",
	"y9gGna9uXI7Go3HTbbCKcfZf+pUxC7RKzcg3l/m6tVqXSN2n75GC037Q7GBftdylj1vJOPuKdFjQ0Sje",
	"Gu2rdv87Hldd14Q6FQdr16pKzp989W61AU6m7b1FcXisM4RRwFP0R9RtQyZfHKz4MCtn9Uro0v9lJRAe",
	"6XfY/zD+lP5zQE9XRu7OYv45wl2CFb5XAKpBJRew7G/IaY0pAYUPFXo/HGMu2i+FglvH4Saynuc5WDUi",
	"9OQQBIG1I2GKXGJh4mywclb+DQAA///rE360+wcAAA==",
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
