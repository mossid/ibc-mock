package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PayloadChannelListening struct {
	Info ChannelInfo
}

var _ Payload = PayloadChannelListening{}

func (pay PayloadChannelListening) Route() string {
	return Route
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

func (pay PayloadChannelListening) QueueID() QueueID {
	return PacketQueueID()
}
