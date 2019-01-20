package ibctest

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/mossid/ibc-mock/x/ibc2"
)

var cdc *codec.Codec

func init() {
	cdc = codec.New()
	codec.RegisterCrypto(cdc)
	ibc.RegisterCodec(cdc)
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgSpeak{}, "ibctest/MsgSpeak", nil)
	cdc.RegisterConcrete(MsgListen{}, "ibctest/MsgListen", nil)
	cdc.RegisterConcrete(MsgSpeakSafe{}, "ibctest/MsgSpeakSafe", nil)
	cdc.RegisterConcrete(MsgPushMessage{}, "ibctest/MsgPushMessage", nil)
	cdc.RegisterConcrete(PacketPushMessage{}, "ibctest/PacketPushMessage", nil)
	//cdc.RegisterConcrete(MsgListenSafe{}, "ibc/MsgListenSafe", nil)
}
