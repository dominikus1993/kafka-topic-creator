package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/dominikus1993/kafka-topic-creator/internal/config"
	log "github.com/sirupsen/logrus"
)

type KafkaTopicCreator struct {
	config *config.Configuration
	admin  sarama.ClusterAdmin
	topics map[string]sarama.TopicDetail
}

func NewKafkaTopicCreator(config *config.Configuration, admin sarama.ClusterAdmin, topics map[string]sarama.TopicDetail) *KafkaTopicCreator {
	return &KafkaTopicCreator{
		config: config,
		topics: topics,
		admin:  admin,
	}
}

func (creator *KafkaTopicCreator) CreateTopicIfNotExists(topic config.TopicConfig) error {
	topicName := topic.GetTopicName(creator.config.Prefix)
	if _, ok := creator.topics[topicName]; ok {
		log.WithField("topic", topicName).Infoln("topic already exists")
		return nil
	}

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
