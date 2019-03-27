package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Address         string `json:"address"`
	Workers         int    `json:"workers"`
	TaskCountOnPage int    `json:"task_count_on_page"`
}

func NewConfig(path string) (*Config, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	err = json.Unmarshal(file, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
