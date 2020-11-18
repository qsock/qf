package types

type LEVEL int8

const (
	ALL LEVEL = iota
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

var LevelText = map[LEVEL]string{
	ALL:     "ALL",
	DEBUG:   "DEBUG",
	INFO:    "INFO",
	WARNING: "WARNING",
	ERROR:   "ERROR",
	FATAL:   "FATAL",
}

const (
	MetaKey = "meta"
)
