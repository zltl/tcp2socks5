// This file is part of tcp2socks5.
//
// tcp2socks5 is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// tcp2socks5 is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Foobar. If not, see <https://www.gnu.org/licenses/>.

package tcp2socks5

import (
	"context"
	"io"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

func Pipe(ctx context.Context, local, socks5, target string) error {
	log.Debugf("listening local: %s", local)
	lis, err := net.Listen("tcp", local)
	if err != nil {
		log.WithError(err).Fatal("cannot listen")
	}
	defer lis.Close()

	go func() {
		<-ctx.Done()
		l := lis.(*net.TCPListener)
		l.Accept()

		log.Debugf("listender %s cancel", local)
		l.SetDeadline(time.Now().Add(time.Microsecond))
	}()

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.WithError(err).Warn("cannot accept")
			if err, ok := err.(net.Error); ok && err.Timeout() {
				continue
			}
			return err
		}

		go func(conn net.Conn) {
			defer func() {
				conn.Close()
				log.Debugf("%s close", conn.RemoteAddr())
			}()
			log.Debugf("%s connected to local %s", conn.RemoteAddr(), local)

			var c net.Conn
			if len(socks5) > 0 {
				log.Debugf("connecting to socks5 %s for %s", socks5, conn.RemoteAddr())
				dailer, err := proxy.SOCKS5("tcp", socks5, nil, &net.Dialer{
					Timeout:   60 * time.Second,
					KeepAlive: 30 * time.Second,
				})
				if err != nil {
					log.WithError(err).Errorf("cannot initialize socks5 proxy %s", socks5)
					return
				}
				log.Debugf("dailing target %s with socks5 %s for %s", target,
					socks5, conn.RemoteAddr())
				c, err = dailer.Dial("tcp", target)
				if err != nil {
					log.WithError(err).WithField("target", target).Error("cannot dial")
					return
				}
			} else {
				c, err = net.Dial("tcp", target)
				if err != nil {
					log.WithError(err).WithField("target", target).Error("cannot dial")
					return
				}
			}
			defer c.Close()
			wg := &sync.WaitGroup{}
			wg.Add(2)
			go pipe(conn, c, wg)
			go pipe(c, conn, wg)
			wg.Wait()
		}(conn)
	}
}

// copy src -> dest
func pipe(src io.Reader, dst io.WriteCloser, wg *sync.WaitGroup) {
	defer wg.Done()
	io.Copy(dst, src)
}
