package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:        "Acker",
		Version:     "v0.0.1",
		HideVersion: true,
		Compiled:    time.Now(),
		HelpName:    "acker",
		Usage:       "A cli to consume/produce messages from/to AMQP servers, e.g. RabbitMQ",
		Commands: []*cli.Command{
			{
				Name:    "consume",
				Aliases: []string{"c"},
				Usage:   "Consume messages from the queue forever",
				Action: func(c *cli.Context) error {
					ConsumeForever(c.String("server"), c.String("channel"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "server",
						Value:   "amqp://guest:guest@localhost:5672/",
						Usage:   "RabbitMQ server address",
						Aliases: []string{"s"},
					},
					&cli.StringFlag{
						Name:     "channel",
						Value:    "",
						Usage:    "Queue channel name to consume from",
						Aliases:  []string{"c"},
						Required: true,
					},
				},
			}, {
				Name:    "produce",
				Aliases: []string{"p"},
				Usage:   "Produce a message to the queue",
				Action: func(c *cli.Context) error {
					Produce(c.String("server"), c.String("channel"), c.String("body"), c.Int("count"))
					return nil
				},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "server",
						Value:   "amqp://guest:guest@localhost:5672/",
						Usage:   "RabbitMQ server address",
						Aliases: []string{"s"},
					},
					&cli.StringFlag{
						Name:     "channel",
						Value:    "",
						Usage:    "Queue channel name to consume from",
						Aliases:  []string{"c"},
						Required: true,
					},
					&cli.StringFlag{
						Name:     "body",
						Value:    "",
						Usage:    "Body of message (as string)",
						Aliases:  []string{"b"},
						Required: true,
					},
					&cli.IntFlag{
						Name:    "count",
						Value:   1,
						Usage:   "Number of messages to produce",
						Aliases: []string{"n"},
					},
				},
			}, {
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "Shows version info",
				Action: func(c *cli.Context) error {
					fmt.Printf("version=%s runtime=%s\n", c.App.Version, runtime.Version())
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	FailOnError(err, "Failed to initialize app")
}
