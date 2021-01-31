package proxy

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type LocalConfig struct {
	Global  LocalGlobalConfig `json:"global,omitempty" yaml:"global"`
	Servers []LocalServerCfg  `json:"servers,omitempty" yaml:"servers"`
}

type LocalGlobalConfig struct {
	MetricsAddr string `json:"metrics_addr,omitempty" yaml:"metrics_addr"`
}

type LocalServerCfg struct {
	Name       string `yaml:"name"`
	Protocol   string `yaml:"protocol"`
	LocalAddr  string `yaml:"localAddr"`
	ServerAddr string `yaml:"serverAddr"`
	CaCert     string `yaml:"caCert"`
	ClientCert string `yaml:"clientCert"`
	ClientKey  string `yaml:"clientKey"`
}

func ParseLocalServerCfg(file string) (LocalConfig, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return LocalConfig{}, err
	}

	result := LocalConfig{}
	err = yaml.Unmarshal(data, &result)
	return result, err
}

type RemoteServerCfg struct {
	Addr       string   `yaml:"addr"`
	CaCert     []string `yaml:"caCert"`
	ServerCert string   `yaml:"serverCert"`
	ServerKey  string   `yaml:"serverKey"`
	Protocols  []string `yaml:"protocols"`
}

func ParseRemoteServerCfg(file string) (*RemoteServerCfg, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	result := &RemoteServerCfg{}
	err = yaml.Unmarshal(data, result)
	return result, err
}
