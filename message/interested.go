package message

type InterestedMsg struct{}

func (*InterestedMsg) Id() MsgId { return Interested }

func (*InterestedMsg) Len() uint { return 1 }

func (*InterestedMsg) Payload() []byte { return []byte{} }
