package constants

import "time"

const P2P_MESSAGE_TYPES_COUNT = 12

const MAX_PEER_CONNECTIONS = 30

const (
	MAX_DOWNLOAD_SLOTS     = 15
	MAX_IN_FLIGHT_PER_PEER = 5
	BLOCK_SIZE             = 16 * 1024
)

const (
	MAX_UPLOAD_SLOTS = 4
	MAX_UPLOAD_RATE  = 5 * 1024 * 1024 // 5 MB/s TODO
)

const (
	READ_TIMEOUT  = 5 * time.Second
	WRITE_TIMEOUT = 5 * time.Second
	TIMEOUT       = 5 * time.Second
)

const KEEP_ALIVE_PERIOD_MINUTES = 2 * time.Minute
