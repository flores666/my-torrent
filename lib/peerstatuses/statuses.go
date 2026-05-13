package peerstatuses

const (
	Discovered        = "discovered"         // получили от tracker
	Connecting        = "connecting"         // начали TCP connect
	Connected         = "connected"          // TCP connect успешен
	HandshakeSent     = "handshake_sent"     // отправили handshake
	HandshakeReceived = "handshake_received" // получили handshake
	Ready             = "ready"              // handshake валиден, можно обмениваться сообщениями
	Failed            = "failed"             // ошибка connect/handshake/read/write
)
