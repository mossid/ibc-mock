package connection

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgOpenConnection{}, "ibc/MsgOpenConnection", nil)
}
