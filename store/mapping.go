package store

import (
	"reflect"

	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/store"
)

type Mapping struct {
	base       Base
	prefix     []byte
	start, end []byte
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

func (m Mapping) store(ctx Context) KVStore {
	if len(m.prefix) != 0 {
		return NewPrefixStore(m.base.store(ctx), m.prefix)
	}
	return m.base.store(ctx)
}

func (m Mapping) Value(key []byte) Value {
	return NewValue(
		NewBaseWithAccessor(
			m.base.cdc,
			m.store,
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

func (m Mapping) GetSafe(ctx Context, key []byte, ptr interface{}) *GetSafeError {
	return m.Value(key).GetSafe(ctx, ptr)
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

func (m Mapping) IsEmpty(ctx Context) bool {
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

func (m Mapping) Range(start, end []byte) Mapping {
	return Mapping{
		base:  m.base,
		start: start,
		end:   end,
	}
}

// go-amino does not support decoding to a non-nil interface
func setnil(ptr interface{}) {
	v := reflect.ValueOf(ptr)
	v.Elem().Set(reflect.Zero(v.Elem().Type()))
}

func (m Mapping) Iterate(ctx Context, ptr interface{}, fn func([]byte) bool) {
	iter := m.store(ctx).Iterator(m.start, m.end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		setnil(ptr)

		v := iter.Value()

		m.base.cdc.MustUnmarshalBinaryBare(v, ptr)

		if fn(iter.Key()) {
			break
		}
	}
}

func (m Mapping) IterateKeys(ctx Context, fn func([]byte) bool) {
	iter := m.store(ctx).Iterator(m.start, m.end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if fn(iter.Key()) {
			break
		}
	}
}

func (m Mapping) ReverseIterate(ctx Context, ptr interface{}, fn func([]byte) bool) {
	iter := m.store(ctx).ReverseIterator(m.start, m.end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		setnil(ptr)

		v := iter.Value()

		m.base.cdc.MustUnmarshalBinaryBare(v, ptr)

		if fn(iter.Key()) {
			break
		}
	}
}

func (m Mapping) ReverseIterateKeys(ctx Context, fn func([]byte) bool) {
	iter := m.store(ctx).ReverseIterator(m.start, m.end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		if fn(iter.Key()) {
			break
		}
	}
}
func (m Mapping) First(ctx Context, ptr interface{}) (key []byte, ok bool) {
	kvp, ok := store.First(m.base.store(ctx), m.start, m.end)
	if !ok {
		return
	}
	key = kvp.Key
	if ptr != nil {
		m.base.cdc.MustUnmarshalBinaryBare(kvp.Value, ptr)
	}
	return
}

func (m Mapping) Last(ctx Context, ptr interface{}) (key []byte, ok bool) {

	kvp, ok := last(m.base.store(ctx), m.start, m.end)
	if !ok {
		return
	}
	key = kvp.Key
	if ptr != nil {
		m.base.cdc.MustUnmarshalBinaryBare(kvp.Value, ptr)
	}
	return
}

func (m Mapping) Clear(ctx Context) {
	var keys [][]byte

	iter := m.base.store(ctx).ReverseIterator(m.start, m.end)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		keys = append(keys, iter.Key())
	}

	store := m.base.store(ctx)
	for _, key := range keys {
		store.Delete(key)
	}
}
