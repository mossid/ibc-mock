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

	// LC verification, state structure is not enforced
	checkpoints store.Indexer // height int64 -> SimpleHeader
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) (k Keeper) {
	base := store.NewBase(cdc, key)

	// structure-enforced keys
	qbase := base.Prefix([]byte{0x00}) // message
	cbase := base.Prefix([]byte{0x01}) // chain config
	pbase := base.Prefix([]byte{0x02}) // port config

	// structure-free keys
	ibase := base.Prefix([]byte{0xA0})  // queue incoming sequence
	cmbase := base.Prefix([]byte{0xA1}) // commit
	ckbase := base.Prefix([]byte{0xA2}) // checkpoint
	sbase := base.Prefix([]byte{0xA3})  // status

	k = Keeper{
		config: store.NewValue(base, []byte{}),
		conn: func(chainID ChainID) conn {
			return conn{
				config:      store.NewValue(cbase, chainID[:]),
				status:      store.NewValue(sbase, chainID[:]),
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
						incoming: store.NewValue(ibase, bz),
						status:   store.NewValue(sbase, bz),
						conn:     k.conn(chainID),
					}
				},
			}
		},
		checkpoints: store.NewIndexer(ckbase, []byte{}, store.BinIndexerEnc),
	}

	return

}

func DummyKeeper() Keeper {
	return NewKeeper(nil, nil)
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
