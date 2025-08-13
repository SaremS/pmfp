package request_middleware

import (
	"net/http"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestAddHeaderFromYAML(t *testing.T) {
	yamlInput := `
header_name: X-Custom-Header
header_value: my-value
`

	var ah AddHeader
	err := yaml.Unmarshal([]byte(yamlInput), &ah)

	if err != nil {
		t.Fatalf("Unmarshal failed with error: %v", err)
	}

	if ah.HeaderName != "X-Custom-Header" {
		t.Errorf("Expected HeaderName to be 'X-Custom-Header', got '%s'", ah.HeaderName)
	}

	if ah.HeaderValue != "my-value" {
		t.Errorf("Expected HeaderValue to be 'my-value', got '%s'", ah.HeaderValue)
	}
}

func TestAddHeaderApply(t *testing.T) {
	request, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	headerModifier := &AddHeader{
		HeaderName:  "Content-Type",
		HeaderValue: "application/json",
	}

	headerModifier.Apply(request)

	if request.Header.Get(headerModifier.HeaderName) != headerModifier.HeaderValue {
		t.Errorf("Expected header %s to be '%s', but got '%s'",
			headerModifier.HeaderName, headerModifier.HeaderValue, request.Header.Get(headerModifier.HeaderName))
	}

	if len(request.Header) != 1 {
		t.Errorf("Expected only one header to be set, but found %d", len(request.Header))
	}
}
