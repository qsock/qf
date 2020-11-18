package ws

// 封包
type Packet struct {
	// gorilla的类型
	t      int
	msg    []byte
	filter HandleFilterFunc
}
