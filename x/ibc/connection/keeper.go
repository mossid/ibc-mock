package connection

import (
	//	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
	ch "github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

type remote struct {
	config  store.Mapping // chainID []byte -> ChainConfig
	height  store.Mapping // chainID []byte -> int64
	commits store.Indexer // (chainID []byte prefix) -> height int64 -> lite.FullCommit
}

type local struct {
	config      store.Value   // ChainConfig
	checkpoints store.Indexer // height uint64 -> lite.FullCommit
}

type Keeper struct {
	cdc *codec.Codec
	key sdk.StoreKey

	remote remote
	local  local

	chstore *store.MultiStore

	perconn ch.PerConnectionTemplate
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) (k Keeper) {
	base := store.NewBase(cdc, key)

	// Implementing as prefixstore for now
	// TODO: change it to rootmultistore
	k = Keeper{
		cdc: cdc,
		key: key,

		remote: remote{
			config:  store.NewMapping(base, []byte{0x00}),
			height:  store.NewMapping(base, []byte{0x01}),
			commits: store.NewIndexer(base, []byte{0x02}, store.BinIndexerEnc),
		},
		local: local{
			config:      store.NewValue(base, []byte{0x20}),
			checkpoints: store.NewIndexer(base, []byte{0x21}, store.BinIndexerEnc),
		},

		chstore: store.NewMultiStore(key),
	}

	tem := k.ChannelTemplate("ibc", ch.PerConnectionChannelType())
	k.perconn = ch.PerConnectionTemplate{&tem}

	return k
}

// -----------------------------------------
// Handler methods

func (k Keeper) existsChain(ctx sdk.Context, id []byte) bool {
	return k.remote.config.Has(ctx, id)
}

func (k Keeper) registerChain(ctx sdk.Context, id []byte, config types.ConnectionConfig) {
	k.remote.config.Set(ctx, id, config)
	k.remote.height.Set(ctx, id, config.Height())
	k.remote.commits.Prefix(id).Set(ctx, config.Height(), config.ROT)

	// Sending
	k.perconn.Channel(id).OpenConnection(ctx, config, id)
}

// ---------------------------------------
//
