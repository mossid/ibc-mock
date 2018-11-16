package channel

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc = codec.New()

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgReceive{}, "ibc/MsgReceive", nil)
	cdc.RegisterConcrete(PayloadChannelListening{}, "ibc/PayloadChannelListening", nil)

}
