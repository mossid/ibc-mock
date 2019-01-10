package store

import (
	//	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
type keypath struct {
	keypath merkle.KeyPath
	prefix  []byte
}

func (path keypath) _prefix(prefix []byte) (res keypath) {
	res.keypath = path.keypath
	res.prefix = make([]byte, len(path.prefix)+len(prefix))
	copy(res.prefix[:len(path.prefix)], path.prefix)
	copy(res.prefix[len(path.prefix):], prefix)
	return
}
*/
type Base struct {
	cdc     *codec.Codec
	storefn func(Context) KVStore
	prefix  []byte
	//	keypath keypath // temporal
}

func NewBase(cdc *codec.Codec, key sdk.StoreKey) Base {
	return Base{
		cdc:     cdc,
		storefn: func(ctx Context) KVStore { return ctx.KVStore(key) },
		/*
		   keypath: keypath{
		   			keypath: new(KeyPath).AppendKey([]byte(key.Name()), merkle.KeyEncodingHex),
		   		},
		*/
	}
}

func (base Base) store(ctx Context) KVStore {
	return NewPrefixStore(base.storefn(ctx), base.prefix)
}

func (base Base) Prefix(prefix []byte) (res Base) {
	res = Base{
		cdc:     base.cdc,
		storefn: base.storefn,
		//keypath: base.keypath._prefix(prefix),
	}
	res.prefix = join(base.prefix, prefix)
	return
}

func (base Base) key(key []byte) []byte {
	return join(base.prefix, key)
}

func join(a, b []byte) (res []byte) {
	res = make([]byte, len(a)+len(b))
	copy(res, a)
	copy(res[len(a):], b)
	return
}

/*
func (base Base) KeyPath() merkle.KeyPath {
	if len(base.keypath.prefix) != 0 {
		return base.keypath.keypath.AppendKey(base.keypath.prefix, merkle.KeyEncodingHex)
	}
	return base.keypath.keypath
}
*/
