package main

import (
	"github.com/Shopify/sarama"
	"github.com/dominikus1993/kafka-topic-creator/internal/config"
	"github.com/dominikus1993/kafka-topic-creator/internal/kafka"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("kafka")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.MergeInConfig()
	err := viper.ReadInConfig()
	if err != nil {
		log.WithError(err).Fatalf("unable to read config")
	}
	if err := viper.BindEnv("KAFKA_BROKER"); err != nil {
		log.WithError(err).Fatal("unable to bind env kafka_broker")
	}
	if err := viper.BindEnv("PREFIX"); err != nil {
		log.WithError(err).Fatal("unable to bind env prefix")
	}
	conf := &config.Configuration{}
	err = viper.Unmarshal(conf)
	if err != nil {
		log.WithError(err).Fatal("unable to decode into struct")
	}
	log.WithField("config", conf).Info("configuration")
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin([]string{conf.Broker}, config)
	if err != nil {
		log.WithError(err).Fatal("unable to create cluster admin")
	}

	topics, err := admin.ListTopics()
	if err != nil {
		log.WithError(err).Fatal("unable to list topics")
	}
	creator := kafka.NewKafkaTopicCreator(conf, admin, topics)

	for _, topicToCreate := range conf.Topics {
		err := creator.CreateTopicIfNotExists(topicToCreate.Topic)
		if err != nil {
			log.WithError(err).Fatal("unable to create topic")
		}
	}
	log.Infoln("done")

}
