package message

import (
	"encoding/binary"
	"fmt"
	"net"

	"github.com/yinfredyue/bittorrent-go/util"
)

const (
	numBytesForLen    = 4
	startOffsetForLen = 0

	numBytesForMsgId    = 1
	startOffsetForMsgId = 4

	startOffsetForPayload = 5
)

// This doesn't include handshake message.
type Msg interface {
	Id() MsgId
	Len() uint
	Payload() []byte
}

func bytesOf(msg Msg) ([]byte, error) {
	totalLen := msg.Len() + 4
	buffer := make([]byte, totalLen)

	// length
	lengthBytes, err := util.EncodeUint32(uint32(msg.Len()), numBytesForLen)
	if err != nil {
		return []byte{}, err
	}
	copy(buffer[startOffsetForLen:], lengthBytes)

	// msg id
	msgIdBytes := []byte{byte(msg.Id())}
	copy(buffer[startOffsetForMsgId:], msgIdBytes)

	// payload
	copy(buffer[startOffsetForPayload:], msg.Payload())

	return buffer, nil
}

func WriteMsg(conn net.Conn, msg Msg) error {
	msgBytes, err := bytesOf(msg)
	if err != nil {
		return err
	}

	numBytesWritten, err := conn.Write(msgBytes)
	if err != nil {
		return err
	}
	if numBytesWritten != len(msgBytes) {
		return fmt.Errorf("expect to write full msg %v bytes, but wrote %v bytes", len(msgBytes), numBytesWritten)
	}

	return nil
}

func ReadMsg(conn net.Conn) (Msg, error) {
	lengthBytes := make([]byte, 4)
	numBytesRead, err := conn.Read(lengthBytes)
	if numBytesRead != 4 {
		return nil, fmt.Errorf("expect to read 4 bytes for length, but read %v bytes", numBytesRead)
	}
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBytes)
	if length == 0 {
		return &KeepAliveMsg{}, nil
	}

	msgIdByte := make([]byte, 1)
	numBytesRead, err = conn.Read(msgIdByte)
	if numBytesRead != 1 {
		return nil, fmt.Errorf("expect to read 1 byte for msg id, but read %v bytes", numBytesRead)
	}
	if err != nil {
		return nil, err
	}

	payload := make([]byte, length-1)
	numBytesRead, err = conn.Read(payload)
	if numBytesRead != int(length)-1 {
		return nil, fmt.Errorf("expect to read %v byte for msg id, but read %v bytes", length-1, numBytesRead)
	}
	if err != nil {
		return nil, err
	}

	switch msgId := MsgIdFrom(msgIdByte[0]); msgId {
	case Choke:
		return nil, fmt.Errorf("msgId: %v", msgId)
	case Unchoke:
		return &UnchokeMsg{}, nil
	case Interested:
		return nil, fmt.Errorf("msgId: %v", msgId)
	case NotInterested:
		return nil, fmt.Errorf("msgId: %v", msgId)
	case Have:
		msg, err := newHaveMsg(payload)
		if err != nil {
			return nil, err
		}
		return &msg, nil
	case Bitfield:
		msg := newBitfieldMsg(payload)
		return &msg, nil
	case Request:
		return nil, fmt.Errorf("msgId: %v", msgId)
	case Piece:
		return nil, fmt.Errorf("msgId: %v", msgId)
	case Cancel:
		return nil, fmt.Errorf("msgId: %v", msgId)
	default:
		return nil, fmt.Errorf("unexpected MsgId")
	}
}
