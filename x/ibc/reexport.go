package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/connection"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

type (
	Keeper = connection.Keeper

	MsgOpenConnection = connection.MsgOpenConnection
	MsgCheckpoint     = connection.MsgCheckpoint
)

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, valset connection.ValidatorSet) (k Keeper) {
	return connection.NewKeeper(cdc, key, valset)
}

func NewHandler(k Keeper) sdk.Handler {
	return connection.NewHandler(k)
}

type (
	ChannelCore = channel.ChannelCore
)

type (
	ConnectionConfig = types.ConnectionConfig
)

const (
	EncodingAmino = types.EncodingAmino
)
