package store

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type List = store.List

func NewList(cdc *codec.Codec, parent sdk.KVStore) store.List {
	return store.NewList(cdc, parent)
}

type Queue = store.Queue

func NewQueue(cdc *codec.Codec, parent sdk.KVStore) store.Queue {
	return store.NewQueue(cdc, parent)
}
