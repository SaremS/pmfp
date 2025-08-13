package config

import (
	"fmt"
	"net/http"
	"net/url"

	"encoding/json"
)

type Config struct {
	ConfigElements map[string]ConfigElement
	ProxyServer *url.URL
}

func (c *Config) ApplyRequestMiddleware(request *http.Request) {
	for _, config := range c.ConfigElements {
		config.ApplyRequestMiddleware(request)
	}
}

func (c *Config) UnmarshalJSON(data []byte) error {
	type ConfigAlias struct {
		ConfigElements []ConfigElement `json:"config_elements"`
		ProxyServer	*string	`json:"proxy_server,omitempty"`
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

	c.ConfigElements = make(map[string]ConfigElement)

	for _, config := range alias.ConfigElements {
		name := config.Name	

		if _, exists := c.ConfigElements[name]; exists {
			return fmt.Errorf("Config Elements must have unique names; element with name '%s' already exists", name)
		}

		c.ConfigElements[name] = config
	}
	
	return nil
}
