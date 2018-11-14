package ibc

import (
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/merkle"

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

func (core ChannelCore) KeyPath() (res merkle.KeyPath) {
	res = new(merkle.KeyPath).
		AppendKey([]byte(core.k.key.Name()), merkle.KeyEncodingURL).
		AppendKey([]byte(core.name), merkle.KeyEncodingURL)

	if len(core.id) != 0 {
		res.AppendKey(amino.MustMarshalBinaryLengthPrefixed(core.id), merkle.KeyEncodingHex)
	}

	return
}

func (core ChannelCore) GetChannelInfo() ChannelInfo {
	return ChannelInfo{
		KeyPath: core.KeyPath(),
	}
}

func QueuePrefix(channelID []byte, queueID QueueID) (res []byte) {
	// Prefixing length to IDs to prevent collision
	if len(channelID) == 0 {
		res = append(res, amino.MustMarshalBinaryLengthPrefixed(channelID)...)
	}
	res = append(res, queueID[:]...)
	return
}

func ChannelQueuePrefix(channelID []byte, queueID QueueID) []byte {
	return append([]byte{0x00}, QueuePrefix(channelID, queueID)...)
}

func (core ChannelCore) ChannelQueue(ctx sdk.Context, queueID QueueID) store.Queue {
	return store.NewQueue(
		core.k.cdc,
		core.store(ctx),
		ChannelQueuePrefix(core.id, queueID),
		store.BinIndexerEnc,
	)
}

func ChannelIncomingKey(channelID []byte, queueID QueueID) []byte {
	return append([]byte{0x01}, QueuePrefix(channelID, queueID)...)
}

func (core ChannelCore) ChannelIncoming(ctx sdk.Context, queueID [QUEUEIDLEN]byte) store.Value {
	return store.NewValue(
		core.k.cdc,
		core.store(ctx),
		ChannelIncomingKey(core.id, queueID),
	)
}

// TODO: rename it to reduce confusion with GetChannelInfo
func ChannelInfoKey(channelID []byte, chainID []byte) (res []byte) {
	res = append(res, amino.MustMarshalBinaryLengthPrefixed(channelID)...)
	res = append(res, chainID...)
	return
}

func (core ChannelCore) ChannelInfo(ctx sdk.Context, chainID []byte) store.Value {
	return store.NewValue(
		core.k.cdc,
		core.store(ctx),
		ChannelInfoKey(core.id, chainID),
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

func (core ChannelCore) Established(ctx sdk.Context, chainID []byte) bool {
	return core.ChannelInfo(ctx, chainID).Exists()
}

func (ch PerConnectionChannel) Open(ctx sdk.Context, info ChannelInfo, id []byte) {
	pay := PayloadChannelListening{info}
	ch.core.send(ctx, pay, PacketQueueID(), id)
}

func (core ChannelCore) Open(ctx sdk.Context, chainID []byte) bool {
	if core.Established(ctx, chainID) {
		return false
	}

	core.k.perConnectionChannel(chainID).Open(ctx, core.GetChannelInfo(), chainID)

	return true
}

// send does not check establishment, used internally
func (core ChannelCore) send(ctx sdk.Context, pay Payload, queueID QueueID, destChain []byte) {
	if pay.Route() != core.name {
		panic("invalid payload")
	}

	queue := core.ChannelQueue(ctx, queueID)
	queue.Push(pay)
}

// Send checks the establishment of the channel before sending any payload
func (core ChannelCore) Send(ctx sdk.Context, pay Payload, queueID QueueID, destChain []byte) bool {
	if !core.Established(ctx, destChain) {
		return false
	}

	core.send(ctx, pay, queueID, destChain)

	return true
}
