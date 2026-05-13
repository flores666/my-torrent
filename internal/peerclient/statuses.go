package peerclient

const (
	PeerStatusDiscovered        = "discovered"         // получили от tracker
	PeerStatusConnecting        = "connecting"         // начали TCP connect
	PeerStatusConnected         = "connected"          // TCP connect успешен
	PeerStatusHandshakeSent     = "handshake_sent"     // отправили handshake
	PeerStatusHandshakeReceived = "handshake_received" // получили handshake
	PeerStatusReady             = "ready"              // handshake валиден, можно обмениваться сообщениями
	PeerStatusFailed            = "failed"             // ошибка connect/handshake/read/write
)
