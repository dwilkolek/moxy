package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-playground/validator/v10"
)

type MoxyConfig struct {
	Tunnel   TunnelConfig           `json:"tunnel" validate:"required"`
	Services map[string]ProxyConfig `json:"services" validate:"required"`
}

type TunnelConfig struct {
	UserAndHost      string `json:"userAndHost" validate:"required"`
	PathToPrivateKey string `json:"pathToPrivateKey" validate:"required"`
	Destination      string `json:"destination" validate:"required"`
}

type ProxyConfig struct {
	Port      int               `json:"port" validate:"required,numeric"`
	Headers   map[string]string `json:"headers"`
	AllowCors bool              `json:"allowCors" validate:"required,numeric"`
}

var Validator = validator.New()

func NewConfig(fileLocation string) (*MoxyConfig, error) {
	jsonFile, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		return &MoxyConfig{}, err
	}

	var objmap MoxyConfig
	json.Unmarshal([]byte(jsonFile), &objmap)
	err = Validator.Struct(objmap)
	if err != nil {
		return &MoxyConfig{}, err
	}
	return &objmap, nil
}
