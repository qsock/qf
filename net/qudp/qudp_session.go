package qudp

func (w *Server) HandleMessage(fn HandleMessageFunc) {
	w.messageHandler = fn
}

func (w *Server) HandleSent(fn HandleMessageFunc) {
	w.messageSentHandler = fn
}

func (w *Server) HandleError(fn HandleErrorFunc) {
	w.errorHandler = fn
}

func (w *Server) HandleConnect(fn HandleSessionFunc) {
	w.connectHandler = fn
}
