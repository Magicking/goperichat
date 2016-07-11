package main

import (
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/urfave/cli"
)

func init() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	go func(){
		select {
			case <-interrupt:
				log.Fatal("Exit")
		}
	}()
}

func startNewListener(id string, c chan LineEntry) *PeriscopeChatListener {
	pm, err := GetPeriscopeMeta(id)
	if err != nil {
		log.Fatal(err)
	}
	cm, err := GetPeriscopeMetaChat(pm.ChatToken)
	if err != nil {
		log.Fatal(err)
	}
	pcl := NewPeriscopeChatListener(*pm, *cm)
	err = pcl.Run(c)
	if err != nil {
		log.Fatal(err)
	}
	return pcl
}

func main() {
	app := cli.NewApp()
	app.Name = "GoPerichat"
	app.Usage = "No usage"

	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name:	"id",
			Usage:	"Stream ID",
			EnvVar:	"STREAM_ID",
		},
	}

	app.Action = func (c* cli.Context) error {
		if c.String("id") == "" {
			return errors.New("No id given")
		}
		sl := NewStatsLive(WordFreq)
		startNewListener(c.String("id"), sl.C)
		return nil
	}
	app.Run(os.Args)

	// For MQTT
	// Connect to MQTT
	// Register Topic
	// On new message launch new recorder
	/*
	opts.AddBroker(server)
	opts.SetUsername(username)
	opts.SetPassword(password)
	*/
}
