package connection

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc = codec.New()

func init() {
	codec.RegisterCrypto(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCheckpoint{}, "ibc/MsgCheckpoint", nil)
	cdc.RegisterConcrete(MsgOpenConnection{}, "ibc/MsgOpenConnection", nil)

	cdc.RegisterInterface((*Validator)(nil), nil)
}
