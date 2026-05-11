package trackerclient

import (
	"crypto/sha1"
	"fmt"
	"io"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/internal/modelbuilder/torrent"
	"my-torrent/lib/constants"
	"my-torrent/lib/encode"
	"net/http"
	"net/url"
)

type arguments struct {
	peerId     string
	port       int
	uploaded   int64
	downloaded int64
	left       int64
	event      string
	infoHash   [20]byte
}

func GetPeers(torrent *torrent.Torrent, port int) *peers.Response {
	args := getArguments(torrent)
	args.port = port
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

		return response
	}

	return nil
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
	idSum := sha1.Sum([]byte(args.peerId))

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

			u.RawQuery = q.Encode() + "&info_hash=" + encode.Percent(args.infoHash[:]) + "&peer_id=" + encode.Percent(idSum[:])

			urls = append(urls, u)
		}
	}

	return urls
}

func getArguments(torrent *torrent.Torrent) arguments {
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
		peerId:     constants.PEER_ID,
		uploaded:   uploaded,
		downloaded: downloaded,
		left:       left,
		event:      event,
		infoHash:   torrent.Info.Hash,
	}

	return args
}
