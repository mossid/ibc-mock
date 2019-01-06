package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/connection"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

func RegisterCodec(cdc *codec.Codec) {
	types.RegisterCodec(cdc)
	channel.RegisterCodec(cdc)
	connection.RegisterCodec(cdc)
}
