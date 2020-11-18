package qhttp

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

var (
	x map[string]*Client = make(map[string]*Client)
)

type Client struct {
	ch     chan *http.Client
	cfg    *Config
	header map[string]string
	query  map[string]string
}

//  添加一个http client
func Add(config *Config) error {
	if !config.check() {
		return ErrConfig
	}
	_, ok := x[config.Name]
	if ok {
		return ErrClientExists
	}
	client := initClient(config)
	for i := 0; i < config.PoolSize; i++ {
		client.ch <- &http.Client{
			Timeout: config.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   20,
				IdleConnTimeout:       60 * time.Second,
				DisableCompression:    true,
				ResponseHeaderTimeout: 10 * time.Second,
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					for _, addr := range config.Addrs {
						c, err := net.DialTimeout("tcp", addr, time.Second)
						if err != nil {
							continue
						}
						return c, nil
					}
					return nil, ErrAllServerDown
				},
			},
		}
	}
	x[config.Name] = client

	return nil
}

func initClient(config *Config) *Client {
	client := new(Client)
	client.cfg = config
	client.ch = make(chan *http.Client, config.PoolSize)
	client.header = make(map[string]string)

	return client
}

func GetClient(name string) (*Client, error) {
	c, ok := x[name]
	if !ok {
		return nil, ErrNoClient
	}
	return c, nil
}

func (c *Client) HClient() *http.Client {
	client := <-c.ch
	c.ch <- client
	return client
}
