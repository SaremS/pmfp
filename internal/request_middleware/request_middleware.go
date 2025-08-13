package request_middleware

import (
	"fmt"
	"net/http"
)

type RequestMiddleware interface {
	Apply(request *http.Request)
}

type RequestMiddlewareFactory func() RequestMiddleware

var RequestMiddlewareRegistry = make(map[string]RequestMiddlewareFactory)

func RegisterRequestMiddleware(typeName string, factory RequestMiddlewareFactory) {
	if _, exists := RequestMiddlewareRegistry[typeName]; exists {
		panic(fmt.Sprintf("RequestMiddleware type '%s' is already registered", typeName))
	}
	RequestMiddlewareRegistry[typeName] = factory
}
