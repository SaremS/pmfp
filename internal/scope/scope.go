package scope

import (
	"fmt"
)

type Scope interface {
	IsInScope(target string) bool
}

type ScopeFactory func() Scope

var ScopeRegistry = make(map[string]ScopeFactory)

func RegisterScope(typeName string, factory ScopeFactory) {
	if _, exists := ScopeRegistry[typeName]; exists {
		panic(fmt.Sprintf("Scope type '%s' is already registered", typeName))
	}
	ScopeRegistry[typeName] = factory
}
