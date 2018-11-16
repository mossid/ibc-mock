package connection

import (
	"github.com/tendermint/tendermint/crypto"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
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

type SnapshotValidatorSet struct {
	byconsaddr store.Mapping
	bypower    store.Sorted
	totalpower store.Value

	valset ValidatorSet
}

var _ ValidatorSet = SnapshotValidatorSet{}

func (valset SnapshotValidatorSet) Snapshot(ctx sdk.Context) {
	// TODO: make it efficient
	valset.byconsaddr.Clear(ctx)
	valset.bypower.Clear(ctx)

	valset.valset.IterateBondedValidatorsByPower(ctx, func(_ int64, v Validator) (stop bool) {
		valset.byconsaddr.Set(ctx, v.GetConsAddr(), v.GetOperator())
		valset.bypower.Set(ctx, uint64(v.GetPower().RoundInt64()), v.GetOperator(), v)
		return
	})
}

func (valset SnapshotValidatorSet) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(int64, Validator) bool) {
	var v Validator
	i := int64(0)
	valset.bypower.IterateDescending(ctx, &v, func(_ []byte) (stop bool) {
		stop = fn(i, v)
		i++
		return
	})
}

func (valset SnapshotValidatorSet) ValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) (res Validator) {
	valset.byconsaddr.Get(ctx, addr, &res)
	return
}

func (valset SnapshotValidatorSet) TotalPower(ctx sdk.Context) (res sdk.Dec) {
	valset.totalpower.Get(ctx, &res)
	return
}

func (valset SnapshotValidatorSet) IsStillValid(ctx sdk.Context, diffreq sdk.Dec) (res bool) {
	required := valset.TotalPower(ctx).Mul(diffreq)
	available := sdk.ZeroDec()

	valset.valset.IterateBondedValidatorsByPower(ctx, func(_ int64, v Validator) (stop bool) {
		available = available.Add(valset.ValidatorByConsAddr(ctx, v.GetConsAddr()).GetPower())
		if available.GTE(required) {
			stop = true
			res = true
		}
		return
	})

	return
}
