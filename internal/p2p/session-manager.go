package p2p

import (
	"errors"
	"fmt"
	"my-torrent/internal/storage"
	"my-torrent/lib/constants"
	"my-torrent/lib/messages"
	"my-torrent/lib/peerstatuses"
	"slices"
	"sync"
)

type SessionManager interface {
	Start(session *Session) error
}

type sessionManager struct {
	sessions       []*Session
	mu             sync.Mutex
	torrentStorage storage.TorrentsStorage
}

func NewSessionManager(ts storage.TorrentsStorage) SessionManager {
	return &sessionManager{
		torrentStorage: ts,
		sessions:       make([]*Session, 0, constants.MAX_PEER_CONNECTIONS),
		mu:             sync.Mutex{},
	}
}

func (s *sessionManager) Start(session *Session) error {
	s.mu.Lock()
	if len(s.sessions) >= constants.MAX_PEER_CONNECTIONS {
		s.mu.Unlock()
		session.Close()
		return errors.New("max peer connections reached")
	}

	s.sessions = append(s.sessions, session)

	s.mu.Unlock()

	_ = s.torrentStorage.UpdatePeerStatus(session.InfoHash[:], session.Ip, session.Port, peerstatuses.Connected)

	go s.readLoop(session)

	return nil
}

func (s *sessionManager) readLoop(session *Session) {
	/*defer*/ func() {
		err := session.Close()
		if err != nil {
			fmt.Println(err)
		}

		s.mu.Lock()
		s.sessions = slices.DeleteFunc(s.sessions, func(item *Session) bool {
			return item == session
		})
		s.mu.Unlock()

		_ = s.torrentStorage.UpdatePeerStatus(session.InfoHash[:], session.Ip, session.Port, peerstatuses.Disconnected)
	}()

	for {
		msg, err := messages.ReadBody(session.Conn)
		if err != nil {
			return
		}

		_ = msg
	}
}
