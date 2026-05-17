package peerstatuses

const (
	Discovered        = "discovered"         // получили от tracker
	Disconnected      = "disconnected"       // закончили TCP connect
	Connected         = "connected"          // TCP connect успешен
	HandshakeReceived = "handshake_received" // получили handshake
	Ready             = "ready"              // handshake валиден, можно обмениваться сообщениями
	Failed            = "failed"             // ошибка connect/handshake/read/write
)
