package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChainID = [8]byte
type PortID = [8]byte

const (
	QueueIdle  Status = iota // ready to send if portconfig exists
	QueueReady               // compatability checked, ready to receive
)

/*
type PortConfig interface {
	RequiredStatus() Status
	IsCompatibleWith(PortConfig) bool
}
*/
type PortConfig struct {
}

// Called when the internal logic deployes initializes new port
func (p port) init(ctx sdk.Context, config PortConfig) bool {
	if p.config.Exists(ctx) {
		return false
	}

	p.config.Set(ctx, config)

	return true
}

func (p port) push(ctx sdk.Context, chainID ChainID, o interface{}) bool {
	if !p.config.Exists(ctx) {
		return false
	}

	if !assertMinStatus(ctx, p.queue(chainID).conn.status, ConnSpeak) {
		return false
	}

	p.queue(chainID).outgoing.Push(ctx, o)

	return true
}

func (p port) ready(
	ctx sdk.Context,
	chainID ChainID, height uint64, proof *merkle.Proof, remoteconfig PortConfig,
) bool {
	var pconfig PortConfig
	err := p.config.GetSafe(ctx, &pconfig)
	if err != nil {
		return false
	}

	if !transitStatus(ctx, p.queue(chainID).status, QueueIdle, QueueReady) {
		return false
	}

	bz, err := p.config.Cdc().MarshalBinaryBare(remoteconfig)
	if err != nil {
		return false
	}

	var cconfig ConnConfig
	err = p.config.GetSafe(ctx, &cconfig)
	if err != nil {
		return false
	}

	path, err := cconfig.ppath(p.id)
	if err != nil {
		return false
	}

	if !p.queue(chainID).conn.verify(ctx, path, height, proof, bz) {
		return false
	}

	// TODO
	/*
		if !pconfig.IsCompatibleWith(remoteconfig) {
			return false
		}
	*/

	return true
}

func (p port) cleanup(ctx sdk.Context, chainID ChainID, proof *merkle.Proof, seq uint64) bool {
	if !assertMinStatus(ctx, p.queue(chainID).status, QueueReady) {
		return false
	}

	return false // TODO
}

/*
func (p port) pull(ctx sdk.Context, chainID ChainID) bool {
	q := p.queue(chainID)

	if !assertMinStatus(ctx, q.conn.status, ConnListen) {
		return false
	}

	if !assertMinStatus(ctx, q.status, QueueReady) {
		return false
	}



	return true
}
*/
