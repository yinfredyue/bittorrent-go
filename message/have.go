package message

import (
	"encoding/binary"
	"fmt"
)

type HaveMsg struct {
	pieceIdx uint32
}

func (*HaveMsg) Id() MsgId { return Have }

func (*HaveMsg) Len() uint { return 5 }

func (m *HaveMsg) Payload() []byte {
	return binary.BigEndian.AppendUint32(nil, m.pieceIdx)
}

func (m *HaveMsg) PieceIndex() uint32 {
	return m.pieceIdx
}

func newHaveMsg(payload []byte) (HaveMsg, error) {
	if len(payload) != 4 {
		return HaveMsg{}, fmt.Errorf("expect payload of length 4, get: %v", payload)
	}

	pieceIdx := binary.BigEndian.Uint32(payload)
	return HaveMsg{pieceIdx: pieceIdx}, nil
}
