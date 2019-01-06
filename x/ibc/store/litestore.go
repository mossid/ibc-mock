package store

/*

import (
	"errors"

	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

type LiteStore interface {
	Get([]byte) ([]byte, error)
	Has([]byte) (bool, error)
	Iterator([]byte, []byte) sdk.Iterator
	ReverseIterator([]byte, []byte) sdk.Iterator
}

// liteStore is just a map from []byte to []byte,
// the proofs are checked when it is constructed
// liteStore does **not** implement sdk.KVStore
// It cannot modify the value, and should return error when there is no proof
type liteStore struct {
	m      map[string][]byte
	ranges []interval
}

var _ LiteStore = liteStore{}

func NewLiteStore(root []byte, proofs []types.Proof, path merkle.KeyPath) (LiteStore, error) {
	kvps, err := kvpairs(root, proofs, path)
	if err != nil {
		return nil, err
	}

	res := make(liteStore)
	for _, kvp := range kvps {
		res[string(kvp.Key)] = kvp.Value
	}

	return res, nil
}

// TODO: extend proofs to rangeproofs
func kvpairs(root []byte, proofs []types.Proof, path merkle.KeyPath) (res cmn.KVPairs, err error) {
	res = make(cmn.KVPairs, 0, len(proofs))

	prt := merkle.DefaultProofRuntime()
	for _, proof := range proofs {
		keypath := path.AppendKey(proof.Key, merkle.KeyEncodingHex).String()

		if proof.Value == nil {
			err = prt.VerifyAbsence(proof.Proof, root, keypath)
		} else {
			err = prt.VerifyValue(proof.Proof, root, keypath, proof.Value)
		}

		if err != nil {
			return
		}
		res = append(res, cmn.KVPair{Key: proof.Key, Value: proof.Value})
	}

	return
}

func (store liteStore) Get(key []byte) ([]byte, error) {
	res, ok := store[string(key)]
	if !ok {
		return nil, errors.New("no proof provided for the key")
	}
	return res, nil
}

func (store liteStore) Has(key []byte) (bool, error) {
	res, ok := store[string(key)]
	if !ok {
		return false, errors.New("no proof provided for the key")
	}
	return res != nil, nil
}

func (store liteStore) Iterator(start, end []byte) sdk.Iterator {

}

func (store liteStore) ReverseIterator(start, end []byte) sdk.Iterator {

}

type interval struct {
	start, end []byte
}

type intervals []interval

func (ivs intervals) insert(iv interval)
*/
