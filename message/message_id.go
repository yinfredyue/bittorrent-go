package message

import "fmt"

type MsgId int

const (
	Choke MsgId = iota
	Unchoke
	Interested
	NotInterested
	Have
	Bitfield
	Request
	Piece
	Cancel
	KeepAlive = 99
)

func MsgIdFrom(v uint8) MsgId {
	switch v {
	case 0:
		return Choke
	case 1:
		return Unchoke
	case 2:
		return Interested
	case 3:
		return NotInterested
	case 4:
		return Have
	case 5:
		return Bitfield
	case 6:
		return Request
	case 7:
		return Piece
	case 8:
		return Cancel
	default:
		panic(fmt.Sprintf("Unexpected msg id: %v", v))
	}
}

func (id MsgId) String() string {
	switch id {
	case Choke:
		return "Choke"
	case Unchoke:
		return "Unchoke"
	case Interested:
		return "Interested"
	case NotInterested:
		return "NotInterested"
	case Have:
		return "Have"
	case Bitfield:
		return "Bitfield"
	case Request:
		return "Request"
	case Piece:
		return "Piece"
	default:
		panic("unexpected msgId")
	}
}
