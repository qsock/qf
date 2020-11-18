package qhttp

func (c *Client) AddQuery(k, v string) {
	c.query[k] = v
}

func (c *Client) RemoveQuery(k string) {
	delete(c.query, k)
}

func (c *Client) GetQuery(k string) (string, bool) {
	v, ok := c.query[k]
	return v, ok
}

func (c *Client) GetAllQuery() map[string]string {
	return c.query
}

func (c *Client) ClearAllQuery() {
	c.query = map[string]string{}
}
