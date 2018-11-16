package channel

import (
	"github.com/tendermint/go-amino"
	//	"github.com/tendermint/tendermint/crypto/merkle"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

/*
func (core ChannelCore) KeyPath() (res merkle.KeyPath) {
	res = new(merkle.KeyPath).
		AppendKey([]byte(core.k.key.Name()), merkle.KeyEncodingURL).
		AppendKey([]byte(core.name), merkle.KeyEncodingURL)

	if len(core.id) != 0 {
		res.AppendKey(amino.MustMarshalBinaryLengthPrefixed(core.id), merkle.KeyEncodingHex)
	}

	return
}
*/
func (core ChannelCore) GetChannelInfo() types.ChannelInfo {
	return types.ChannelInfo{
		ChannelType: core.ty,
		//		KeyPath:     core.KeyPath(),
	}
}

type remote struct {
	info   store.Mapping // (chainID []byte prefix) -> (infokey []byte) -> info T
	height store.Mapping // (chainID []byte prefix) -> (queueID QueueID) -> height uint64
}

type local struct {
	info    store.Mapping // (infotype []byte) -> info <T>
	packets store.Queue   // (chainID []byte prefix) -> (queueID QueueID prefix) -> sequence uint64 -> Packet
}

// (channelName []byte prefix) -> (channelID []byte optional prefix) -> ChannelCore
type ChannelCore struct {
	remote remote
	local  local

	perconn PerConnectionTemplate
	name    string
	ty      types.ChannelType
}

func NewChannelCore(base store.Base, perconn PerConnectionTemplate, name string, ty types.ChannelType) ChannelCore {
	return ChannelCore{
		remote: remote{
			height: store.NewMapping(base, []byte{0x00}),
		},
		local: local{
			info:    store.NewMapping(base, []byte{0x10}),
			packets: store.NewQueue(base, []byte{0x11}, store.BinIndexerEnc),
		},
		perconn: perconn,

		name: name,
		ty:   ty,
	}

}

type ChannelTemplate struct {
	cdc     *codec.Codec
	rkey    sdk.StoreKey
	perconn PerConnectionTemplate
	key     sdk.StoreKey
	ty      types.ChannelType
}

func NewChannelTemplate(cdc *codec.Codec, rkey sdk.StoreKey, perconn PerConnectionTemplate, key sdk.StoreKey, ty types.ChannelType) ChannelTemplate {
	return ChannelTemplate{
		cdc:     cdc,
		rkey:    rkey,
		perconn: perconn,
		key:     key,
		ty:      ty,
	}
}

func (tem ChannelTemplate) ChannelCore(id []byte) ChannelCore {
	if len(id) == 0 {
		panic("cannot use empty channel id")
	}

	base := store.NewBaseWithAccessor(
		tem.cdc,
		func(ctx sdk.Context) sdk.KVStore {
			return store.NewPrefixStore(
				ctx.KVStore(tem.rkey).(*store.MultiStore).GetKVStore(tem.key)(ctx), // TODO: change to multistore
				amino.MustMarshalBinaryLengthPrefixed(id),
			)
		},
	)

	return NewChannelCore(base, tem.perconn, tem.key.Name(), tem.ty)
}

func (core ChannelCore) instantiated(ctx sdk.Context) bool {
	return core.local.info.Has(ctx, []byte{})
}

func (core ChannelCore) instantiate(ctx sdk.Context) {
	infomap := core.local.info
	infomap.Set(ctx, []byte{}, []byte{})
	/*
		for _, pair := range core.GetChannelInfo().Pairs() {
			infomap.Set(ctx, pair.Key, pair.Info)
		}*/
}

func (core ChannelCore) established(ctx sdk.Context, chainID []byte) bool {
	return core.remote.info.Has(ctx, chainID)
}

func (core ChannelCore) Open(ctx sdk.Context, chainID []byte) bool {
	if core.established(ctx, chainID) {
		return false
	}

	if !core.instantiated(ctx) {
		core.instantiate(ctx)
	}

	core.perconn.Channel(chainID).OpenChannel(ctx, core.GetChannelInfo(), chainID)

	return true
}

// send does not check establishment, used internally
func (core ChannelCore) send(ctx sdk.Context, pay types.Payload, queueID types.QueueID, destChain []byte) {
	if pay.Route() != core.name {
		panic("invalid payload")
	}

	core.local.packets.Prefix(queueID[:]).Push(ctx, pay)
}

// Send checks the establishment of the channel before sending any payload
func (core ChannelCore) Send(ctx sdk.Context, pay types.Payload, queueID types.QueueID, destChain []byte) bool {
	if !core.established(ctx, destChain) {
		return false
	}

	core.send(ctx, pay, queueID, destChain)

	return true
}
