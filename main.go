package main

import (
	"os"

	"github.com/dominikus1993/kafka-topic-creator/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "bootstrapservers",
				Aliases: []string{"bs"},
				EnvVars: []string{"KAFKA_BROKER"},
			},
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				EnvVars: []string{"KAFKA_FILE"},
				Value:   "./kafka.yaml",
			},
		},
		Action: cmd.CreateTopicsIfNotExists,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	log.Infoln("done")
}
