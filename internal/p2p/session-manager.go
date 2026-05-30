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
	"time"
)

type SessionManager interface {
	Download(session *Session) error
	Upload(session *Session) error
}

type keepAliveFunc struct {
	ticker *time.Ticker
	f      func()
}

type sessionManager struct {
	dowloads          []*Session
	dm                sync.Mutex
	uploads           []*Session
	um                sync.Mutex
	keepAliveRequests map[*Session]keepAliveFunc
	km                sync.Mutex
	torrentStorage    storage.TorrentsStorage
}

func NewSessionManager(ts storage.TorrentsStorage) SessionManager {
	return &sessionManager{
		torrentStorage:    ts,
		uploads:           make([]*Session, 0, constants.MAX_UPLOAD_SLOTS),
		dowloads:          make([]*Session, 0, constants.MAX_DOWNLOAD_SLOTS),
		keepAliveRequests: make(map[*Session]keepAliveFunc),
		um:                sync.Mutex{},
		dm:                sync.Mutex{},
	}
}

func (s *sessionManager) Upload(session *Session) error {
	s.um.Lock()
	if len(s.uploads) >= constants.MAX_UPLOAD_SLOTS {
		s.um.Unlock()
		session.Close()
		return errors.New("max peer upload connections reached")
	}

	s.uploads = append(s.uploads, session)

	s.um.Unlock()

	_ = s.torrentStorage.UpdatePeerStatus(session.InfoHash[:], session.Ip, session.Port, peerstatuses.Connected)

	go s.uploadLoop(session)

	return nil
}

func (s *sessionManager) Download(session *Session) error {
	s.dm.Lock()
	if len(s.dowloads) >= constants.MAX_DOWNLOAD_SLOTS {
		s.dm.Unlock()
		// todo: do not close to remember them
		session.Close()
		return errors.New("max peer download connections reached")
	}

	s.dowloads = append(s.dowloads, session)

	s.dm.Unlock()

	_ = s.torrentStorage.UpdatePeerStatus(session.InfoHash[:], session.Ip, session.Port, peerstatuses.Connected)

	go s.downloadLoop(session)

	return nil
}

func (s *sessionManager) uploadLoop(session *Session) {
	/*defer*/ func() {
		err := session.Close()
		if err != nil {
			fmt.Println(err)
		}

		s.um.Lock()
		s.uploads = slices.DeleteFunc(s.uploads, func(item *Session) bool {
			return item == session
		})
		s.um.Unlock()

		_ = s.torrentStorage.UpdatePeerStatus(session.InfoHash[:], session.Ip, session.Port, peerstatuses.Disconnected)
	}()

	stateMachine(session)
}

func (s *sessionManager) downloadLoop(session *Session) {
	/*defer*/ func() {
		err := session.Close()
		if err != nil {
			fmt.Println(err)
		}

		s.dm.Lock()
		s.dowloads = slices.DeleteFunc(s.dowloads, func(item *Session) bool {
			return item == session
		})
		s.dm.Unlock()

		_ = s.torrentStorage.UpdatePeerStatus(session.InfoHash[:], session.Ip, session.Port, peerstatuses.Disconnected)
	}()

	s.keepAlive(session)
	stateMachine(session)
}

func stateMachine(session *Session) {
	for {
		select {
		case <-session.Done:
			fmt.Printf("Session with peer %s is closed gracefuly, stopping state machine\n", session.PeerId)
			return
		default:
			msg, err := messages.ReadBody(session.Conn)
			if err != nil {
				return
			}

			_ = msg
		}
	}
}

func (s *sessionManager) sendMessage(session *Session, message []byte) {
	s.keepAlive(session)

	_, err := session.Write(message)
	if err != nil {
		fmt.Println("Error while sending message, closing connection")
		session.Close()
	}
}

func (s *sessionManager) keepAlive(session *Session) {
	s.km.Lock()
	if item, ok := s.keepAliveRequests[session]; !ok {
		t := time.NewTicker(constants.KEEP_ALIVE_PERIOD_MINUTES)

		s.keepAliveRequests[session] = keepAliveFunc{
			f: func() {
				defer t.Stop()

				for {
					select {
					case <-session.Done:
						s.km.Lock()
						delete(s.keepAliveRequests, session)
						s.km.Unlock()
						return
					case <-t.C:
						fmt.Printf("Sending keep alive message to %s", session.Conn.RemoteAddr())

						_, err := session.Write(make([]byte, 4))
						if err != nil {
							session.Close()
							fmt.Printf("Error while sending keep alive message to %s, session is closing", session.Conn.RemoteAddr())
						}
					}
				}
			},
			ticker: t,
		}

		s.keepAliveRequests[session].f()
	} else {
		item.ticker.Reset(constants.KEEP_ALIVE_PERIOD_MINUTES)
	}

	s.km.Unlock()
}
