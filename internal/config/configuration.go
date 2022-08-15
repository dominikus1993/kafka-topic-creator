package config

import "fmt"

type TopicConfig struct {
	Name        string `yaml:"name"`
	Partitions  int32  `yaml:"partitions"`
	Replication int16  `yaml:"replication"`
	Retention   string `yaml:"retention"`
}

func (cfg *TopicConfig) GetTopicName(prefix string) string {
	if prefix == "" {
		return cfg.Name
	}
	return fmt.Sprintf("%s.%s", prefix, cfg.Name)
}

type Configuration struct {
	Broker string `yaml:"broker" mapstructure:"KAFKA_BROKER"`
	Prefix string `yaml:"prefix" mapstructure:"PREFIX"`
	Topics []struct {
		Topic TopicConfig `yaml:"topic"`
	} `yaml:"topics"`
}
