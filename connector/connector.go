package connector

import (
	"context"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"time"
)

// NewConnection init connector using either web or unix socket
func NewConnection(network, address string) *http.Client {
	tr := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			conn, err := net.Dial(network, address)
			if err != nil {
				log.Fatal(err)
			}
			return conn, nil
		},
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr, Timeout: 10 * time.Second}
}
