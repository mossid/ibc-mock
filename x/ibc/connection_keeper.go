package ibc

import (
	"github.com/tendermint/go-amino"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type Keeper struct {
	key sdk.StoreKey
	cdc *codec.Codec

	// Using pointer to avoid recursive type
	perconn *ChannelTemplate

	chs map[string]ChannelCore
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	// Implementing as prefixstore for now
	// TODO: change it to rootmultistore
	k := Keeper{
		key: key,
		cdc: cdc,
	}

	perconn := k.ChannelTemplate("ibc")
	k.perconn = &perconn

	return k
}

// ---------------------------------------
// Remote Chain management keys

func ChainConfigKey(srcChain []byte) []byte {
	return append([]byte{0x00}, []byte(srcChain)...)
}

func (k Keeper) chainConfig(ctx sdk.Context, srcChain []byte) store.Value {
	return store.NewValue(
		k.cdc,
		ctx.KVStore(k.key),
		ChainConfigKey(srcChain),
	)
}

func ChainHeightKey(srcChain []byte) []byte {
	return append([]byte{0x01}, []byte(srcChain)...)
}

func (k Keeper) chainHeight(ctx sdk.Context, srcChain []byte) store.Value {
	return store.NewValue(
		k.cdc,
		ctx.KVStore(k.key),
		ChainHeightKey(srcChain),
	)
}

func ChainCommitsPrefix(srcChain []byte) []byte {
	return append([]byte{0x02}, amino.MustMarshalBinaryLengthPrefixed(srcChain)...)
}

func (k Keeper) chainCommits(ctx sdk.Context, srcChain []byte) store.Indexer {
	return store.NewIndexer(
		k.cdc,
		ctx.KVStore(k.key),
		ChainCommitsPrefix(srcChain),
		store.BinIndexerEnc,
	)
}

// ----------------------------------------
// Remote chain per connection channel
func (k Keeper) perConnectionChannel(destChain []byte) PerConnectionChannel {
	return PerConnectionChannel{
		core: k.perconn.Instanciate(destChain),
	}
}

// ----------------------------------------
// Local chain management keys
func CheckpointPrefix() []byte {
	return []byte{0x20}
}

func (k Keeper) checkpointList(ctx sdk.Context) store.Indexer {
	return store.NewIndexer(
		k.cdc,
		ctx.KVStore(k.key),
		CheckpointPrefix(),
		store.BinIndexerEnc,
	)
}

// -----------------------------------------
// Handler methods

func (k Keeper) existsChain(ctx sdk.Context, id []byte) bool {
	return k.chainConfig(ctx, id).Exists()
}

func (ch PerConnectionChannel) Register(ctx sdk.Context, config ConnectionConfig, id []byte) {
	pay := PayloadConnectionListening{config, id}
	ch.core.send(ctx, pay, PacketQueueID(), id)
}

func (k Keeper) registerChain(ctx sdk.Context, id []byte, config ConnectionConfig) {
	k.chainConfig(ctx, id).Set(config)
	k.chainHeight(ctx, id).Set(config.Height())
	k.chainCommits(ctx, id).Set(config.Height(), config.ROT)

	// Sending
	k.perConnectionChannel(id).register(ctx, config, id)
}

// ---------------------------------------
//
