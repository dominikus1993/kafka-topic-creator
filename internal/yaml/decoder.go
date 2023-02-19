package yaml

import (
	"github.com/dominikus1993/kafka-topic-creator/internal/config"
	"gopkg.in/yaml.v3"
)

func DecodeConfiguration(data *[]byte) (*config.Configuration, error) {
	var result config.Configuration
	err := yaml.Unmarshal(*data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
