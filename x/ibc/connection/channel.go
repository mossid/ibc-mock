package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
	ch "github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

func (k Keeper) ChannelCore(name string, ty types.ChannelType) ch.ChannelCore {
	if name == "" {
		panic("cannot use empty ch.name")
	}

	key := sdk.NewKVStoreKey(name)
	k.chstore.MountStore(key)

	base := store.NewBaseWithAccessor(
		k.cdc,
		func(ctx sdk.Context) sdk.KVStore {
			// TODO: change to ctx.MultiStore()
			return ctx.KVStore(k.key).(*store.MultiStore).GetKVStore(key)(ctx)
		},
	)

	return ch.NewChannelCore(base, k.perconn, name, ty)
}

func (k Keeper) ChannelTemplate(name string, ty types.ChannelType) ch.ChannelTemplate {
	if name == "" {
		panic("cannot use empty ch.name")
	}

	key := sdk.NewKVStoreKey(name)
	k.chstore.MountStore(key)

	return ch.NewChannelTemplate(k.cdc, k.key, k.perconn, key, ty)
}
