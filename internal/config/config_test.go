package config

import (
	"github.com/sarems/pmfp/internal/scope"
	rm "github.com/sarems/pmfp/internal/request_manipulation"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestConfigUnmarshalScope(t *testing.T) {
	yamlInput := `
scope:
- type: exact_match
  target_host: http://api.example.com
- type: exact_match
  target_host: http://api2.example.com
request_manipulators:
- type: add_header
  header_name: X-Test
  header_value: Test
`
	var cfg Config
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
	if len(cfg.Manipulators) != expectedLen {
		t.Fatalf("Expected %d request manipulators, but got %d", expectedLen, len(cfg.Manipulators))
	}

	exactManipulator, ok := cfg.Manipulators[0].(*rm.AddHeader)
	if !ok {
		t.Fatalf("Expected request manipulator to be of type *rm.AddHeader, but got %T", cfg.Manipulators[0])
	}

	expectedHeaderName := "X-Test"
	if exactManipulator.HeaderName != expectedHeaderName {
		t.Errorf("Expected HeaderName to be '%s', but got '%s'", expectedHeaderName, exactManipulator.HeaderName)
	}

	expectedHeaderValue := "Test"
	if exactManipulator.HeaderValue != expectedHeaderValue {
		t.Errorf("Expected HeaderValue to be '%s', but got '%s'", expectedHeaderValue, exactManipulator.HeaderValue)
	}
}
