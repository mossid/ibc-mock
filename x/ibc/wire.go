package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/connection"
)

func RegisterCodec(cdc *codec.Codec) {
	channel.RegisterCodec(cdc)
	connection.RegisterCodec(cdc)
}
