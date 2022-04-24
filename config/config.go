package config

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
)

type Config struct {
	Dbname     string `json:"dbname"`
	User       string `json:"user"`
	Password   string `json:"password"`
	Port       string `json:"port"`
	URL        string `json:"url"`
	Type       string `json:"type"`
	Secret     string `json:"secret"`
	Host       string `json:"host"`
	HostPort   string `json:"hostport"`
	Headers    bool   `json:"headers"`
	ExpireTime int    `json:"expireTime"`
}

var (
	configPath = flag.String("c", "", "path to config json file")
)

// New validates flag params and returns a new configuration
func New() (*Config, error) {
	if configPath == nil || len(*configPath) <= 0 {
		return nil, errors.New("no config path provided")
	}
	c, err := loadJSONConfig(*configPath)
	if err != nil {
		return nil, errors.New("load config")
	}

	return c, err
}

// loadJSONConfig loads file with passed path and parses JSON
func loadJSONConfig(p string) (*Config, error) {
	data, err := os.Open(p)
	if err != nil {
		return nil, errors.New("open json config file")
	}
	d := json.NewDecoder(data)
	var c Config
	if err := d.Decode(&c); err != nil {
		return nil, errors.New("parse json config")
	}
	return &c, nil
}
