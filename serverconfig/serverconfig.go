package serverconfig

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

//******************************************************************************

// Config type.
type Config struct {
	ServerPort          string `json:"server-port"`
	DataDirectory       string `json:"data-directory"`
	WireguardConfigPath string `json:"wireguard-config-path"`
	StaticDir           string `json:"static-dir"`
	LogFileAddress      string `json:"log-file-address"`
	EnableHTTPS         bool   `json:"enable-https"`
	CertAddress         string `json:"cert-address"`
	CertKeyAddress      string `json:"cert-key-address"`
}

//******************************************************************************

// LoadConfig loads the json config file
func LoadConfig(path string) (Config, error) {
	config := Config{}
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		return config, err
	}
	if len(fileData) <= 0 {
		return config, fmt.Errorf("config file is empty")
	}
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		return config, err
	}
	return config, nil
}

//******************************************************************************
