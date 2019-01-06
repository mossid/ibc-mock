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

	MsgReceive = channel.MsgReceive

	PayloadConnectionListening = channel.PayloadConnectionListening
	PayloadChannelListening    = channel.PayloadChannelListening
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
	Header           = types.Header
	Packet           = types.Packet
	ConnectionConfig = types.ConnectionConfig
)

const (
	EncodingAmino = types.EncodingAmino
)
