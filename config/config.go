package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type MoxyConfig struct {
	Tunnel   TunnelConfig           `json:"tunnel"`
	Services map[string]ProxyConfig `json:"services"`
}

type TunnelConfig struct {
	UserAndHost      string `json:"userAndHost"`
	PathToPrivateKey string `json:"pathToPrivateKey"`
	Destination      string `json:"destination"`
}

type ProxyConfig struct {
	Port      int               `json:"port"`
	Headers   map[string]string `json:"headers"`
	AllowCors bool              `json:"allowCors"`
}

func NewConfig(arg string) (*MoxyConfig, error) {
	fileLocation := findFileLocation(arg)
	jsonFile, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return &MoxyConfig{}, err
	}
	var objmap MoxyConfig
	json.Unmarshal([]byte(jsonFile), &objmap)
	return &objmap, nil
}

func findFileLocation(arg string) string {
	configFile := arg
	if len(arg) > 0 {
		configFile = arg
	}
	log.Default().Printf("Config file: %s \n", configFile)
	return configFile

}
