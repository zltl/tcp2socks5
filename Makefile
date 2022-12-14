# This file is part of tcp2socks5.
#
# tcp2socks5 is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# any later version.
#
# tcp2socks5 is distributed in the hope that it will be useful, but
# WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU
# General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with Foobar. If not, see <https://www.gnu.org/licenses/>.

.PHONY: tcp2socks5 clean

VERSION="1.0.0"

tcp2socks5:
	go build -ldflags '-X main.VERSION=${VERSION}' ./cmd/tcp2socks5.go

clean:
	rm -rf tcp2socks5
