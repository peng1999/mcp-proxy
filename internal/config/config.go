package config

import (
	"gopkg.in/yaml.v3"
	"os"

	//
	"mcp-proxy/api"
)

// replaceDefaults TODO
func replaceDefaults(config *api.Configuration) {

	if config.Server.Options.CacheThresholdBytes == 0 {
		config.Server.Options.CacheThresholdBytes = api.DefaultCacheThresholdBytes
	}

	if config.Server.Options.PaginationDefaultPageSize == 0 {
		config.Server.Options.PaginationDefaultPageSize = api.DefaultPaginationDefaultPageSize
	}

	if config.Server.Options.PaginationMaxPageSize == 0 {
		config.Server.Options.PaginationMaxPageSize = api.DefaultPaginationMaxPageSize
	}
}

// Marshal TODO
func Marshal(config api.Configuration) (bytes []byte, err error) {
	bytes, err = yaml.Marshal(config)
	return bytes, err
}

// Unmarshal TODO
func Unmarshal(bytes []byte) (config api.Configuration, err error) {
	err = yaml.Unmarshal(bytes, &config)
	return config, err
}

// ReadFile TODO
func ReadFile(filepath string) (config api.Configuration, err error) {
	var fileBytes []byte
	fileBytes, err = os.ReadFile(filepath)
	if err != nil {
		return config, err
	}

	// Expand environment variables present in the config
	// This will cause expansion in the following way: field: "$FIELD" -> field: "value_of_field"
	fileExpandedEnv := os.ExpandEnv(string(fileBytes))

	config, err = Unmarshal([]byte(fileExpandedEnv))

	replaceDefaults(&config)

	return config, err
}
