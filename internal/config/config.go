package config

import (
	"fmt"

	"encoding/json"

	rm "github.com/sarems/pmfp/internal/request_manipulation"
	"github.com/sarems/pmfp/internal/scope"
)

type ScopeArray []scope.Scope
type ManipulatorArray []rm.RequestManipulator

type Config struct {
	Scope       ScopeArray       
	Manipulators ManipulatorArray 
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type ConfigAlias struct {
		Scope []json.RawMessage `json:"scope"`
		Manipulators []json.RawMessage `json:"request_manipulators"`
	}

	var alias ConfigAlias
	if err := json.Unmarshal(data, &alias); err != nil {
		return err
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


