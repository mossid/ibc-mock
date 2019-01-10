package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc *codec.Codec

func init() {
	cdc = codec.New()
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgOpen{}, "ibc/MsgOpen", nil)
	cdc.RegisterConcrete(MsgReceive{}, "ibc/MsgReceive", nil)

	cdc.RegisterConcrete(Header{}, "ibc/Header", nil)
	cdc.RegisterInterface((*Payload)(nil), nil)
	cdc.RegisterConcrete(Proof{}, "ibc/Proof", nil)
	cdc.RegisterConcrete(Packet{}, "ibc/Packet", nil)
}
