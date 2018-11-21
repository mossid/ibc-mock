package store

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Base struct {
	cdc   *codec.Codec
	store func(Context) KVStore
}

func NewBase(cdc *codec.Codec, key sdk.StoreKey) Base {
	return Base{
		cdc:   cdc,
		store: func(ctx Context) KVStore { return ctx.KVStore(key) },
	}
}

func NewBaseWithAccessor(cdc *codec.Codec, store func(Context) KVStore) Base {
	return Base{
		cdc:   cdc,
		store: store,
	}
}
