package main

import (
	"io"
	"net"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"

	cli "github.com/urfave/cli/v2"
)

// copy src -> dest
func pipe(src io.Reader, dst io.WriteCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	defer dst.Close()
	io.Copy(dst, src)
}

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
				Value:   "127.0.0.1:1080",
				Usage:   "Socks5 proxy address",
				Aliases: []string{"s"},
			},
			&cli.StringFlag{
				Name:     "target",
				Usage:    "Forwarding target, Example: www.google.com:80",
				Aliases:  []string{"t"},
				Required: true,
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
	lis, err := net.Listen("tcp", local)
	if err != nil {
		log.WithError(err).Fatal("cannot listen")
	}
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.WithError(err).Warn("cannot accept")
		}
		go func(conn net.Conn) {
			defer conn.Close()
			dailer, err := proxy.SOCKS5("tcp", socks5, nil, &net.Dialer{
				Timeout:   60 * time.Second,
				KeepAlive: 30 * time.Second,
			})
			if err != nil {
				log.WithError(err).Warn("cannot initialize socks5 proxy")
				return
			}
			c, err := dailer.Dial("tcp", target)
			if err != nil {
				log.WithError(err).WithField("target", target).Warn("cannot dial")
				return
			}
			wg := &sync.WaitGroup{}
			wg.Add(2)
			go pipe(conn, c, wg)
			go pipe(c, conn, wg)
			wg.Wait()
		}(conn)
	}
}
