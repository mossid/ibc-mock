package store

import (
	"io"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Rootmultistore mockup, implemented as prefixstore, to mimic recursive multistore
// Move to RootMultiStore after recursive multistore has been implemented

type MultiStore struct {
	key sdk.StoreKey

	prefixes   map[sdk.StoreKey][]byte
	keysByName map[string]sdk.StoreKey
}

var _ sdk.KVStore = (*MultiStore)(nil) // TODO: change to multistore

func NewMultiStore(key sdk.StoreKey) *MultiStore {
	return &MultiStore{
		key:        key,
		prefixes:   make(map[sdk.StoreKey][]byte),
		keysByName: make(map[string]sdk.StoreKey),
	}
}

func (ms *MultiStore) MountStore(key sdk.StoreKey) {
	if key == nil {
		panic("nil key")
	}
	if _, ok := ms.prefixes[key]; ok {
		panic("dup key")
	}
	if _, ok := ms.keysByName[key.Name()]; ok {
		panic("dup name")
	}

	ms.prefixes[key] = []byte(key.Name())
	ms.keysByName[key.Name()] = key
}

// XXX: currently impossible to embed prefix information to storekey
func (ms *MultiStore) GetKVStore(key sdk.StoreKey) func(Context) KVStore {
	return func(ctx Context) KVStore {
		return NewPrefixStore(ctx.KVStore(ms.key), ms.prefixes[key])
	}
}

// Dosen't matter, *MultiStore will not be used as a KVStore anyway...
// Implementing KVStore so we can retrieve it from the context
// We don't have ctx.MultiStore() yet
func (ms *MultiStore) CacheWrap() sdk.CacheWrap    { panic("mockup") }
func (ms *MultiStore) GetStoreType() sdk.StoreType { panic("mockup") }
func (ms *MultiStore) CacheWrapWithTrace(w io.Writer, tc sdk.TraceContext) sdk.CacheWrap {
	panic("mockup")
}
func (ms *MultiStore) Get(key []byte) []byte                       { panic("mockup") }
func (ms *MultiStore) Set(key []byte, value []byte)                { panic("mockup") }
func (ms *MultiStore) Has(key []byte) bool                         { panic("mockup") }
func (ms *MultiStore) Delete(key []byte)                           { panic("mockup") }
func (ms *MultiStore) Iterator([]byte, []byte) sdk.Iterator        { panic("mockup") }
func (ms *MultiStore) ReverseIterator([]byte, []byte) sdk.Iterator { panic("mockup") }

func (ms *MultiStore) Prefix([]byte) sdk.KVStore                   { panic("mockup") }
func (ms *MultiStore) Gas(sdk.GasMeter, sdk.GasConfig) sdk.KVStore { panic("mockup") }
