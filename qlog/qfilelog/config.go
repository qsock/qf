package qfilelog

import "path"

type Config struct {
	FileSize int    `toml:"filesize" json:"filesize"`
	FileNum  int    `toml:"filenum" json:"filenum"`
	FileName string `toml:"filename" json:"filename"`
	Dir      string `toml:"dir" json:"dir"`
	UseGzip  bool   `toml:"use_gzip" json:"use_gzip"`

	fileName string
	fileSize int64
}

func DefaultConfig() Config {
	c := Config{}
	c.FileSize = 256
	c.FileNum = 50
	c.FileName = "info"
	c.UseGzip = true
	c.Dir = "./logs"
	return c
}

func (c *Config) check() {
	if c.FileSize == 0 {
		c.FileSize = 128
	}
	if c.FileNum == 0 {
		c.FileNum = 10
	}
	if c.FileName == "" {
		c.FileName = "INFO"
	}
	if c.Dir == "" {
		c.Dir = "./logs"
	}

	c.fileSize = int64(c.FileSize * 1024 * 1024)
	c.fileName = path.Join(c.Dir, c.FileName+".log")
}
