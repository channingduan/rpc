package config

import (
	"encoding/json"
	"io/ioutil"
)

func Register(path string) (*Config, error) {

	cfg, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return parseJson(cfg)
}

func parseJson(cfg []byte) (*Config, error) {

	var err error
	var config Config
	err = json.Unmarshal(cfg, &config)

	return &config, err
}
