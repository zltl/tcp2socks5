// This file is part of tcp2socks5.
//
// tcp2socks5 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// any later version.
//
// tcp2socks5 is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/zltl/tcp2socks5"

	cli "github.com/urfave/cli/v2"
)

var (
	VERSION = "-"
)

func main() {
	app := &cli.App{
		Name:    "tcp2socks5",
		Usage:   "Tunnel tcp port to socks5 proxy",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{
				Name:  "liaotonglang",
				Email: "liaotonglang@gmail.com",
			},
		},
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
				Usage:    "Forwarding target, example: www.google.com:80",
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
