package tracker

import (
	"crypto/sha1"
	"fmt"
	"io"
	"my-torrent/lib/encode"
	"net/http"
	"net/url"
)

type Client interface {
	Announce(args arguments) ([]byte, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) Announce(args arguments) ([]byte, error) {
	urls := buildAnnounceUrls(args)

	for _, tracker := range urls {
		address := tracker.String()
		fmt.Printf("Sending request to %s\n", address)

		response, err := sendRequest(address)
		if err != nil {
			fmt.Printf("Tracker %s responded with error = %s, trying next\n", tracker, err.Error())
			continue
		}

		fmt.Printf("Request to tracker %s successfuly sent!\n", address)

		return response, nil
	}

	return nil, nil
}

func sendRequest(address string) ([]byte, error) {
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
	return body, nil
}

func buildAnnounceUrls(args arguments) []*url.URL {
	urls := make([]*url.URL, 0)
	args.announcelist = append(args.announcelist, []string{args.announce})
	idSum := sha1.Sum([]byte(args.peerId))

	for _, announceList := range args.announcelist {
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
