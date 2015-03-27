package fhirterm

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	HttpPort               int        `json:"http_port"`
	HttpHost               string     `json:"http_host"`
	HttpCorsAllowedOrigins []string   `json:"http_cors_allowed_origins"`
	Databases              []string   `json:"databases"`
	Storage                JsonObject `json:"storage"`
}

func ReadConfig(path string) (*Config, error) {
	file, e := ioutil.ReadFile(path)
	if e != nil {
		return nil, e
	}

	var cfg Config
	e = json.Unmarshal(file, &cfg)

	if e != nil {
		return nil, fmt.Errorf("Could not parse config file: %s", e)
	}

	return &cfg, nil
}
