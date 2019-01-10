package ibc

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type State struct {
	QueuePrefix      []byte
	ConnConfigPrefix []byte
	PortConfigPrefix []byte
}

func DefaultState() State {
	return State{
		QueuePrefix:      []byte{0x00},
		ConnConfigPrefix: []byte{0x01},
		PortConfigPrefix: []byte{0x02},
	}
}

func (s State) Bases(base store.Base) (store.Base, store.Base, store.Base) {
	return base.Prefix(s.QueuePrefix), base.Prefix(s.ConnConfigPrefix), base.Prefix(s.PortConfigPrefix)
}

type queue struct {
	// state structure is enforced
	outgoing store.Queue

	// state structure is not enforced
	incoming store.Value
	status   store.Value

	// internal reference
	conn conn
}

type port struct {
	// state structure is enforced
	config store.Value
	queue  func(ChainID) queue

	meta queue
}

type conn struct {
	// state structure is enforced
	config store.Value // ChainConfig

	// state structure is not enforced
	status  store.Value
	commits store.Indexer

	// internal reference
	localconfig store.Value

	// misc
	cdc *codec.Codec
}

type Keeper struct {
	// Message sending, state structure is enforced
	config store.Value
	conn   func(ChainID) conn
	port   func(PortID) port

	// misc
	state State
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, state State) (k Keeper) {
	base := store.NewBase(cdc, key)

	// structure-enforced keys
	enfbase := base.Prefix([]byte{0x00})
	qbase, cbase, pbase := state.Bases(enfbase)

	// structure-free keys
	freebase := base.Prefix([]byte{0x01})
	cmbase := freebase.Prefix([]byte{0x00}) // commit
	csbase := freebase.Prefix([]byte{0x01}) // conn status
	qsbase := freebase.Prefix([]byte{0x02}) // queue status

	k = Keeper{
		config: store.NewValue(enfbase, []byte{}),
		conn: func(chainID ChainID) conn {
			return conn{
				config:      store.NewValue(cbase, chainID[:]),
				status:      store.NewValue(csbase, chainID[:]),
				commits:     store.NewIndexer(cmbase, chainID[:], store.BinIndexerEnc),
				localconfig: k.config,
				cdc:         cdc,
			}
		},
		port: func(portID PortID) port {
			return port{
				config: store.NewValue(pbase, portID[:]),
				queue: func(chainID ChainID) queue {
					bz := bytes.Join([][]byte{chainID[:], portID[:]}, nil)
					return queue{
						outgoing: store.NewQueue(qbase, bz, store.BinIndexerEnc),
						incoming: store.NewValue(qbase, bz),
						status:   store.NewValue(qsbase, bz),
						conn:     k.conn(chainID),
					}
				},
			}
		},
		state: state,
	}

	return

}

func DummyKeeper(state State) Keeper {
	return NewKeeper(nil, nil, state)
}

/*
func (k Keeper) VerifyPacketProof(
	ctx sdk.Context,
	chainID ChainID, queueID QueueID,
	packet []byte, proof merkle.ProofOperators,
) bool {
	var config QueueConfig
	k.chain(chainID).queue(queueID).config.Get(ctx, &config)

	var incoming uint64
	k.conn(chainID).queue(queueID).incoming.Get(ctx, &incoming)

	var root []byte
	k.conn(chainID).commits.Last(ctx, &root)

	return proof.VerifyValue(root, config.KeyPath.String(), packet) == nil
}

func (k Keeper) PushPacket(
	ctx sdk.Context,
	chainID ChainID, queueID QueueID,
	packet []byte,
) {
	k.conn(chainID).queue(queueID).outgoing.Push(ctx, packet)
}

func (k Keeper) QueueStatus(
	ctx sdk.Context,
	chainID ChainID, queueID QueueID,
) (status QueueStatus) {
	k.conn(chainID).queue(queueID).status.GetIfExists(ctx, &status)
	return
}*/
