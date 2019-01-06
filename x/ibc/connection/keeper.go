package connection

import (
	"bytes"

	abci "github.com/tendermint/tendermint/abci/types"
	cmn "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
	ch "github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

type remote struct {
	config  store.Value          // ChainConfig
	height  store.Value          // int64
	chainid store.Value          // []byte
	commits store.Indexer        // height int64 -> lite.FullCommit
	valsets store.Indexer        // height int64 -> tmtypes.ValidatorSet
	channel func([]byte) channel // channelID []byte -> channel
}

type local struct {
	checkpoints store.Indexer // height uint64 -> SimpleHeader
	config      store.Value   // ChainConfig
}

type channel struct {
	keypath store.Value // string
}

type Keeper struct {
	// keeper base elements
	cdc *codec.Codec
	key sdk.StoreKey

	// storage accessor
	remote func([]byte) remote // chainID []byte -> remote
	local  local

	// local multistore
	chstore *store.MultiStore

	// external reference
	perconn ch.PerConnectionTemplate
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, valset ValidatorSet) (k Keeper) {
	base := store.NewBase(cdc, key)

	k = Keeper{
		cdc: cdc,
		key: key,

		remote: func(chainID []byte) remote {
			base2 := base.Prefix(chainID)
			return remote{
				config:  store.NewValue(base2, []byte{0x00}),
				height:  store.NewValue(base2, []byte{0x01}),
				chainid: store.NewValue(base2, []byte{0x02}),
				commits: store.NewIndexer(base2, []byte{0x03}, store.BinIndexerEnc),
				valsets: store.NewIndexer(base2, []byte{0x04}, store.BinIndexerEnc),
			}
		},
		local: local{
			checkpoints: store.NewIndexer(base, []byte{0x20}, store.BinIndexerEnc),
			config:      store.NewValue(base, []byte{0x21}),
		},

		// Implementing as store.MultiStore(which is just a prefixstore) for now
		// TODO: change it to rootmultistore
		chstore: store.NewMultiStore(key),
	}

	tem := k.ChannelTemplate("ibc", ch.PerConnectionChannelType())
	k.perconn = ch.PerConnectionTemplate{&tem}

	return k
}

// -----------------------------------------
// Handler methods

func (k Keeper) checkpoint(ctx sdk.Context, height uint64, header abci.Header) {
	k.local.checkpoints.Set(ctx, height, simpleFromABCI(header))
}

func (k Keeper) existsChain(ctx sdk.Context, id []byte) bool {
	return k.remote(id).config.Exists(ctx)
}

func (k Keeper) registerChain(ctx sdk.Context, id []byte, config types.ConnectionConfig) {

	k.remote(id).config.Set(ctx, config)
	k.remote(id).height.Set(ctx, config.Height())
	k.remote(id).commits.Set(ctx, config.Height(), config.ROT)

	// Sending
	k.perconn.Channel(id).OpenConnection(ctx, config, id)
}

func (k Keeper) establishedChain(ctx sdk.Context, id []byte) bool {
	return k.remote(id).chainid.Exists(ctx)
}

func (k Keeper) checkpointedConfig(ctx sdk.Context, rot *tmtypes.Header) bool {
	var stored SimpleHeader
	if rot.Height < 0 { // TODO: move to ValidateBasic
		return false
	}
	k.local.checkpoints.Get(ctx, uint64(rot.Height), &stored)
	return simpleFromTmtypes(*rot).Equal(stored)
}

func (k Keeper) establishChain(ctx sdk.Context, id []byte) {
	k.remote(id).chainid.Set(ctx, id)
}

func (k Keeper) lastHeader(ctx sdk.Context, id []byte) SimpleHeader {
	var height int64
	k.remote(id).height.Get(ctx, &height)
	var fc lite.FullCommit
	k.remote(id).commits.Get(ctx, uint64(height), &fc)
	return simpleFromTmtypes(*fc.SignedHeader.Header)
}

// ---------------------------------------
//

type SimpleHeader struct {
	ChainID            string
	Height             int64
	ValidatorsHash     cmn.HexBytes
	NextValidatorsHash cmn.HexBytes
	AppHash            cmn.HexBytes
}

func (h1 SimpleHeader) Equal(h2 SimpleHeader) bool {
	return (h1.ChainID == h2.ChainID &&
		h1.Height == h2.Height &&
		bytes.Equal(h1.ValidatorsHash, h2.ValidatorsHash) &&
		bytes.Equal(h1.NextValidatorsHash, h2.NextValidatorsHash) &&
		bytes.Equal(h1.AppHash, h2.AppHash))
}

func simpleFromABCI(header abci.Header) SimpleHeader {
	return SimpleHeader{
		ChainID:            header.ChainID,
		Height:             header.Height,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		AppHash:            header.AppHash,
	}
}

func simpleFromTmtypes(header tmtypes.Header) SimpleHeader {
	return SimpleHeader{
		ChainID:            header.ChainID,
		Height:             header.Height,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		AppHash:            header.AppHash,
	}
}
