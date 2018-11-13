package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewPrefixStore(parent sdk.KVStore, prefix []byte) sdk.KVStore {
	return parent.Prefix(prefix)
}
