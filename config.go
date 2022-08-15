package main

import "fmt"

type TopicConfig struct {
	Name        string `yaml:"name"`
	Partitions  int32  `yaml:"partitions"`
	Replication int16  `yaml:"replication"`
	Retention   string `yaml:"retention"`
}

func (cfg *TopicConfig) GetTopicName(envPrefix string) string {
	if envPrefix == "" {
		return cfg.Name
	}
	return fmt.Sprintf("%s.%s", envPrefix, cfg.Name)
}

type Configuration struct {
	Env    string `yaml:"env"`
	Topics []struct {
		Topic TopicConfig `yaml:"topic"`
	} `yaml:"topics"`
}
