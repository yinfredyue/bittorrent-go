package message

type UnchokeMsg struct{}

func (*UnchokeMsg) Id() MsgId { return Unchoke }

func (*UnchokeMsg) Len() uint { return 1 }

func (*UnchokeMsg) Payload() []byte { return []byte{} }
