package config

import (
	"fmt"
	"net/http"
	"net/url"

	"encoding/json"

	rm "github.com/sarems/pmfp/internal/request_middleware"
	"github.com/sarems/pmfp/internal/scope"
)

type ScopeArray []scope.Scope
type RequestMiddlewareArray []rm.RequestMiddleware

type ConfigElement struct {
	Name               string
	Scope              ScopeArray
	RequestMiddlewares RequestMiddlewareArray
	ProxyServer        *url.URL
}

func (c *ConfigElement) isInScope(target string) bool {
	for _, scoper := range c.Scope {
		if scoper.IsInScope(target) {
			return true
		}
	}
	return false
}

func (c *ConfigElement) ApplyRequestMiddleware(request *http.Request) {
	target := request.URL.Host

	if c.isInScope(target) {
		for _, requestMiddleware := range c.RequestMiddlewares {
			requestMiddleware.Apply(request)
		}
	}
}

func (c *ConfigElement) UnmarshalJSON(data []byte) error {
	type ConfigElementAlias struct {
		Name               string            `json:"name"`
		Scope              []json.RawMessage `json:"scope"`
		RequestMiddlewares []json.RawMessage `json:"request_middlewares"`
		ProxyServer        *string           `json:"proxy_server,omitempty"`
	}

	var alias ConfigElementAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
	}

	if alias.ProxyServer != nil {
		parsedUrl, err := url.Parse(*alias.ProxyServer)
		if err != nil {
			return fmt.Errorf("failed to decode proxy server: %w", err)
		}
		c.ProxyServer = parsedUrl
	}

	c.Name = alias.Name

	for _, rawScope := range alias.Scope {
		var peeker struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawScope, &peeker); err != nil {
			return fmt.Errorf("failed to decode scope type: %w", err)
		}

		factory, found := scope.ScopeRegistry[peeker.Type]
		if !found {
			return fmt.Errorf("unknown scope type: '%s'", peeker.Type)
		}

		scope := factory()

		if err := json.Unmarshal(rawScope, scope); err != nil {
			return err
		}
		c.Scope = append(c.Scope, scope)
	}

	for _, rawMiddleware := range alias.RequestMiddlewares {
		var peeker struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawMiddleware, &peeker); err != nil {
			return fmt.Errorf("failed to decode manipulator type: %w", err)
		}

		factory, found := rm.RequestMiddlewareRegistry[peeker.Type]
		if !found {
			return fmt.Errorf("unknown request manipulator type: '%s'", peeker.Type)
		}

		requestMiddleware := factory()

		if err := json.Unmarshal(rawMiddleware, requestMiddleware); err != nil {
			return err
		}
		c.RequestMiddlewares = append(c.RequestMiddlewares, requestMiddleware)
	}
	return nil
}
