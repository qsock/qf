package qsyslog

type Config struct {
	Network string `toml:"network" json:"network"`
	Addr    string `toml:"addr" json:"addr"`
	Tag     string `toml:"tag" json:"tag"`
}
