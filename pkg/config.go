package pkg

import (
	"fmt"
	"github.com/jinzhu/configor"
	"os"
)

type Config struct {
	Study   Study
	Modules []Module
}

type Study struct {
	Name    string
	Storage ConfigMap
}

type ModuleType string

const (
	ModuleTypeBinary ModuleType = "binary"
	ModuleTypeDocker ModuleType = "docker"
)

type ConfigMap map[string]any

func (cm *ConfigMap) GetString(name string) string {
	return (*cm)[name].(string)
}

func (cm *ConfigMap) GetStringSlice(name string) []string {
	if (*cm)[name] == nil {
		return nil
	}
	v := (*cm)[name].([]interface{})
	res := make([]string, len(v))
	for _, elem := range v {
		res = append(res, elem.(string))
	}
	return res
}

func (cm *ConfigMap) GetMap(name string) ConfigMap {
	return (*cm)[name].(ConfigMap)
}

func (cm *ConfigMap) GetStringMap(name string) map[string]string {
	var ok bool
	v := cm.GetMap(name)
	res := make(map[string]string)
	for k, elem := range v {
		if res[k], ok = elem.(string); !ok {
			res[k] = fmt.Sprintf("%v", elem)
		}
	}
	return res
}

type Module struct {
	Name          string
	Type          ModuleType
	Configuration ConfigMap
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}
	return cfg, configor.Load(cfg, os.Getenv("CONFIG_FILE"))
}
