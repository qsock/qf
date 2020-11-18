package qhttp

func (c *Client) AddHeader(k, v string) {
	c.header[k] = v
}

func (c *Client) RemoveHeader(k string) {
	delete(c.header, k)
}

func (c *Client) GetHeader(k string) (string, bool) {
	v, ok := c.header[k]
	return v, ok
}

func (c *Client) GetAllHeader() map[string]string {
	return c.header
}

func (c *Client) ClearAllHeader() {
	c.header = map[string]string{}
}
