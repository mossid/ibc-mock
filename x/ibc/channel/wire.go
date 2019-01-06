package channel

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

var cdc = codec.New()

func init() {
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgReceive{}, "ibc/MsgReceive", nil)
	cdc.RegisterConcrete(PayloadConnectionListening{}, "ibc/PayloadConnectionListening", nil)
	cdc.RegisterConcrete(PayloadChannelListening{}, "ibc/PayloadChannelListening", nil)

}
