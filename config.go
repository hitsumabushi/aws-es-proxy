package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config provide config file struct
type Config struct {
	ServerMap map[string]endpoint `json:"server_map"`
	ListPath string `json:"list_path"`
}

type endpoint struct {
	Region string `json:"region"`
	Host   string `json:"host"`
}

func (e *endpoint) String() string {
	return fmt.Sprintf("region=%s,host=%s", e.Region, e.Host)
}

func loadConfig(file string) Config {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		panic(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
