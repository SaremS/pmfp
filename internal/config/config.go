package config

import (
	"fmt"
	"net/http"
	"net/url"

	"encoding/json"

	rm "github.com/sarems/pmfp/internal/request_manipulation"
	"github.com/sarems/pmfp/internal/scope"
)

type ScopeArray []scope.Scope
type ManipulatorArray []rm.RequestManipulator

type Config struct {
	Scope        ScopeArray
	Manipulators ManipulatorArray
	ProxyServer  *url.URL
}

func (c *Config) isInScope(target string) bool {
	for _, scoper := range c.Scope {
		if scoper.IsInScope(target) {
			return true
		}
	}
	return false
}

func (c *Config) ApplyManipulation(request *http.Request) {
	target := request.URL.Host

	if c.isInScope(target) {
		for _, manipulator := range c.Manipulators {
			manipulator.Apply(request)
		}
	}
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type ConfigAlias struct {
		Scope        []json.RawMessage `json:"scope"`
		Manipulators []json.RawMessage `json:"request_manipulators"`
		ProxyServer  *string           `json:"proxy_server,omitempty"`
	}

	var alias ConfigAlias
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

	for _, rawManipulator := range alias.Manipulators {
		var peeker struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(rawManipulator, &peeker); err != nil {
			return fmt.Errorf("failed to decode manipulator type: %w", err)
		}

		factory, found := rm.RequestManipulatorRegistry[peeker.Type]
		if !found {
			return fmt.Errorf("unknown request manipulator type: '%s'", peeker.Type)
		}

		requestManipulator := factory()

		if err := json.Unmarshal(rawManipulator, requestManipulator); err != nil {
			return err
		}
		c.Manipulators = append(c.Manipulators, requestManipulator)
	}
	return nil
}
