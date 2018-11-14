package store

import (
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Mapping struct {
	cdc   *codec.Codec
	store sdk.KVStore
}

func NewMapping(cdc *codec.Codec, store sdk.KVStore, prefix []byte) Mapping {
	return Mapping{
		cdc:   cdc,
		store: NewPrefixStore(store, prefix),
	}
}

func NewRawPrefix(space byte, prefixes ...[]byte) (res []byte) {
	res = append(res, space)
	for _, p := range prefixes {
		res = append(res, p...)
	}
	return
}

func NewPrefix(space byte, prefixes ...[]byte) (res []byte) {
	res = append(res, space)
	for _, p := range prefixes {
		res = append(res, amino.MustMarshalBinaryLengthPrefixed(p)...)
	}
	return
}

func (m Mapping) Value(key []byte) Value {
	return NewValue(m.cdc, m.store, key)
}

func (m Mapping) Get(key []byte, ptr interface{}) {
	m.Value(key).Get(ptr)
}

func (m Mapping) GetIfExists(key []byte, ptr interface{}) {
	m.Value(key).GetIfExists(ptr)
}

func (m Mapping) Set(key []byte, o interface{}) {
	m.Value(key).Set(o)
}

func (m Mapping) Has(key []byte) bool {
	return m.Value(key).Exists()
}

func (m Mapping) Delete(key []byte) {
	m.Value(key).Delete()
}

func (m Mapping) IsEmpty() (ok bool) {
	iter := m.store.Iterator(nil, nil)
	defer iter.Close()
	return iter.Valid()
}

func (m Mapping) Iterate(ptr interface{}, fn func([]byte) bool) {
	iter := m.store.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		v := iter.Value()
		m.cdc.MustUnmarshalBinaryBare(v, ptr)

		if fn(iter.Key()) {
			break
		}
	}
}

func (m Mapping) First(ptr interface{}) (key []byte, ok bool) {
	kvp, ok := store.First(m.store, nil, nil)
	if !ok {
		return
	}
	key = kvp.Key
	if ptr != nil {
		m.cdc.MustUnmarshalBinaryBare(kvp.Value, ptr)
	}
	return
}

func (m Mapping) Last(ptr interface{}) (key []byte, ok bool) {
	kvp, ok := store.Last(m.store, nil, nil)
	if !ok {
		return
	}
	key = kvp.Key
	if ptr != nil {
		m.cdc.MustUnmarshalBinaryBare(kvp.Value, ptr)
	}
	return
}
