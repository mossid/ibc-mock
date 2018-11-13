package store

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Value struct {
	cdc   *codec.Codec
	store sdk.KVStore
	key   []byte
}

func NewValue(cdc *codec.Codec, store sdk.KVStore, key []byte) Value {
	return Value{
		cdc:   cdc,
		store: store,
		key:   key,
	}
}

func (v Value) Get(ptr interface{}) {
	v.cdc.MustUnmarshalBinaryBare(v.store.Get(v.key), ptr)
}

func (v Value) Set(o interface{}) {
	v.store.Set(v.key, v.cdc.MustMarshalBinaryBare(o))
}

func (v Value) Exists() bool {
	return v.store.Has(v.key)
}

//...
