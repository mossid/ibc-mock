package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(Header{}, "ibc/Header", nil)
	cdc.RegisterConcrete(Packet{}, "ibc/Packet", nil)
	cdc.RegisterInterface((*Payload)(nil), nil)
	cdc.RegisterConcrete(ChannelInfo{}, "ibc/ChannelInfo", nil)
}
