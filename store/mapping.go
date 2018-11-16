package store

import (
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/store"
)

type Mapping struct {
	base   Base
	prefix []byte
}

func NewMapping(base Base, prefix []byte) Mapping {
	return Mapping{
		base:   base,
		prefix: prefix,
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
	return NewValue(
		NewBaseWithAccessor(
			m.base.Codec,
			func(ctx Context) KVStore {
				return NewPrefixStore(m.base.store(ctx), m.prefix)
			},
		),
		key,
	)
}

func (m Mapping) Get(ctx Context, key []byte, ptr interface{}) {
	m.Value(key).Get(ctx, ptr)
}

func (m Mapping) GetIfExists(ctx Context, key []byte, ptr interface{}) {
	m.Value(key).GetIfExists(ctx, ptr)
}

func (m Mapping) Set(ctx Context, key []byte, o interface{}) {
	m.Value(key).Set(ctx, o)
}

func (m Mapping) Has(ctx Context, key []byte) bool {
	return m.Value(key).Exists(ctx)
}

func (m Mapping) Delete(ctx Context, key []byte) {
	m.Value(key).Delete(ctx)
}

func (m Mapping) IsEmpty(ctx Context) (ok bool) {
	iter := m.base.store(ctx).Iterator(nil, nil)
	defer iter.Close()
	return iter.Valid()
}

func (m Mapping) Prefix(prefix []byte) Mapping {
	return NewMapping(
		m.base,
		append(m.prefix, prefix...),
	)
}

func (m Mapping) Iterate(ctx Context, ptr interface{}, fn func([]byte) bool) {
	iter := m.base.store(ctx).Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		v := iter.Value()
		m.base.MustUnmarshalBinaryBare(v, ptr)

		if fn(iter.Key()) {
			break
		}
	}
}

func (m Mapping) First(ctx Context, ptr interface{}) (key []byte, ok bool) {
	kvp, ok := store.First(m.base.store(ctx), nil, nil)
	if !ok {
		return
	}
	key = kvp.Key
	if ptr != nil {
		m.base.MustUnmarshalBinaryBare(kvp.Value, ptr)
	}
	return
}

func (m Mapping) Last(ctx Context, ptr interface{}) (key []byte, ok bool) {
	kvp, ok := store.Last(m.base.store(ctx), nil, nil)
	if !ok {
		return
	}
	key = kvp.Key
	if ptr != nil {
		m.base.MustUnmarshalBinaryBare(kvp.Value, ptr)
	}
	return
}
