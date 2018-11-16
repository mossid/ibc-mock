package ibc

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/connection"
)

type AddressPowerPair struct {
	Address sdk.ConsAddress
	Power   sdk.Dec
}

type MockValidatorSet struct {
	Power      sdk.Dec
	ByPower    []AddressPowerPair
	Validators map[string]connection.Validator
}

var _ connection.ValidatorSet = MockValidatorSet{}

func NewMockValidatorSet(vals ...connection.Validator) (res MockValidatorSet) {
	res = MockValidatorSet{
		Power:      sdk.NewDec(0),
		ByPower:    make([]AddressPowerPair, len(vals)),
		Validators: make(map[string]connection.Validator),
	}

	for i, val := range vals {
		res.Power = res.Power.Add(val.GetPower())
		res.ByPower[i] = AddressPowerPair{val.GetConsAddr(), val.GetPower()}
		res.Validators[string(val.GetConsAddr())] = val
	}

	sort.Slice(res.ByPower, func(i, j int) bool { return res.ByPower[i].Power.LT(res.ByPower[j].Power) })

	return
}

func (valset MockValidatorSet) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(int64, connection.Validator) bool) {
	for i, pair := range valset.ByPower {
		if fn(int64(i), valset.Validators[string(pair.Address)]) {
			break
		}
	}
}

func (valset MockValidatorSet) TotalPower(ctx sdk.Context) sdk.Dec {
	return valset.Power
}

func (valset MockValidatorSet) ValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) connection.Validator {
	return valset.Validators[string(addr)]
}
