package main

import (
	"os"

	"github.com/dominikus1993/kafka-topic-creator/cmd"
	"github.com/dominikus1993/kafka-topic-creator/internal/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v2"
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

	app := &cli.App{
		Action: cmd.CreateTopicsIfNotExists(conf),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	log.Infoln("done")
}
