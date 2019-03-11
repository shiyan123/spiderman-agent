package app

import (
	"encoding/json"
	"os"
	"spiderman-agent/common/model"
	"spiderman-agent/utils"
)

type Config struct {
	Server       *ServerConfig
	CenterServer *CenterServerConfig
}

type ServerConfig struct {
	Address      string
	Port         int
	RegisterIp   string
	RegisterName string
}

type CenterServerConfig struct {
	URL string
	Ips []string
}

func LoadConfig(configPath string) (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func InitCenterConfig(url string, config *Config) (err error) {
	resp, err := utils.HttpGet(url)
	if err != nil {
		return
	}
	var ips model.CenterIpsResp
	err = json.Unmarshal(resp, &ips)
	if err != nil {
		return
	}
	config.CenterServer.Ips = ips.IpList
	return
}

func ConvertConfig(config *Config, params *ProcParams) (err error) {
	config.Server.RegisterIp = params.RegisterIp
	config.Server.RegisterName = params.RegisterName
	return
}
