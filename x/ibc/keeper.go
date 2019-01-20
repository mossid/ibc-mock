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
	id PortID

	// state structure is enforced
	config store.Value
	queue  func(ChainID) queue

	cdc *codec.Codec
}

func PortPrefix() []byte {
	id := ChainID{}
	return id[:]
}

type User interface {
	Satisfy(User) bool
}

type conn struct {
	id ChainID

	// state structure is enforced
	config store.Value // ChainConfig

	// state structure is not enforced
	status  store.Value
	commits store.Indexer
	owner   store.Value

	// internal reference
	localconfig store.Value

	cdc *codec.Codec
}

type Keeper struct {
	// Message sending, state structure is enforced
	config store.Value
	conn   func(ChainID) conn
	port   func(PortID) port

	cdc *codec.Codec
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey) Keeper {
	/*
		if msgCdc != nil {
			if cdc != msgCdc {
				panic("NewKeeper() called with different cdc")
			}
		}
	*/
	// XXX: unsafe
	msgCdc = cdc
	base := store.NewBase(cdc, key)
	return newKeeper(cdc, base.Prefix([]byte{0x00}), base.Prefix([]byte{0x01}))
}

func DummyKeeper() Keeper {
	return newKeeper(nil, store.NewBase(nil, nil), store.NewBase(nil, nil))
}

func newKeeper(cdc *codec.Codec, enfbase store.Base, freebase store.Base) (k Keeper) {
	// structure-free keys
	cmbase := freebase.Prefix([]byte{0x00}) // commit
	csbase := freebase.Prefix([]byte{0x01}) // conn status
	qsbase := freebase.Prefix([]byte{0x02}) // queue status
	cobase := freebase.Prefix([]byte{0x03}) // conn ownership

	k = Keeper{
		config: store.NewValue(enfbase, []byte{}),
		conn: func(chainID ChainID) conn {
			bz := chainID[:]
			return conn{
				id:          chainID,
				config:      store.NewValue(enfbase, bz),
				status:      store.NewValue(csbase, bz),
				commits:     store.NewIndexer(cmbase, bz, store.BinIndexerEnc),
				owner:       store.NewValue(cobase, bz),
				localconfig: k.config,
				cdc:         cdc,
			}
		},
		port: func(portID PortID) port {
			return port{
				id:     portID,
				config: store.NewValue(enfbase, append(PortPrefix(), portID[:]...)),
				queue: func(chainID ChainID) queue {
					bz := bytes.Join([][]byte{chainID[:], portID[:]}, nil)
					return queue{
						outgoing: store.NewQueue(enfbase, bz, store.BinIndexerEnc),
						incoming: store.NewValue(enfbase, bz),
						status:   store.NewValue(qsbase, bz),
						conn:     k.conn(chainID),
					}
				},
				cdc: cdc,
			}
		},
		cdc: cdc,
	}

	return

}

// keeperConfig is stored in Keeper.Config, working as a magic number for IBC root
type keeperConfig struct {
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
