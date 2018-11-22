package connection

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Validator interface {
	GetOperator() sdk.ValAddress
	GetConsAddr() sdk.ConsAddress
	GetConsPubKey() crypto.PubKey
	GetPower() sdk.Dec
}

type ValidatorSet interface {
	IterateBondedValidatorsByPower(sdk.Context,
		func(index int64, validator Validator) (stop bool))

	TotalPower(sdk.Context) sdk.Dec
	ValidatorByConsAddr(sdk.Context, sdk.ConsAddress) Validator
}

type StakeValidatorSet struct {
	valset sdk.ValidatorSet
}

var _ ValidatorSet = StakeValidatorSet{}

func (valset StakeValidatorSet) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(int64, Validator) bool) {
	valset.valset.IterateBondedValidatorsByPower(ctx, func(index int64, v sdk.Validator) bool {
		return fn(index, v)
	})
}

func (valset StakeValidatorSet) TotalPower(ctx sdk.Context) sdk.Dec {
	return valset.valset.TotalPower(ctx)
}

func (valset StakeValidatorSet) ValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) Validator {
	return valset.valset.ValidatorByConsAddr(ctx, addr)
}

func FromStakeValidatorSet(valset sdk.ValidatorSet) ValidatorSet {
	return StakeValidatorSet{valset}
}
