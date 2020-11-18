package ws

type HandleMessageFunc func(*Session, []byte)
type HandleErrorFunc func(*Session, error)
type HandleCloseFunc func(*Session, int, string)
type HandleSessionFunc func(*Session)
type HandleFilterFunc func(*Session) bool
type HandleEmptyFunc func()
