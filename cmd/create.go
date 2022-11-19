package cmd

import (
	"github.com/Shopify/sarama"
	"github.com/dominikus1993/kafka-topic-creator/internal/config"
	"github.com/dominikus1993/kafka-topic-creator/internal/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func CreateTopicsIfNotExists(conf *config.Configuration) func(context *cli.Context) error {
	return func(context *cli.Context) error {
		config := sarama.NewConfig()
		admin, err := sarama.NewClusterAdmin([]string{conf.Broker}, config)
		if err != nil {
			log.WithError(err).Errorln("unable to create cluster admin")
			return cli.Exit("unable to create cluster admin", 1)
		}

		topics, err := admin.ListTopics()
		if err != nil {
			log.WithError(err).Errorln("unable to list topics")
			return cli.Exit("unable to list topics", 1)
		}
		creator := kafka.NewKafkaTopicCreator(conf, admin, topics)

		for _, topicToCreate := range conf.Topics {
			err := creator.CreateTopicIfNotExists(topicToCreate.Topic)
			if err != nil {
				log.WithError(err).WithField("topic", topicToCreate).Errorln("unable to create topic")
			}
		}

		return cli.Exit("topic created", 0)
	}
}
