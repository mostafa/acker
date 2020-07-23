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
		Version:     "v0.0.3",
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
					ConsumeForever(
						c.String("server"),
						c.String("queue"),
						c.Bool("existing_queue"),
						c.Bool("autoack"),
						c.Bool("recover"),
						c.Bool("current-consumer"))
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
						Name:     "queue",
						Value:    "",
						Usage:    "Queue name to consume from",
						Aliases:  []string{"q"},
						Required: true,
					},
					&cli.BoolFlag{
						Name:     "existing-queue",
						Value:    true,
						Usage:    "Connect to an existing queue",
						Aliases:  []string{"e"},
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "autoack",
						Value:    true,
						Usage:    "Automatically acknowledges messages upon consumption",
						Aliases:  []string{"a"},
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "recover",
						Value:    false,
						Usage:    "Recover nack messages on the queue before consumption",
						Aliases:  []string{"r"},
						Required: false,
					},
					&cli.BoolFlag{
						Name:     "current-consumer",
						Value:    false,
						Usage:    "Recover nack messages on the queue before consumption in this CLI consumer",
						Aliases:  []string{"u"},
						Required: false,
					},
				},
			}, {
				Name:    "produce",
				Aliases: []string{"p"},
				Usage:   "Produce a message to the queue",
				Action: func(c *cli.Context) error {
					Produce(
						c.String("server"),
						c.String("queue"),
						c.Bool("existing_queue"),
						c.String("body"),
						c.Int("count"))
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
						Name:     "queue",
						Value:    "",
						Usage:    "Queue name to produce to",
						Aliases:  []string{"q"},
						Required: true,
					},
					&cli.BoolFlag{
						Name:     "existing-queue",
						Value:    true,
						Usage:    "Connect to an existing queue",
						Aliases:  []string{"e"},
						Required: false,
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
