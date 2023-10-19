package kafka

import (
	"errors"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/dominikus1993/kafka-topic-creator/internal/config"
	log "github.com/sirupsen/logrus"
)

type KafkaTopicCreator struct {
	config *config.Configuration
	admin  sarama.ClusterAdmin
}

func NewKafkaTopicCreator(config *config.Configuration, admin sarama.ClusterAdmin) *KafkaTopicCreator {
	return &KafkaTopicCreator{
		config: config,
		admin:  admin,
	}
}

func (creator *KafkaTopicCreator) CreateTopicsIfNotExists() error {
	var err error = nil

	brokerTopics, err := creator.admin.ListTopics()
	if err != nil {
		log.WithError(err).Errorln("unable to list topics")
		return err
	}
	for _, topicToCreate := range creator.config.Topics {
		topic := topicToCreate.Topic
		topicName := topic.GetTopicName(creator.config.Prefix)

		if _, ok := brokerTopics[topicName]; ok {
			log.WithField("topic", topicName).Infoln("topic already exists")
		}

		err := creator.createTopicIfNotExists(topic)
		if err != nil {
			err = errors.Join(err, errors.New(fmt.Sprintf("unable to create topic: %s", topic.Name)))
		}
	}
	return err
}

func (creator *KafkaTopicCreator) createTopicIfNotExists(topic config.TopicConfig) error {
	topicName := topic.GetTopicName(creator.config.Prefix)

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     topic.Partitions,
		ReplicationFactor: topic.Replication,
	}

	if topic.Retention != "" {
		topicDetail.ConfigEntries = map[string]*string{"retention.ms": &topic.Retention}
	}

	err := creator.admin.CreateTopic(topicName, topicDetail, false)

	if err != nil {
		return fmt.Errorf("unable to create topic: %w", err)
	}

	log.WithField("topic", topicName).Infoln("topic created")
	return nil
}
