package tracker

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"my-torrent/internal/bencode"
	"my-torrent/internal/modelbuilder/peers"
	"my-torrent/lib/encode"
	"net/http"
	"net/url"
)

type Client interface {
	Announce(args arguments) (*AnnounceResponse, error)
}

type client struct {
}

func NewClient() Client {
	return &client{}
}

func (c *client) Announce(args arguments) (*AnnounceResponse, error) {
	urls := buildAnnounceUrls(args)

	for _, tracker := range urls {
		address := tracker.String()
		fmt.Printf("Sending request to %s\n", address)

		body, err := sendRequest(address)
		if err != nil {
			fmt.Printf("Tracker %s responded with error = %s, trying next\n", tracker, err.Error())
			continue
		}

		fmt.Printf("Request to tracker %s successfuly sent!\n", address)

		response, err := parseAnnounceResponse(body)
		if err != nil {
			return nil, nil
		}

		return response, nil
	}

	return nil, nil
}

func parseAnnounceResponse(response []byte) (*AnnounceResponse, error) {
	fmt.Println("Decoding peers response")

	res, err := bencode.Decode(response)
	if err != nil {
		return nil, err
	}

	if len(res) <= 0 {
		return nil, errors.New("parsed bencode contains 0 elements")
	}

	fmt.Println("Building peers response struct")

	if dict, ok := res[0].Value.(bencode.BDict); ok {
		return buildAnnounceResponse(dict)
	}

	return nil, errors.New("could not read response")
}

func buildAnnounceResponse(root bencode.BDict) (*AnnounceResponse, error) {
	response := AnnounceResponse{}
	for k, v := range root {
		switch k {
		case intervalKey:
			if val, err := bencode.GetTypedValue[int64](v); err == nil {
				response.Interval = int(val)
			}
		case peersKey:
			if val, err := bencode.GetTypedValue[[]any](v); err == nil {
				response.Peers, err = peers.BuildPeers(val)
				if err != nil {
					return nil, err
				}
			}
		default:
			return nil, fmt.Errorf("could not determine field with name = %s\n", k)
		}
	}

	return &response, nil
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
