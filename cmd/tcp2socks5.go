package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zltl/tcp2socks5"

	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "tcp2socks5",
		Usage: "Tunnel tcp port to socks5 proxy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "local",
				Value:   "0.0.0.0:4444",
				Usage:   "TCP Address to listen",
				Aliases: []string{"l"},
			},
			&cli.StringFlag{
				Name:    "socks5",
				Usage:   "Socks5 proxy address, example: 127.0.0.1:1080",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:     "target",
				Usage:    "Forwarding target, Example: www.google.com:80",
				Aliases:  []string{"t"},
				Required: true,
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "Set log level: trace, debug, info, warn, error",
				Value: "warn",
			},
		},
		Action: start,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func start(c *cli.Context) error {
	local := c.String("local")
	socks5 := c.String("socks5")
	target := c.String("target")

	level := c.String("log-level")
	lev, err := log.ParseLevel(level)
	if err != nil {
		log.WithField("error", err).Errorf("unkoiwn log-level: %s", level)
		return err
	}
	log.SetLevel(lev)
	log.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		DisableQuote:    true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339Nano,
	})

	return tcp2socks5.Pipe(context.TODO(), local, socks5, target)
}
