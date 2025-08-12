package scope

import (
	"sigs.k8s.io/yaml"
	"testing"
)

func TestExactMatchFromYAML(t *testing.T) {
	yamlInput := `
target_host: http://test.test
`
	var em ExactMatch
	err := yaml.Unmarshal([]byte(yamlInput), &em)

	if err != nil {
		t.Fatalf("Unmarshal failed with error: %v", err)
	}

	if em.TargetHost != "http://test.test" {
	t.Errorf("Execpted TargetHost to be 'http://test.test', go %s", em.TargetHost)
	}
}

func TestExactMatchIsInScope(t *testing.T) {
	exactMatch := &ExactMatch{
		TargetHost: "http://test.test",
	}

	if !exactMatch.IsInScope("http://test.test") {
		t.Fatalf("Hostname 'http://test.test' is expected to be in scope")
	}

	if exactMatch.IsInScope("http://test.fail") {
		t.Fatalf("Hostname 'http://test.fail' is expected to NOT be in scope")
	}
}

