package main

import (
	"errors"
	"strings"

	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func getKafkaBrokers(env string) ([]string, error) {
	brokers := viper.GetString(env)
	if brokers == "" {
		return make([]string, 0), errors.New("brokers urls are empty")
	}
	kafkaBrokers := strings.Split(brokers, ",")
	return kafkaBrokers, nil
}

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.BindEnv("KAFKA_BROKERS")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	conf := &Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.WithError(err).Fatal("unable to decode into struct")
	}
	log.WithField("config", conf).Info("configuration")
	config := sarama.NewConfig()

	brokers, err := getKafkaBrokers("KAFKA_BROKERS")
	if err != nil {
		log.WithError(err).Fatalln("KAFKA_BROKERS should not be empty")
	}
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		log.WithError(err).Fatal("unable to create cluster admin")
	}

	topics, err := admin.ListTopics()
	if err != nil {
		log.WithError(err).Fatal("unable to list topics")
	}

	for _, topicToCreate := range conf.Topics {
		topicName := topicToCreate.Topic.GetTopicName(conf.Env)
		if _, ok := topics[topicName]; ok {
			log.WithField("topic", topicName).Infoln("topic already exists")
			continue
		}

		topicDetail := &sarama.TopicDetail{
			NumPartitions:     topicToCreate.Topic.Partitions,
			ReplicationFactor: topicToCreate.Topic.Replication,
		}

		if topicToCreate.Topic.Retention != "" {
			topicDetail.ConfigEntries = map[string]*string{"retention.ms": &topicToCreate.Topic.Retention}
		}

		err = admin.CreateTopic(topicName, topicDetail, false)

		if err != nil {
			log.WithError(err).Fatal("unable to create topic")
			continue
		}

		log.WithField("topic", topicName).Infoln("topic created")
	}
	log.Infoln("done")

}
