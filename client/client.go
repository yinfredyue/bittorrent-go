package client

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/netip"
	"reflect"
	"strconv"

	"github.com/jackpal/bencode-go"
	"github.com/yinfredyue/bittorrent-go/torrent"
	"github.com/yinfredyue/bittorrent-go/util"
)

type Client struct {
	torrent torrent.Torrent
	peers   []*connectedPeer
}

func loadPeerAddrPorts(t torrent.Torrent) ([]netip.AddrPort, error) {
	// send GET request to tracker
	infoHash, err := t.InfoHash()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("GET", t.Tracker, nil)
	if err != nil {
		return nil, err
	}

	params := req.URL.Query()
	params.Add("info_hash", string(infoHash))
	params.Add("peer_id", string(NewPeerId()))
	params.Add("port", "6881")
	params.Add("uploaded", "0")
	params.Add("downloaded", "0")
	params.Add("left", strconv.Itoa(t.Info.Length))
	params.Add("compact", "1")
	req.URL.RawQuery = params.Encode()

	cli := http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	decodedObj, err := bencode.Decode(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	decodedDict, ok := decodedObj.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("fail to decode tracker response")
	}

	var peers []netip.AddrPort
	switch peersRaw := decodedDict["peers"].(type) {
	case string:
		// Compact
		// Each peer is represented by 6 bytes.
		// First 4 bytes is IP, where each byte is a number in the IP.
		// Last 2 bytes is port, in big-endian order.
		for i := 0; i < len(peersRaw); i += 6 {
			addr, ok := netip.AddrFromSlice([]byte(peersRaw)[i : i+4])
			if !ok {
				return nil, fmt.Errorf("fail to parse peer addr")
			}
			port := binary.BigEndian.Uint16([]byte(peersRaw[i+4 : i+6]))
			addrPort := netip.AddrPortFrom(addr, port)
			peers = append(peers, addrPort)
		}
	case [](interface{}):
		// Not compact
		// Each peer is represented as a dict.
		for _, peerRaw := range peersRaw {
			peerRawDict := peerRaw.(map[string]interface{})
			ipStr := peerRawDict["ip"].(string)
			addr, err := netip.ParseAddr(ipStr)
			if err != nil {
				return nil, err
			}
			port := peerRawDict["port"].(int64)
			addrPort := netip.AddrPortFrom(addr, uint16(port))
			peers = append(peers, addrPort)
		}
	default:
		log.Fatalf("Unexpected case: %v", reflect.TypeOf(peersRaw))
	}

	return peers, nil
}

func (cli *Client) ConnectToPeers() error {
	peerAddrPorts, err := loadPeerAddrPorts(cli.torrent)
	if err != nil {
		return err
	}

	peers := make([]*connectedPeer, 0)
	infoHash, err := cli.torrent.InfoHash()
	if err != nil {
		return err
	}

	peerChan := make(chan *connectedPeer)
	defer close(peerChan)
	for _, addrPort := range peerAddrPorts {
		go func(addrPort netip.AddrPort) {
			peer, err := connectToPeer(addrPort, infoHash)
			if err != nil {
				util.DPrintf("cannot connect to peer @ %v, err: %v", addrPort, err)
				peerChan <- nil
			} else {
				peerChan <- &peer
			}
		}(addrPort)
	}

	for i := 0; i < len(peerAddrPorts); i++ {
		if peer := <-peerChan; peer != nil {
			peers = append(peers, peer)
		}
	}

	cli.peers = peers
	log.Printf("Connected to %v peers!", len(peers))
	return nil
}

func NewClient(torrent torrent.Torrent) Client {
	return Client{torrent: torrent}
}

func Download() error {
	return nil
}
