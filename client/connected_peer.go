package client

import (
	"fmt"
	"net"
	"net/netip"
	"reflect"

	"github.com/yinfredyue/bittorrent-go/util"
)

const (
	pstr            = "BitTorrent protocol"
	pstrLen         = len(pstr)
	handshakeMsgLen = 49 + pstrLen
)

type connectedPeer struct {
	addrPort netip.AddrPort
	conn     net.Conn
}

func newHandshakeMsg(infoHash []byte) []byte {
	pstrLen := []byte{byte(pstrLen)}
	pstr := ([]byte)(pstr)
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	peerId := NewPeerId()

	return util.ConcatBytes([]([]byte){pstrLen, pstr, reserved, infoHash, peerId})
}

// Create a TCP connection and finish handshake
func connectToPeer(addrPort netip.AddrPort, infoHash []byte) (connectedPeer, error) {
	conn, err := net.Dial("tcp", addrPort.String())
	if err != nil {
		return connectedPeer{}, err
	}

	// send handshake message
	handshakeMsg := newHandshakeMsg(infoHash)
	numBytesWritten, err := conn.Write(handshakeMsg)
	if err != nil {
		return connectedPeer{}, err
	}
	if numBytesWritten != handshakeMsgLen {
		return connectedPeer{}, fmt.Errorf("only written %v bytes, but expect %v bytes", numBytesWritten, handshakeMsgLen)
	}

	// receive handshake response
	respBytes := make([]byte, handshakeMsgLen)
	numBytesRead, err := conn.Read(respBytes)
	if err != nil {
		return connectedPeer{}, err
	}
	if numBytesRead != handshakeMsgLen {
		return connectedPeer{}, fmt.Errorf("only read %v bytes, but expect %v bytes", numBytesRead, handshakeMsgLen)
	}

	pstrLenByte := respBytes[0]
	if pstrLenByte != 19 || string(respBytes[1:20]) != pstr || !reflect.DeepEqual(respBytes[28:48], infoHash) {
		return connectedPeer{}, fmt.Errorf("malformed handshake response?")
	}

	return connectedPeer{addrPort: addrPort, conn: conn}, nil
}
