package requests

import (
	"fmt"
	"io"
	"log"
	"my-torrent/internal/peers"
	"my-torrent/internal/torrent"
	"my-torrent/internal/utils/encoding"
	"net"
	"net/http"
	"net/url"

	"github.com/gofrs/uuid/v5"
)

var peerId string

type arguments struct {
	peerId     string
	port       int
	uploaded   int64
	downloaded int64
	left       int64
	event      string
	infoHash   string
}

func GetPeers(torrent *torrent.Torrent) (*peers.Response, net.Listener) {
	args, listener := getArguments(torrent)
	urls := buildUrls(torrent.Announce, torrent.AnnounceList, args)

	for _, tracker := range urls {
		address := tracker.String()
		fmt.Printf("Sending request to %s\n", address)

		response, err := sendRequest(address)
		if err != nil {
			fmt.Printf("Tracker %s responded with error = %s, trying next\n", tracker, err.Error())
			continue
		}

		fmt.Printf("Request to tracker %s successfuly sent!\n", address)

		return response, listener
	}

	return nil, nil
}

func sendRequest(address string) (*peers.Response, error) {
	httpResponse, err := http.Get(address)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != http.StatusOK {
		//todo handle status code
		return nil, fmt.Errorf("status code is %d\n", httpResponse.StatusCode)
	}

	body, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	defer httpResponse.Body.Close()

	response, err := peers.Parse(body)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func buildUrls(announce string, announceTierList [][]string, args arguments) []*url.URL {
	urls := make([]*url.URL, 0)
	announceTierList = append(announceTierList, []string{announce})

	for _, announceList := range announceTierList {
		for _, a := range announceList {
			u, err := url.Parse(a)
			if err != nil {
				continue
			}

			q := u.Query()
			q.Add("port", fmt.Sprintf("%d", args.port))
			q.Add("uploaded", fmt.Sprintf("%d", args.uploaded))
			q.Add("downloaded", fmt.Sprintf("%d", args.downloaded))
			q.Add("left", fmt.Sprintf("%d", args.left))
			q.Add("event", args.event)

			u.RawQuery = q.Encode() + "&info_hash=" + encoding.Sha1(args.infoHash) + "&peer_id=" + encoding.Sha1(peerId)

			urls = append(urls, u)
		}
	}

	return urls
}

func getArguments(torrent *torrent.Torrent) (arguments, net.Listener) {
	peerIdUuid, _ := uuid.NewV7()
	peerId = fmt.Sprintf("%x", peerIdUuid)[:20]

	port, listener := getFreePort()
	if listener == nil {
		log.Fatalln("No available ports")
	}

	event := "started" // stopped completed

	var uploaded int64 = 0
	var downloaded int64 = 0

	left := torrent.Info.Length
	if left == 0 {
		for _, f := range torrent.Info.Files {
			left = f.Length
		}
	}

	args := arguments{
		port:       port,
		peerId:     peerId,
		uploaded:   uploaded,
		downloaded: downloaded,
		left:       left,
		event:      event,
		infoHash:   torrent.Info.Hash,
	}

	return args, listener
}

func getFreePort() (int, net.Listener) {
	var listener net.Listener
	var err error

	const start int = 6881
	const end int = 6889

	for port := start; port <= end; port++ {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))

		if err == nil {
			fmt.Println("Using port:", port)
			return port, listener
		}
	}

	return 0, nil
}
