package connector

import (
	"context"
	"net"
	"net/http"
	"time"
)

// NewConnection init connector using either web or unix socket
func NewConnection(network, address string) *http.Client {
	tr := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial(network, address)
		},
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}

	return &http.Client{Transport: tr, Timeout: 10 * time.Second}
}
