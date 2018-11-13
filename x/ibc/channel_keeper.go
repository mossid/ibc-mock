package connection

import (
	"github.com/tendermint/go-amino"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

func (core ChannelCore) store(ctx sdk.Context) sdk.KVStore {
	// TODO: switch this to RootMultiStore
	return store.NewPrefixStore(
		ctx.KVStore(core.k.key),
		[]byte(core.name),
	)
}

func QueueID(channelID []byte, queueID [QUEUEIDLEN]byte) (res []byte) {
	// Prefixing length to IDs to prevent collision
	if len(channelID) == 0 {
		res = append(res, amino.MustMarshalBinaryLengthPrefixed(channelID)...)
	}
	res = append(res, queueID[:]...)
	return
}

func ChannelQueuePrefix(channelID []byte, queueID [QUEUEIDLEN]byte) []byte {
	return append([]byte{0x00}, QueueID(channelID, queueID)...)
}

func (core ChannelCore) ChannelQueue(ctx sdk.Context, queueID [QUEUEIDLEN]byte) store.Queue {
	return store.NewQueue(
		core.k.cdc,
		store.NewPrefixStore(
			core.store(ctx),
			ChannelQueuePrefix(core.id, queueID),
		),
	)
}

func ChannelIncomingKey(channelID []byte, queueID [QUEUEIDLEN]byte) []byte {
	return append([]byte{0x01}, QueueID(channelID, queueID)...)
}

func (core ChannelCore) ChannelIncoming(ctx sdk.Context, queueID [QUEUEIDLEN]byte) store.Value {
	return store.NewValue(
		core.k.cdc,
		core.store(ctx),
		ChannelIncomingKey(core.id, queueID),
	)
}

type ChannelCore struct {
	k    Keeper
	name string
	id   []byte
}

func (k Keeper) ChannelCore(name string) ChannelCore {
	if name == "" {
		panic("cannot use empty channel name")
	}
	if _, ok := k.chs[name]; ok {
		panic("name already occupied")
	}

	return ChannelCore{
		k:    k,
		name: name,
		id:   nil,
	}
}

type ChannelTemplate struct {
	k    Keeper
	name string
}

func (k Keeper) ChannelTemplate(name string) ChannelTemplate {
	if name == "" {
		panic("cannot use empty channel name")
	}
	if _, ok := k.chs[name]; ok {
		panic("name already occupied")
	}

	return ChannelTemplate{
		k:    k,
		name: name,
	}
}

func (tem ChannelTemplate) Instanciate(id []byte) ChannelCore {
	if len(id) == 0 {
		panic("cannot use empty channel id")
	}

	return ChannelCore{
		k:    tem.k,
		name: tem.name,
		id:   id,
	}
}

func (core ChannelCore) Established(ctx sdk.Context) bool {

}

func (core ChannelCore) Open(ctx sdk.Context, chainID string) bool {
	ch := core.k.perConnectionChannel(ctx, chainID)
	ch.Send()
}

func (core ChannelCore) Send(ctx sdk.Context, pay Payload, queueID [QUEUEIDLEN]byte, destChain []byte) {
	queue := core.ChannelQueue(ctx, core.id, queueID)
	queue.Push(pay)
}
