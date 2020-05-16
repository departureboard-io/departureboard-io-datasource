package main

import (
	"encoding/json"
	"fmt"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

// DataSourceSettings are the configurable settings for the departureboard.io data source.
type DataSourceSettings struct {
	APIEndpoint string `json:"apiEndpoint"`
	APIKey      string `json:"apiKey"`
}

// LoadSettings gets the DataSourceSettings from the plugin context.
func LoadSettings(ctx backend.PluginContext) (DataSourceSettings, error) {
	settings := DataSourceSettings{}

	if ctx.DataSourceInstanceSettings == nil {
		return settings, fmt.Errorf("error reading settings: instance settings are nil")
	}
	err := json.Unmarshal(ctx.DataSourceInstanceSettings.JSONData, &settings)
	if err != nil {
		return settings, fmt.Errorf("error reading settings: %s", err.Error())
	}
	settings.APIKey = ctx.DataSourceInstanceSettings.DecryptedSecureJSONData["apiKey"]

	return settings, nil
}
