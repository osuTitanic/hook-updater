package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	GitHub struct {
		Owner      string `json:"owner"`
		Repository string `json:"repository"`
	} `json:"github"`
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	ReleaseTypes []ReleaseType `json:"releaseTypes"`
	DataFolder   string
}

func (config *Config) ReleaseFolder() string {
	return filepath.Join(config.DataFolder, "releases")
}

type ReleaseType struct {
	Filename   string `json:"filename"`
	Identifier string `json:"identifier"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func EnsureConfigExists() error {
	_, err := os.Stat("config.json")
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return err
	}

	input, err := os.Open("config.example.json")
	if err != nil {
		return err
	}
	defer input.Close()

	output, err := os.Create("config.json")
	if err != nil {
		return err
	}
	defer output.Close()

	_, err = io.Copy(output, input)
	return err
}
