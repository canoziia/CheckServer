package config

import (
	"fmt"

	ini "gopkg.in/ini.v1"
)

// Config 配置文件
type Config map[string]map[string]string

// GetConfig 获取conf.ini内容
func GetConfig(configPath string) (config Config) {
	config = make(Config, 0)
	var cfg *(ini.File)
	var err error
	cfg, err = ini.Load(configPath)
	if err != nil {
		panic(fmt.Sprintf("Load Config %s Failed: %s", configPath, err.Error()))
	}

	sectionSlice := cfg.SectionStrings()
	for _, sectionStr := range sectionSlice {
		keySlice := cfg.Section(sectionStr).KeyStrings()
		sectionMap := make(map[string]string, 0)
		for _, key := range keySlice {
			sectionMap[key] = cfg.Section(sectionStr).Key(key).String()
		}
		config[sectionStr] = sectionMap
	}
	return config
}
