package cmd

import (
	"strings"

	"github.com/Shopify/sarama"
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

	topics, err := admin.ListTopics()
	if err != nil {
		log.WithError(err).Errorln("unable to list topics")
		return cli.Exit("unable to list topics", 1)
	}
	creator := kafka.NewKafkaTopicCreator(configuration, admin, topics)

	for _, topicToCreate := range configuration.Topics {
		err := creator.CreateTopicIfNotExists(topicToCreate.Topic)
		if err != nil {
			log.WithError(err).WithField("topic", topicToCreate).Errorln("unable to create topic")
		}
	}

	return cli.Exit("topic created", 0)
}
