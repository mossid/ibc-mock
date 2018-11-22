package ibc

import (
	"bytes"
	"sort"

	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
	"github.com/mossid/ibc-mock/x/ibc/connection"
)

// reimplementing tmtypes.MockPV to make it marshallable
type mockPV struct {
	PrivKey crypto.PrivKey
}

var _ tmtypes.PrivValidator = (*mockPV)(nil)

func newMockPV() *mockPV {
	return &mockPV{ed25519.GenPrivKey()}
}

func (pv *mockPV) GetAddress() tmtypes.Address {
	return pv.PrivKey.PubKey().Address()
}

func (pv *mockPV) GetPubKey() crypto.PubKey {
	return pv.PrivKey.PubKey()
}

func (pv *mockPV) SignVote(chainID string, vote *tmtypes.Vote) error {
	signBytes := vote.SignBytes(chainID)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig
	return nil
}

func (pv *mockPV) SignProposal(string, *tmtypes.Proposal) error {
	panic("not needed")
}

func (pv *mockPV) SignHeartbeat(string, *tmtypes.Heartbeat) error {
	panic("not needed")
}

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

type MockValidators []MockValidator

var _ sort.Interface = MockValidators{}

func (vals MockValidators) Len() int {
	return len(vals)
}

func (vals MockValidators) Less(i, j int) bool {
	return bytes.Compare([]byte(vals[i].GetConsAddr()), []byte(vals[j].GetConsAddr())) == -1
}

func (vals MockValidators) Swap(i, j int) {
	it := vals[j]
	vals[j] = vals[i]
	vals[i] = it
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
