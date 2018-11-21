package ibc

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
)

func registerCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
	cdc.RegisterInterface((*tmtypes.PrivValidator)(nil), nil)
	cdc.RegisterConcrete(mockPV{}, "test/mockPV", nil)
	cdc.RegisterConcrete(MockValidator{}, "test/MockValidator", nil)
}

func getMockApp(t *testing.T) (*mock.App, Keeper, MockValidatorSet, []MockValidator) {
	mApp := mock.NewApp()

	registerCodec(mApp.Cdc)

	key := sdk.NewKVStoreKey("ibc")
	valsetkey := sdk.NewKVStoreKey("valset")
	vals := getValidators()
	valset := NewMockValidatorSet(mApp.Cdc, valsetkey)
	keeper := NewKeeper(mApp.Cdc, key, valset)

	mApp.Router().AddRoute("ibc", NewHandler(keeper))
	mApp.SetInitChainer(getInitChainer(mApp, keeper, valset, vals))
	mApp.SetAnteHandler(nil) // overriding antehandler to bypass auth logic
	mApp.SetBeginBlocker(nil)
	mApp.SetEndBlocker(nil)

	require.NoError(t, mApp.CompleteSetup(key, valsetkey))

	return mApp, keeper, valset, vals
}

func getInitChainer(mapp *mock.App, k Keeper, valset MockValidatorSet, vals []MockValidator) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
		fmt.Println(vals)
		valsetInitGenesis(ctx, valset, vals)

		res.Validators = make([]abci.ValidatorUpdate, len(vals))
		for i, val := range vals {
			res.Validators[i] = abci.ValidatorUpdate{
				PubKey: tmtypes.TM2PB.PubKey(val.GetConsPubKey()),
				Power:  val.GetPower().RoundInt64(),
			}
		}
		return
	}
}

const valsetnum = 10

func getValidators() []MockValidator {
	vals := make([]MockValidator, valsetnum)
	for i := range vals {
		vals[i] = MockValidator{newMockPV(), sdk.NewDec(int64(i))}
	}

	return vals
}

type node struct {
	*mock.App
	k      Keeper
	valset MockValidatorSet

	vals    []MockValidator
	commits map[int64]lite.FullCommit
	height  int64

	seq        int64
	lastheader abci.Header
}

func getNode(t *testing.T) node {
	app, k, valset, vals := getMockApp(t)

	return node{
		App:     app,
		k:       k,
		valset:  valset,
		vals:    vals,
		commits: make(map[int64]lite.FullCommit),
	}
}

func (node *node) simulateSigning() {
	tmheader := tmtypes.Header{
		Height:  node.lastheader.Height,
		AppHash: node.lastheader.AppHash,
	}

	precommits := node.sign(tmheader)

	shdr := tmtypes.SignedHeader{
		Header: &tmheader,
		Commit: &tmtypes.Commit{
			Precommits: precommits,
		},
	}

	fc := lite.NewFullCommit(shdr, node.tmvalset(), node.tmvalset()) // TODO: implement dynamic
	node.commits[node.lastheader.Height] = fc
}

func (node *node) signCheckDeliver(t *testing.T, msg sdk.Msg, expSimPass, expPass bool) sdk.Result {
	tx := mock.GenTx([]sdk.Msg{msg}, []int64{0}, []int64{node.seq}, node.vals[0].MockPV.PrivKey)
	res := node.Simulate(tx)

	if expSimPass {
		require.Equal(t, sdk.ABCICodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.ABCICodeOK, res.Code, res.Log)
	}

	node.height++

	node.BeginBlock(abci.RequestBeginBlock{Header: node.lastheader})
	res = node.Deliver(tx)

	if expPass {
		require.Equal(t, sdk.ABCICodeOK, res.Code, res.Log)
	} else {
		require.NotEqual(t, sdk.ABCICodeOK, res.Code, res.Log)
	}

	node.EndBlock(abci.RequestEndBlock{Height: node.height})

	cres := node.Commit()

	node.lastheader = abci.Header{
		Height:  node.height,
		AppHash: cres.Data,
	}

	node.simulateSigning()

	return res
}

func (node *node) tmvalset() *tmtypes.ValidatorSet {
	vals := make([]*tmtypes.Validator, len(node.vals))
	for i, val := range node.vals {
		vals[i] = &tmtypes.Validator{
			Address:     val.MockPV.GetAddress(),
			PubKey:      val.MockPV.GetPubKey(),
			VotingPower: val.Power.RoundInt64(),
		}
	}

	return &tmtypes.ValidatorSet{
		Validators: vals,
	}
}

func (node *node) sign(header tmtypes.Header) []*tmtypes.Vote {
	precommits := make([]*tmtypes.Vote, len(node.vals))
	for i, val := range node.vals {
		precommits[i] = &tmtypes.Vote{
			BlockID: tmtypes.BlockID{
				Hash: header.Hash(),
			},
			ValidatorAddress: val.MockPV.GetAddress(),
			ValidatorIndex:   i,
			Height:           node.height,
			Type:             tmtypes.PrecommitType,
		}
		val.MockPV.SignVote("", precommits[i])
	}
	return precommits
}

func (node *node) startChain(t *testing.T) {
	node.InitChain(abci.RequestInitChain{})
}

func (node *node) openConnection(t *testing.T, target *node) {
	// send MsgOpenConnection
}

func (node *node) establishConnection(t *testing.T) {
	// relay PayloadConnectionListening
}

func (node *node) openChannel(t *testing.T) {
	// call MyChannel.Open()
}

func (node *node) establishChannel(t *testing.T) {
	// relay PayloadChannelListening
}

func (node *node) sendPackets(t *testing.T) {
	// send []MsgReceive, relay []MyPayload
}

func (node *node) updateConnection(t *testing.T) {
	// send MsgUpdateConnection
}

func TestNode(t *testing.T) {
	node := getNode(t)
	node.startChain(t)

	for i := 0; i < 10; i++ {
		node.signCheckDeliver(t, MsgCheckpoint{}, true, true)
	}
}

func TestSimulateSigning(t *testing.T) {
	node := getNode(t)
	node.startChain(t)

	verifier := lite.NewBaseVerifier("", 0, node.tmvalset())

	for i := 0; i < 10; i++ {
		node.signCheckDeliver(t, MsgCheckpoint{}, true, true)
		err := verifier.Verify(node.commits[int64(i)+1].SignedHeader)
		require.NoError(t, err)
	}
}

/*
func TestLifecycle(t *testing.T) {
	node := getNode(t)
	node.openConnection(t)
	node.establishConnection(t)
	node.openChannel(t)
	node.establishChannel(t)
	node.sendPackets(t)
	node.updateConnection(t)
	node.sendPackets(t)
}

func TestLongDuration(t *testing.T) {
	node := getNode(t)
	node.openConnection(t)
	node.establishConnection(t)
	node.openChannel(t)
	node.establishChannel(t)

	for i := 0; i < 1000; i++ {
		node.sendPackets(t)
		node.updateConnection(t)
	}
}*/
