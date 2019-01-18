package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ChainID = []byte
type PortID = []byte

const (
	QueueIdle        Status = iota // Port is not opened
	QueueEstablished               // Compatability checked, ready to use
)

type PortConfig interface {
	RequiredStatus() Status
	IsCompatibleWith(PortConfig) bool
}

func (p port) push(ctx sdk.Context, chainID ChainID, data []byte) bool {
	if !assertStatus(ctx, p.queue(chainID).status, QueueEstablished) {
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

func (p port) establish(
	ctx sdk.Context,
	chainID ChainID, proof *merkle.Proof, remoteconfig PortConfig) bool {
	if !transitStatus(ctx, p.queue(chainID).status, QueueIdle, QueueEstablished) {
		return false
	}

	conn := p.queue(chainID).conn
	if !assertStatus(ctx, conn.status, ConnEstablished) {
		return false
	}

	var config PortConfig
	p.config.Get(ctx, &config)
	if !config.IsCompatibleWith(remoteconfig) {
		return false
	}

	return true
}
