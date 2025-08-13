package config

import (
	rm "github.com/sarems/pmfp/internal/request_middleware"
	"github.com/sarems/pmfp/internal/scope"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestConfigElementUnmarshalScope(t *testing.T) {
	yamlInput := `
scope:
- type: exact_match
  target_host: http://api.example.com
- type: exact_match
  target_host: http://api2.example.com
request_middlewares:
- type: add_header
  header_name: X-Test
  header_value: Test
`
	var cfg ConfigElement
	err := yaml.Unmarshal([]byte(yamlInput), &cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed with error: %v", err)
	}

	expectedLen := 2
	if len(cfg.Scope) != expectedLen {
		t.Fatalf("Expected %d scope checkers, but got %d", expectedLen, len(cfg.Scope))
	}

	exactScope, ok := cfg.Scope[0].(*scope.ExactMatch)
	if !ok {
		t.Fatalf("Expected first scope to be of type *scope.ExactMatch, but got %T", cfg.Scope[0])
	}

	expectedHost := "http://api.example.com"
	if exactScope.TargetHost != expectedHost {
		t.Errorf("Expected TargetHost to be '%s', but got '%s'", expectedHost, exactScope.TargetHost)
	}

	exactScope, ok = cfg.Scope[1].(*scope.ExactMatch)
	if !ok {
		t.Fatalf("Expected second scope to be of type *scope.ExactMatch, but got %T", cfg.Scope[1])
	}

	expectedHost = "http://api2.example.com"
	if exactScope.TargetHost != expectedHost {
		t.Errorf("Expected Prefix to be '%s', but got '%s'", expectedHost, exactScope.TargetHost)
	}

	//---Manipulators
	expectedLen = 1
	if len(cfg.RequestMiddlewares) != expectedLen {
		t.Fatalf("Expected %d request manipulators, but got %d", expectedLen, len(cfg.RequestMiddlewares))
	}

	exactManipulator, ok := cfg.RequestMiddlewares[0].(*rm.AddHeader)
	if !ok {
		t.Fatalf("Expected request manipulator to be of type *rm.AddHeader, but got %T", cfg.RequestMiddlewares[0])
	}

	expectedHeaderName := "X-Test"
	if exactManipulator.HeaderName != expectedHeaderName {
		t.Errorf("Expected HeaderName to be '%s', but got '%s'", expectedHeaderName, exactManipulator.HeaderName)
	}

	expectedHeaderValue := "Test"
	if exactManipulator.HeaderValue != expectedHeaderValue {
		t.Errorf("Expected HeaderValue to be '%s', but got '%s'", expectedHeaderValue, exactManipulator.HeaderValue)
	}

	if cfg.ProxyServer != nil {
		t.Errorf("Expected ProxyServer to be nil, but it isn't")
	}
}

func TestConfigElementUnmarshalScopeWithProxy(t *testing.T) {
	yamlInput := `
scope:
- type: exact_match
  target_host: http://api.example.com
- type: exact_match
  target_host: http://api2.example.com
request_middlewares:
- type: add_header
  header_name: X-Test
  header_value: Test
proxy_server: http://localhost:8000
`
	var cfg ConfigElement
	err := yaml.Unmarshal([]byte(yamlInput), &cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed with error: %v", err)
	}

	expectedLen := 2
	if len(cfg.Scope) != expectedLen {
		t.Fatalf("Expected %d scope checkers, but got %d", expectedLen, len(cfg.Scope))
	}

	exactScope, ok := cfg.Scope[0].(*scope.ExactMatch)
	if !ok {
		t.Fatalf("Expected first scope to be of type *scope.ExactMatch, but got %T", cfg.Scope[0])
	}

	expectedHost := "http://api.example.com"
	if exactScope.TargetHost != expectedHost {
		t.Errorf("Expected TargetHost to be '%s', but got '%s'", expectedHost, exactScope.TargetHost)
	}

	exactScope, ok = cfg.Scope[1].(*scope.ExactMatch)
	if !ok {
		t.Fatalf("Expected second scope to be of type *scope.ExactMatch, but got %T", cfg.Scope[1])
	}

	expectedHost = "http://api2.example.com"
	if exactScope.TargetHost != expectedHost {
		t.Errorf("Expected Prefix to be '%s', but got '%s'", expectedHost, exactScope.TargetHost)
	}

	//---Manipulators
	expectedLen = 1
	if len(cfg.RequestMiddlewares) != expectedLen {
		t.Fatalf("Expected %d request middlewares, but got %d", expectedLen, len(cfg.RequestMiddlewares))
	}

	exactManipulator, ok := cfg.RequestMiddlewares[0].(*rm.AddHeader)
	if !ok {
		t.Fatalf("Expected request middleware to be of type *rm.AddHeader, but got %T", cfg.RequestMiddlewares[0])
	}

	expectedHeaderName := "X-Test"
	if exactManipulator.HeaderName != expectedHeaderName {
		t.Errorf("Expected HeaderName to be '%s', but got '%s'", expectedHeaderName, exactManipulator.HeaderName)
	}

	expectedHeaderValue := "Test"
	if exactManipulator.HeaderValue != expectedHeaderValue {
		t.Errorf("Expected HeaderValue to be '%s', but got '%s'", expectedHeaderValue, exactManipulator.HeaderValue)
	}

	expectedHostname := "localhost"
	if cfg.ProxyServer.Hostname() != expectedHostname {
		t.Errorf("ProxyServer hostname expected to be '%s', but got '%s'", expectedHostname, cfg.ProxyServer.Hostname())
	}

	expectedPort := "8000"
	if cfg.ProxyServer.Port() != expectedPort {
		t.Errorf("ProxyServer port expected to be '%s', but got '%s'", expectedPort, cfg.ProxyServer.Port())
	}
}
