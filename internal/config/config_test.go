package config

import (
	"sigs.k8s.io/yaml"
	"testing"
)


func TestConfigUnmarshal(t *testing.T) {
	yamlInput := `
config_elements:
- name: element_1 
  scope:
  - type: exact_match
    target_host: http://api.example.com
  request_middlewares:
  - type: add_header
    header_name: X-Test
    header_value: Test
- name: element_2
  scope:
  - type: exact_match
    target_host: http://api2.example.com
  request_middlewares:
  - type: add_header
    header_name: X-Test-2
    header_value: Test2
`

	var cfg Config
	err := yaml.Unmarshal([]byte(yamlInput), &cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed with error: %v", err)
	}
}
