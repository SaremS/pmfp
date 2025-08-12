package request_manipulation

import (
	"fmt"
	"net/http"
)

type RequestManipulator interface {
	Apply(request *http.Request)
}

type RequestManipulatorFactory func() RequestManipulator

var RequestManipulatorRegistry = make(map[string]RequestManipulatorFactory)

func RegisterRequestManipulator(typeName string, factory RequestManipulatorFactory) {
	if _, exists := RequestManipulatorRegistry[typeName]; exists {
		panic(fmt.Sprintf("RequestManipulator type '%s' is already registered", typeName))
	}
	RequestManipulatorRegistry[typeName] = factory
}
