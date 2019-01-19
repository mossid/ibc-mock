package ibctest

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mock"

	"github.com/mossid/ibc-mock/store"
	"github.com/mossid/ibc-mock/x/ibc2"
)

var customid = ibc.ChainID{0x01, 0x23, 0x34, 0x56, 0x78, 0x89, 0xAB, 0xCD}
var portid = ibc.PortID{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}

func registerCodec(cdc *codec.Codec) {
	ibc.RegisterCodec(cdc)
	RegisterCodec(cdc)
	cdc.RegisterInterface((*tmtypes.PrivValidator)(nil), nil)
	cdc.RegisterConcrete(mockPV{}, "test/mockPV", nil)
	cdc.RegisterConcrete(MockValidator{}, "test/MockValidator", nil)
}

func getMockApp(t *testing.T) (*mock.App, keeper, MockValidatorSet, []MockValidator) {
	mApp := mock.NewApp()

	registerCodec(mApp.Cdc)

	key := sdk.NewKVStoreKey("ibc")
	valsetkey := sdk.NewKVStoreKey("valset")
	vals := getValidators()
	valset := NewMockValidatorSet(mApp.Cdc, valsetkey)
	ibck := ibc.NewKeeper(mApp.Cdc, key)
	keeper := newKeeper(ibck.Channel(Route))

	mApp.Router().AddRoute("ibctest", NewHandler(keeper))
	mApp.SetInitChainer(getInitChainer(ibck, valset, vals))
	// overriding antehandler to bypass auth logic
	mApp.SetAnteHandler(ibc.NewAnteHandler(ibck))
	/*
		mApp.SetBeginBlocker(func(ctx sdk.Context, req abci.RequestBeginBlock) (res abci.ResponseBeginBlock) {
			ibc.BeginBlock(ctx, keeper)
			return
		})
	*/
	mApp.SetEndBlocker(nil)

	require.NoError(t, mApp.CompleteSetup(key, valsetkey))

	return mApp, keeper, valset, vals
}

func getInitChainer(k ibc.Keeper, valset MockValidatorSet, vals []MockValidator) sdk.InitChainer {
	return func(ctx sdk.Context, req abci.RequestInitChain) (res abci.ResponseInitChain) {
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
		vals[i] = MockValidator{newMockPV(), sdk.NewDec(int64(i + 1))}
	}

	sort.Sort(MockValidators(vals))

	return vals
}

type node struct {
	*mock.App
	k      keeper
	valset MockValidatorSet

	vals    []MockValidator
	commits map[int64]lite.FullCommit
	height  int64

	seq        uint64
	lastheader abci.Header

	connconfig ibc.ConnConfig

	lastcommit int64
}

func getNode(t *testing.T) *node {
	app, k, valset, vals := getMockApp(t)

	return &node{
		App:     app,
		k:       k,
		valset:  valset,
		vals:    vals,
		commits: make(map[int64]lite.FullCommit),
	}
}

func (node *node) simulateSigning() {
	tmheader := tmtypes.Header{
		Height:             node.lastheader.Height,
		AppHash:            node.lastheader.AppHash,
		ValidatorsHash:     node.tmvalset().Hash(),
		NextValidatorsHash: node.tmvalset().Hash(), // TODO: implement dynamic
	}

	precommits := node.sign(tmheader)

	shdr := tmtypes.SignedHeader{
		Header: &tmheader,
		Commit: &tmtypes.Commit{
			BlockID: tmtypes.BlockID{
				Hash: tmheader.Hash(),
			},
			Precommits: precommits,
		},
	}

	fc := lite.NewFullCommit(shdr, node.tmvalset(), node.tmvalset()) // TODO: implement dynamic
	node.commits[node.lastheader.Height] = fc
}

func (node *node) execblock(t *testing.T, fn func()) {
	node.height++

	node.BeginBlock(abci.RequestBeginBlock{Header: node.lastheader})

	if fn != nil {
		fn()
	}

	node.EndBlock(abci.RequestEndBlock{Height: node.height})

	cres := node.Commit()

	node.lastheader = abci.Header{
		Height:             node.height,
		AppHash:            cres.Data,
		ValidatorsHash:     node.tmvalset().Hash(),
		NextValidatorsHash: node.tmvalset().Hash(), // TODO: implement dynamic
	}

	node.simulateSigning()
}

func (node *node) signCheckDeliver(t *testing.T, msg sdk.Msg, expPass bool) (res sdk.Result) {
	node.execblock(t, func() {
		tx := mock.GenTx([]sdk.Msg{msg}, []uint64{0}, []uint64{node.seq}, node.vals[0].MockPV.PrivKey)
		res = node.Deliver(tx)

		if expPass {
			require.Equal(t, sdk.CodeOK, res.Code, res.Log)
		} else {
			require.NotEqual(t, sdk.CodeOK, res.Code, res.Log)
		}
	})

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

	return tmtypes.NewValidatorSet(vals)
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

func (node *node) speakConn(t *testing.T, target *node) sdk.Result {
	msg := MsgSpeak{
		ChainID: customid,
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())

	return res
}

func (node *node) listenConn(t *testing.T, target *node) sdk.Result {
	node.connconfig = ibc.ConnConfig{
		ROT: target.commits[1],
		RootKeyPath: new(merkle.KeyPath).
			AppendKey([]byte("ibc"), merkle.KeyEncodingHex).
			AppendKey([]byte{0x00}, merkle.KeyEncodingHex).String(),
		ChainID: customid,
	}

	msg := MsgListen{
		ChainID: customid,
		Config:  node.connconfig,
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())

	return res
}

func (node *node) update(t *testing.T, target *node) sdk.Result {
	msg := MsgUpdate{
		ChainID: customid,
		Commits: []lite.FullCommit{target.commits[target.lastheader.Height]},
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())

	node.lastcommit = target.lastheader.Height

	return res
}

func (node *node) speakSafeConn(t *testing.T, target *node) sdk.Result {
	qres := target.Query(abci.RequestQuery{
		Path:   "/store/ibc/key",
		Data:   append([]byte{0x00}, customid[:]...),
		Height: node.lastcommit,
		Prove:  true,
	})

	require.True(t, qres.IsOK())

	var config ibc.ConnConfig
	node.App.Cdc.MustUnmarshalBinaryBare(qres.Value, &config)

	msg := MsgSpeakSafe{
		ChainID:           customid,
		RegisteredChainID: customid,
		Height:            uint64(qres.Height),
		Proof:             qres.Proof,
		RemoteConfig:      config,
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())

	return sdk.Result{}
}

func (node *node) init(t *testing.T) {
	msg := MsgInit{
		PortID: portid,
		Config: ibc.PortConfig{}, // TODO
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())
}

func (node *node) push(t *testing.T, message string) sdk.Result {
	msg := MsgPushMessage{
		ChainID: customid,
		PortID:  portid,
		Message: message,
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())

	return res
}

func (node *node) pull(t *testing.T, target *node, seq uint64, message string) sdk.Result {
	qres := target.Query(abci.RequestQuery{
		Path: "/store/ibc/key",
		Data: append(append(append(
			[]byte{0x00},
			customid[:]...),
			portid[:]...),
			store.EncodeIndex(seq, store.BinIndexerEnc)...,
		),
		Height: node.lastcommit,
		Prove:  true,
	})

	require.True(t, qres.IsOK())

	var packet ibc.Packet
	node.App.Cdc.MustUnmarshalBinaryBare(qres.Value, &packet)

	msg := ibc.MsgPull{
		ChainID: customid,
		PortID:  portid,
		Packets: []ibc.Packet{packet},
		Proof:   qres.Proof,
		Height:  uint64(qres.Height),
	}

	res := node.signCheckDeliver(t, msg, true)

	require.True(t, res.IsOK())

	return sdk.Result{}
}

func TestNode(t *testing.T) {
	node := getNode(t)
	node.startChain(t)

	verifier := lite.NewBaseVerifier("", 0, node.tmvalset())

	for i := 0; i < 10; i++ {
		node.execblock(t, nil)
		err := verifier.Verify(node.commits[int64(i)+1].SignedHeader)
		require.NoError(t, err)
	}
}

func TestBasicConn(t *testing.T) {
	node1, node2 := getNode(t), getNode(t)
	fmt.Println("node1 execblock")
	node1.execblock(t, nil)
	fmt.Println("node2 execblock")
	node2.execblock(t, nil)
	fmt.Println("node1 speak")
	node1.speakConn(t, node2)
	fmt.Println("node2 speak")
	node2.speakConn(t, node1)
	fmt.Println("node1 listen")
	node1.listenConn(t, node2)
	fmt.Println("node2 listen")
	node2.listenConn(t, node1)
	fmt.Println("node1 update")
	node1.update(t, node2)
	fmt.Println("node2 update")
	node2.update(t, node1)

	fmt.Println("node1 init")
	node1.init(t)
	fmt.Println("node2 init")
	node2.init(t)
	fmt.Println("[node1 ->] node2")
	node1.push(t, "testmessage")
	fmt.Println("node2 update")
	node2.update(t, node1)
	fmt.Println("node1 [-> node2]")
	node2.pull(t, node1, 0, "testmessage")
}

/*
func TestLifecycle(t *testing.T) {
	node := getNode(t)
	node.openConn(t)
	node.establishConn(t)
	node.openChannel(t)
	node.establishChannel(t)
	node.sendPackets(t)
	node.updateConn(t)
	node.sendPackets(t)
}

func TestLongDuration(t *testing.T) {
	node := getNode(t)
	node.openConn(t)
	node.establishConn(t)
	node.openChannel(t)
	node.establishChannel(t)

	for i := 0; i < 1000; i++ {
		node.sendPackets(t)
		node.updateConn(t)
	}
}*/
