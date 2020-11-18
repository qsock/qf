package ws

import (
	"errors"
)

var (
	ErrSessionClosed = errors.New("qws:session has closed yet!")
	ErrSessionFulled = errors.New("qws:session buffer is full!")
	ErrServerDown    = errors.New("qws:server has down!")
)
