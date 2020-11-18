package qhttp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

func (c *Client) PostForm(path string, vals map[string]interface{}) (int, []byte, error) {
	if _, ok := c.GetHeader("Content-Type"); !ok {
		c.AddHeader("Content-Type", "application/x-www-form-urlencoded")
	}
	urlStr := c.cfg.Schema + "://" + c.cfg.Host + path
	var vs = url.Values{}

	for k, v := range vals {
		vs.Add(k, fmt.Sprintf("%v", v))
	}
	var bufferBody io.Reader
	if len(vs) > 0 {
		bufferBody = bytes.NewBuffer([]byte(vs.Encode()))
	}
	req, err := http.NewRequest("POST", urlStr, bufferBody)
	if err != nil {
		return 500, nil, err
	}
	return c.DoRequest(req)
}

func (c *Client) GetWithMap(path string, vals map[string]interface{}) (int, []byte, error) {
	var vs = url.Values{}
	for k, v := range vals {
		vs.Add(k, fmt.Sprintf("%v", v))
	}
	return c.Get(path, vs)
}

func (c *Client) Get(path string, vs url.Values) (int, []byte, error) {
	for k, v := range c.query {
		vs.Add(k, v)
	}
	urlStr := c.cfg.Schema + "://" + c.cfg.Host + path
	if len(vs) > 0 {
		urlStr = "?" + vs.Encode()
	}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return 500, nil, err
	}
	return c.DoRequest(req)
}

func (c *Client) Post(path string, body []byte) (int, []byte, error) {
	if _, ok := c.GetHeader("Content-Type"); !ok {
		c.AddHeader("Content-Type", "application/json")
	}
	urlStr := c.cfg.Schema + "://" + c.cfg.Host + path
	if len(c.query) > 0 {
		vs := url.Values{}
		for k, v := range c.query {
			vs.Add(k, v)
		}
		urlStr = "?" + vs.Encode()
	}

	var bufferBody *bytes.Buffer
	if body == nil || len(body) <= 0 {
		bufferBody = bytes.NewBuffer([]byte("{}"))
	} else {
		bufferBody = bytes.NewBuffer(body)
	}
	req, err := http.NewRequest("POST", urlStr, bufferBody)
	if err != nil {
		return 500, nil, err
	}
	return c.DoRequest(req)

}

func (c *Client) DoRequest(req *http.Request) (int, []byte, error) {
	var (
		bb []byte
	)

	req.Host = c.cfg.Host
	req.URL.Host = c.cfg.Host
	req.URL.Scheme = c.cfg.Schema
	for k, v := range c.header {
		req.Header.Add(k, v)
	}
	req.Header.Set("Connection", "keep-alive")

	resp, err := c.HClient().Do(req)
	if err != nil {
		return 500, nil, err
	}
	defer resp.Body.Close()
	bb, err = ioutil.ReadAll(resp.Body)
	return resp.StatusCode, bb, err
}
