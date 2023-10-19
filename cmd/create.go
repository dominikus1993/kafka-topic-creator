package cmd

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/dominikus1993/kafka-topic-creator/internal/kafka"
	"github.com/dominikus1993/kafka-topic-creator/internal/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func CreateTopicsIfNotExists(context *cli.Context) error {
	broker := context.String("bootstrapservers")
	yamlFilePath := context.String("file")
	yamlFile, err := yaml.LoadYamlFile(yamlFilePath)
	if err != nil {
		return cli.Exit("can't load kafka configuratrion yaml file", 1)
	}
	configuration, err := yaml.DecodeConfiguration(yamlFile)
	if err != nil {
		return cli.Exit("can't parse kafka configuratrion yaml file", 1)
	}

	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(strings.Split(broker, ","), config)
	if err != nil {
		log.WithError(err).Errorln("unable to create cluster admin")
		return cli.Exit("unable to create cluster admin", 1)
	}
	creator := kafka.NewKafkaTopicCreator(configuration, admin)

	err = creator.CreateTopicsIfNotExists()
	if err != nil {
		return cli.Exit(err, 1)
	}
	return cli.Exit("topic created", 0)
}
