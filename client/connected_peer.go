package client

import (
	"fmt"
	"net"
	"net/netip"
	"reflect"

	"github.com/bits-and-blooms/bitset"
	"github.com/yinfredyue/bittorrent-go/message"
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
	serving  *bitset.BitSet

	// peer state
	waitingForFstMsg bool // waiting for the first message after handshake
	choked           bool
}

func newHandshakeMsg(infoHash []byte) []byte {
	pstrLen := []byte{byte(pstrLen)}
	pstr := ([]byte)(pstr)
	reserved := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	peerId := NewPeerId()

	return util.ConcatBytes([]([]byte){pstrLen, pstr, reserved, infoHash, peerId})
}

func (p *connectedPeer) handleMsg(msg message.Msg) error {
	msgId := msg.Id()
	util.DPrintf("%v: %v", p.addrPort, msgId.String())

	switch msgId {
	case message.Choke:
		p.choked = true
	case message.Unchoke:
		p.choked = false
	case message.Interested:
		return fmt.Errorf("msgId: %v", msgId)
	case message.NotInterested:
		return fmt.Errorf("msgId: %v", msgId)
	case message.Have:
		msg := msg.(*message.HaveMsg)
		pieceIdx := msg.PieceIndex()
		if p.serving == nil {
			p.serving = bitset.New(524288)
		}
		p.serving.Set(uint(pieceIdx))
	case message.Bitfield:
		if !p.waitingForFstMsg {
			return fmt.Errorf("not waiting first message: unexpected bitfield message")
		}

		msg := msg.(*message.BitfieldMsg)
		bits := msg.Bits()
		p.serving = bitset.From(bits)
	case message.Request:
		return fmt.Errorf("msgId: %v", msgId)
	case message.Piece:
		return fmt.Errorf("msgId: %v", msgId)
	case message.Cancel:
		return fmt.Errorf("msgId: %v", msgId)
	default:
		return fmt.Errorf("unexpected case")
	}

	p.waitingForFstMsg = false
	return nil
}

// This function does the following:
// - Form a connection to the peer;
// - Handshake
// - Receive any Have and Bitfield message
// - Send Interested and wait handshake
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

	peer := connectedPeer{
		addrPort:         addrPort,
		conn:             conn,
		serving:          nil,
		waitingForFstMsg: true,
		choked:           true,
	}

	// send interest
	err = message.WriteMsg(peer.conn, &message.InterestedMsg{})
	if err != nil {
		return peer, err
	}

	// handle have/bitfield messages, and wait to be unchoked
	for peer.choked {
		msg, err := message.ReadMsg(conn)
		if err != nil {
			return peer, err
		}

		peer.handleMsg(msg)
	}

	util.DPrintf("%v unchoked and ready!\n", peer.addrPort)

	return peer, nil
}
