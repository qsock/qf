package qetcd

import "errors"

var (
	ErrConf          = errors.New("conf error")
	ErrHasRegistered = errors.New("address has registed")
)
