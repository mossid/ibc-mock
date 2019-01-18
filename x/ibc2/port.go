package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChainID = []byte
type PortID = []byte

const (
	QueueIdle  Status = iota // Port is not opened
	QueueReady               // Compatability checked, ready to use
)

type PortConfig interface {
	RequiredStatus() Status
	IsCompatibleWith(PortConfig) bool
}

func (p port) push(ctx sdk.Context, chainID ChainID, data []byte) bool {
	if !assertStatus(ctx, p.queue(chainID).status, QueueReady) {
		return false
	}

	if !assertStatus(ctx, p.queue(chainID).conn.status, ConnOpen, ConnReady) {
		return false
	}

	p.queue(chainID).outgoing.Push(ctx, data)

	return true
}

// Called when the internal logic deployes initializes new port
func (p port) init(ctx sdk.Context, config PortConfig) bool {
	if p.config.Exists(ctx) {
		return false
	}

	p.config.Set(ctx, config)

	return true
}

func (p port) ready(
	ctx sdk.Context,
	chainID ChainID, proof *merkle.Proof, remoteconfig PortConfig) bool {
	if !transitStatus(ctx, p.queue(chainID).status, QueueIdle, QueueReady) {
		return false
	}

	if !assertStatus(ctx, p.queue(chainID).conn.status, ConnOpen, ConnReady) {
		return false
	}

	var pconfig PortConfig
	err := p.config.GetSafe(ctx, &pconfig)
	if err != nil {
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

	if !p.queue(chainID).conn.verify(ctx, path, proof, bz) {
		return false
	}

	if !pconfig.IsCompatibleWith(remoteconfig) {
		return false
	}

	return true
}
