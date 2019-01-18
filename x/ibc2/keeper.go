package ibc

import (
	"bytes"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type queue struct {
	// state structure is enforced
	outgoing store.Queue
	incoming store.Value

	// state structure is not enforced
	status store.Value

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

	// state structure is not enforced
	checkpoints store.Mapping // tmtypes.Header.Hash() -> bool
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	base := store.NewBase(cdc, key)
	return newKeeper(base.Prefix([]byte{0x00}), base.Prefix([]byte{0x01}))
}

func DummyKeeper() Keeper {
	return newKeeper(store.NewBase(nil, nil), store.NewBase(nil, nil))
}

func newKeeper(enfbase store.Base, freebase store.Base) (k Keeper) {
	// structure-free keys
	cmbase := freebase.Prefix([]byte{0x00}) // commit
	csbase := freebase.Prefix([]byte{0x01}) // conn status
	qsbase := freebase.Prefix([]byte{0x02}) // queue status
	ckbase := freebase.Prefix([]byte{0x03}) // checkpoints

	k = Keeper{
		config: store.NewValue(enfbase, []byte{}),
		conn: func(chainID ChainID) conn {
			bz := chainID[:]
			return conn{
				config:      store.NewValue(enfbase, bz),
				status:      store.NewValue(csbase, bz),
				commits:     store.NewIndexer(cmbase, bz, store.BinIndexerEnc),
				localconfig: k.config,
				cdc:         cdc,
			}
		},
		port: func(portID PortID) port {
			return port{
				config: store.NewValue(enfbase, append(ChainID{}, portID[:]...)),
				queue: func(chainID ChainID) queue {
					bz := bytes.Join([][]byte{chainID[:], portID[:]}, nil)
					return queue{
						outgoing: store.NewQueue(enfbase, bz, store.BinIndexerEnc),
						incoming: store.NewValue(enfbase, bz),
						status:   store.NewValue(qsbase, bz),
						conn:     k.conn(chainID),
					}
				},
			}
		},
		checkpoints: store.NewMapping(ckbase, []byte{}),
	}

	return

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
