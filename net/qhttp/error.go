package qhttp

import "errors"

var (
	ErrNoClient      = errors.New("no http client")
	ErrConfig        = errors.New("error no config")
	ErrAllServerDown = errors.New("all server down")
	ErrClientExists  = errors.New("client exists")
)
