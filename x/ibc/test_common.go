package ibc

import (
	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
	"github.com/mossid/ibc-mock/x/ibc/connection"
)

// MockValset
type MockValidator struct {
	MockPV *mockPV
	Power  sdk.Dec
}

func NewMockValidator() MockValidator {
	return MockValidator{
		MockPV: newMockPV(),
		Power:  sdk.NewDec(0),
	}
}

var _ connection.Validator = MockValidator{}

func (val MockValidator) GetOperator() sdk.ValAddress {
	return sdk.ValAddress(val.MockPV.GetAddress())
}

func (val MockValidator) GetConsAddr() sdk.ConsAddress {
	return sdk.GetConsAddress(val.MockPV.GetPubKey())
}

func (val MockValidator) GetConsPubKey() crypto.PubKey {
	return val.MockPV.GetPubKey()
}

func (val MockValidator) GetPower() sdk.Dec {
	return val.Power
}

type AddressPowerPair struct {
	Address sdk.ConsAddress
	Power   sdk.Dec
}

type MockValidatorSet struct {
	power      store.Value
	byPower    store.Sorted
	validators store.Mapping
}

func NewMockValidatorSet(cdc *codec.Codec, key sdk.StoreKey) MockValidatorSet {
	base := store.NewBase(cdc, key)

	return MockValidatorSet{
		power:      store.NewValue(base, []byte{0x00}),
		byPower:    store.NewSorted(base, []byte{0x01}, store.BinIndexerEnc),
		validators: store.NewMapping(base, []byte{0x02}),
	}
}

func valsetInitGenesis(ctx sdk.Context, valset MockValidatorSet, vals []MockValidator) {
	Power := sdk.NewDec(0)

	for _, val := range vals {
		Power = Power.Add(val.GetPower())
		valset.byPower.Set(ctx, uint64(val.GetPower().RoundInt64()), val.GetConsAddr(), val.GetConsAddr())
		valset.validators.Set(ctx, val.GetConsAddr(), val)
	}

	valset.power.Set(ctx, Power)
}

var _ connection.ValidatorSet = MockValidatorSet{}

func (valset MockValidatorSet) IterateBondedValidatorsByPower(ctx sdk.Context, fn func(int64, connection.Validator) bool) {
	var addr sdk.ConsAddress
	var index int64

	valset.byPower.IterateDescending(ctx, &addr, func(key []byte) (stop bool) {
		var val MockValidator
		valset.validators.Get(ctx, addr, &val)
		stop = fn(index, val)
		index++
		return
	})
}

func (valset MockValidatorSet) TotalPower(ctx sdk.Context) (res sdk.Dec) {
	valset.power.Get(ctx, &res)
	return
}

func (valset MockValidatorSet) ValidatorByConsAddr(ctx sdk.Context, addr sdk.ConsAddress) (res connection.Validator) {
	valset.validators.Get(ctx, addr, &res)
	return
}
