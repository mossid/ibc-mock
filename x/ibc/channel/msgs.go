package channel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

type MsgReceive struct {
	Packets []types.Packet
	Proof   types.Proof
	Relayer sdk.AccAddress
}

var _ sdk.Msg = MsgReceive{}

func (msg MsgReceive) Route() string {
	return msg.Packets[0].Route()
}

func (msg MsgReceive) Type() string {
	return msg.Packets[0].Type()
}

func (msg MsgReceive) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgReceive) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}

func (msg MsgReceive) ValidateBasic() sdk.Error {
	if len(msg.Packets) == 0 {
		return ErrEmptyMsg(DefaultCodespace)
	}

	p0 := msg.Packets[0]
	route0, ty0, qid0 := p0.Route(), p0.Type(), p0.QueueID()

	for _, p := range msg.Packets[1:] {
		if p.Route() != route0 {
			return ErrUnmatchingPacket(DefaultCodespace, "invalid route")
		}
		if p.Type() != ty0 {
			return ErrUnmatchingPacket(DefaultCodespace, "invalid type")
		}
		if p.QueueID() != qid0 {
			return ErrUnmatchingPacket(DefaultCodespace, "invalid queue id")
		}
	}

	return nil
}

type PayloadChannelListening struct {
	Info types.ChannelInfo
}

var _ types.Payload = PayloadChannelListening{}

func (pay PayloadChannelListening) Route() string {
	return "ibc"
}

func (pay PayloadChannelListening) Type() string {
	return "payload_channel_listening"
}

func (pay PayloadChannelListening) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (pay PayloadChannelListening) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(pay) // TODO
}

func (pay PayloadChannelListening) QueueID() types.QueueID {
	return types.PacketQueueID()
}

type PayloadConnectionListening struct {
	Config  types.ConnectionConfig
	ChainID []byte
}

var _ types.Payload = PayloadConnectionListening{}

func (pay PayloadConnectionListening) Route() string {
	return "ibc"
}

func (pay PayloadConnectionListening) Type() string {
	return "payload_connection_listening"
}

func (pay PayloadConnectionListening) ValidateBasic() sdk.Error {
	return nil
}

func (pay PayloadConnectionListening) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(pay) //TODO
}

func (pay PayloadConnectionListening) QueueID() types.QueueID {
	return types.PacketQueueID()
}
