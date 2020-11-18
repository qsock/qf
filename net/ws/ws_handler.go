package ws

func (w *Server) HandleMessage(fn HandleMessageFunc) {
	w.messageHandler = fn
}
func (w *Server) HandleMessageBinary(fn HandleMessageFunc) {
	w.messageHandlerBinary = fn
}
func (w *Server) HandleSent(fn HandleMessageFunc) {
	w.messageSentHandler = fn
}
func (w *Server) HandleSentBinary(fn HandleMessageFunc) {
	w.messageSentHandlerBinary = fn
}
func (w *Server) HandleError(fn HandleErrorFunc) {
	w.errorHandler = fn
}
func (w *Server) HandleClose(fn HandleCloseFunc) {
	w.closeHandler = fn
}
func (w *Server) HandleConnect(fn HandleSessionFunc) {
	w.connectHandler = fn
}

func (w *Server) HandleDisconnect(fn HandleSessionFunc) {
	w.disconnectHandler = fn
}
func (w *Server) HandleShutdown(fn HandleEmptyFunc) {
	w.shutdownHandler = fn
}
