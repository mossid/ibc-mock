package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var msgCdc *codec.Codec

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*User)(nil), nil)
	cdc.RegisterInterface((*Packet)(nil), nil)
	cdc.RegisterConcrete(ChannelUser{}, "ibc/ChannelUser", nil)
	cdc.RegisterConcrete(ConnConfig{}, "ibc/ConnConfig", nil)
	cdc.RegisterConcrete(PortConfig{}, "ibc/PortConfig", nil)
}
