package message

type KeepAliveMsg struct{}

func (*KeepAliveMsg) Id() MsgId { return KeepAlive }

func (*KeepAliveMsg) Len() uint { return 0 }

func (*KeepAliveMsg) Payload() []byte { return []byte{} }
