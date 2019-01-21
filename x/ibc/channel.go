package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"

	sdk "github.com/cosmos/cosmos-sdk/types"
	//	"github.com/mossid/ibc-mock/store"
)

type ChannelUser struct {
	Module  string
	Address sdk.AccAddress
}

var _ User = ChannelUser{}

func (chuser ChannelUser) Satisfy(req User) bool {
	chreq, ok := req.(ChannelUser)
	if !ok {
		return false
	}
	if chuser.Module != chreq.Module {
		return false
	}
	if chuser.Address == nil {
		return true
	}
	return chuser.Address.Equals(chreq.Address)
}

type Channel struct {
	portID PortID
	route  string
	k      Keeper
}

func (k Keeper) Channel(route string) (res Channel) {
	res = Channel{
		route: route,
		k:     k,
	}
	copy(res.portID[:], tmhash.NewTruncated().Sum([]byte(route)))
	return
}

func (ch Channel) user(user sdk.AccAddress) ChannelUser {
	return ChannelUser{ch.route, user}
}

func (ch Channel) PortID() PortID {
	return ch.portID
}

func (ch Channel) Speak(ctx sdk.Context, id ChainID, owner sdk.AccAddress) bool {
	return ch.k.conn(id).speak(ctx, ch.user(owner))
}

/*
func (ch Channel) Speak(ctx sdk.Context,
	rot lite.FullCommit, root string,
	customID string, owner sdk.AccAddress,
) bool {
	id := UniqueID(rot.SignedHeader.Header, customID)
	if ch.k.Conn(id).owner.Exists(ctx) {
		return false
	}
	if !ch.k.Conn(id).speak(ctx) {
		return false
	}
	ch.k.Conn(id).owner.Set(ctx, ChannelOwner{ch.route, owner})
	return true
}
*/

func (ch Channel) Listen(ctx sdk.Context, id ChainID, config ConnConfig, user sdk.AccAddress) bool {
	return ch.k.conn(id).listen(ctx, config, ch.user(user))
}

func (ch Channel) SpeakSafe(ctx sdk.Context, id ChainID, height uint64, proof *merkle.Proof, remoteconfig ConnConfig, user sdk.AccAddress) bool {
	return ch.k.conn(id).speakSafe(ctx, height, proof, remoteconfig, ch.user(user))
}

func (ch Channel) ListenSafe(ctx sdk.Context, id ChainID) bool {
	return ch.k.conn(id).listenSafe(ctx)
}

func (ch Channel) Open(ctx sdk.Context, id ChainID) bool {
	return true
}

func (ch Channel) Init(ctx sdk.Context, id PortID, config PortConfig) bool {
	return ch.k.port(id).init(ctx, config)
}

func (ch Channel) Push(ctx sdk.Context, chainID ChainID, p Packet) bool {
	if p.Route() != ch.route {
		panic("invalid payload")
	}
	return ch.k.port(ch.portID).push(ctx, chainID, p)
}

func (ch Channel) Update(ctx sdk.Context, id ChainID, source Source, user sdk.AccAddress) bool {
	return ch.k.conn(id).update(ctx, source, ch.user(user))
}

type StdChannel struct {
	Packet  Channel
	Receipt Channel
}
