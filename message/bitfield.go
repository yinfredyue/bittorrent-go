package message

import "encoding/binary"

type BitfieldMsg struct {
	bits []byte
}

func (*BitfieldMsg) Id() MsgId {
	return Bitfield
}

func (m *BitfieldMsg) Len() uint {
	return 1 + uint(len(m.bits))
}

func (m *BitfieldMsg) Payload() []byte {
	return m.bits
}

func bytesToUint64s(b []byte) []uint64 {
	numBytesToPad := 8 - (len(b) % 8)
	padding := make([]byte, numBytesToPad)
	b = append(b, padding...)

	numUint64s := len(b) / 8
	res := make([]uint64, numUint64s)

	for i := 0; i < numUint64s; i++ {
		res[i] = binary.BigEndian.Uint64(b[8*i : (8 * (i + 1))])
	}

	return res
}

func (m *BitfieldMsg) Bits() []uint64 {
	return bytesToUint64s(m.bits)
}

func newBitfieldMsg(bits []byte) BitfieldMsg {
	return BitfieldMsg{bits: bits}
}
