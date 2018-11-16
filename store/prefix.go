package store

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewPrefixStore(parent sdk.KVStore, prefix []byte) sdk.KVStore {
	return parent.Prefix(prefix)
}

type ProofOpWrapperPrefix struct {
	merkle.ProofOperator

	prefix []byte
}

var _ merkle.ProofOperator = ProofOpWrapperPrefix{}

func (op ProofOpWrapperPrefix) GetKey() []byte {
	return append(op.prefix, op.ProofOperator.GetKey()...)
}

func (op ProofOpWrapperPrefix) ProofOp() merkle.ProofOp {
	panic("not yet implemented")
}
