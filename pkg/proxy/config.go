package proxy

import (
	"gopkg.in/yaml.v2"
)

// FIXME: group configurations, currently certificate config and goproxy specific config are mixed like a mess
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

func ParseLocalServerCfg(data []byte) (LocalConfig, error) {
	result := LocalConfig{}
	err := yaml.Unmarshal(data, &result)
	return result, err
}

type RemoteServerCfg struct {
	Addr       string   `yaml:"addr"`
	CaCert     string   `yaml:"caCert"`
	ServerCert string   `yaml:"serverCert"`
	Domains    []string `yaml:"domains"`
	ServerKey  string   `yaml:"serverKey"`
	Protocols  []string `yaml:"protocols"`
}

func ParseRemoteServerCfg(data []byte) (*RemoteServerCfg, error) {
	result := &RemoteServerCfg{}
	err := yaml.Unmarshal(data, result)
	return result, err
}
