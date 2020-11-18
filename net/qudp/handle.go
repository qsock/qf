package qudp

type HandleMessageFunc func(*Session, []byte)
type HandleErrorFunc func(*Session, error)
type HandleSessionFunc func(*Session)
