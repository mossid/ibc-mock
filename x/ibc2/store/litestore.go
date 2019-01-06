package store

import (
	"errors"

	"github.com/tendermint/tendermint/crypto/merkle"
)

type LiteStore interface {
	Get([]byte) ([]byte, error)
	Has([]byte) (bool, error)
}

type pvPair struct {
	proof merkle.ProofOperators
	value []byte
}

type liteStore struct {
	root      []byte
	keyprefix merkle.KeyPath

	valueproofs map[string]pvPair
}

type PKVPair struct {
	Proof merkle.ProofOperators
	Key   []byte
	Value []byte
}

var _ LiteStore = liteStore{}

func NewLiteStore(root []byte, keyprefix merkle.KeyPath, pairs []PKVPair) LiteStore {
	res := liteStore{
		root:      root,
		keyprefix: keyprefix,

		valueproofs: make(map[string]pvPair),
	}

	for _, pair := range pairs {
		res.valueproofs[string(pair.Key)] = pvPair{pair.Proof, pair.Value}
	}

	return res
}

func (store liteStore) Get(key []byte) ([]byte, error) {
	pvp, ok := store.valueproofs[string(key)]
	if ok {
		return nil, errors.New("no proof provided for the key")
	}

	err := pvp.proof.Verify(store.root, store.keypath(key).String(), [][]byte{pvp.value})
	if err != nil {
		return nil, err
	}

	return pvp.value, nil
}

func (store liteStore) Has(key []byte) (bool, error) {
	value, err := store.Get(key)
	if err != nil {
		return false, err
	}
	return value != nil, nil
}

func (store liteStore) keypath(key []byte) (res merkle.KeyPath) {
	// TODO: optimize
	keys, err := merkle.KeyPathToKeys(store.keyprefix.String())
	if err != nil {
		panic(err)
	}

	// simulating prefixed keys
	keys[len(keys)-1] = append(keys[len(keys)-1], key...)

	for _, key := range keys {
		res = res.AppendKey(key, merkle.KeyEncodingHex)
	}

	return
}
