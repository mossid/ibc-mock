package ibc

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ConnStatus = byte

const (
	ConnIdle        ConnStatus = iota // not opened
	ConnOpen                          // waiting for request
	ConnReady                         // waiting for response
	ConnEstablished                   // ready to send packets
)

type ConnConfig struct {
	ROT         lite.FullCommit // Root-of-trust fullcommit from the other chain
	ChainID     ChainID         // ChainID of this chain registered on the other chain
	RootKeyPath merkle.KeyPath  // Root keypath of ibc module on the other chain
}

func (config ConnConfig) Height() uint64 {
	height := config.ROT.SignedHeader.Header.Height
	if height < 0 {
		panic("invalid header")
	}
	return uint64(height)
}

func (config ConnConfig) extendKeyPath(key []byte) (res merkle.KeyPath, err error) {
	// TODO: optimize
	keys, err := merkle.KeyPathToKeys(config.RootKeyPath.String())
	if err != nil {
		return
	}

	// simulating key prefixing
	keys[len(keys)-1] = append(keys[len(keys)-1], key...)

	for _, key := range keys {
		res = res.AppendKey(key, merkle.KeyEncodingHex)
	}

	return
}

func (config ConnConfig) mpath(portID PortID, index uint64) (merkle.KeyPath, error) {
	k := DummyKeeper()
	return config.extendKeyPath(k.port(portID).queue(config.ChainID).outgoing.Key(index))
}

func (config ConnConfig) cpath() (merkle.KeyPath, error) {
	k := DummyKeeper()
	return config.extendKeyPath(k.conn(config.ChainID).config.Key())
}

func (config ConnConfig) ppath(portID PortID) (merkle.KeyPath, error) {
	k := DummyKeeper()
	return config.extendKeyPath(k.port(portID).config.Key())
}

func (c conn) transitStatus(ctx sdk.Context, from, to ConnStatus) bool {
	var current ConnStatus
	c.status.GetIfExists(ctx, &current)
	if current != from {
		return false
	}
	c.status.Set(ctx, to)
	return true
}

func (c conn) open(ctx sdk.Context, rot lite.FullCommit) bool {
	if !c.transitStatus(ctx, ConnIdle, ConnOpen) {
		return false
	}

	c.config.Set(ctx, ConnConfig{ROT: rot})

	return true
}

func (c conn) update(ctx sdk.Context, source Source, shdr tmtypes.SignedHeader) bool {
	verifier := lite.NewDynamicVerifier(
		shdr.ChainID,
		provider{c, ctx, []byte(shdr.ChainID)},
		source,
	)

	err := verifier.Verify(shdr)
	return err == nil
}

func (c conn) ready(ctx sdk.Context,
	root merkle.KeyPath, proof *merkle.Proof,
	chainID ChainID, remoteconfig ConnConfig,
) bool {
	if !c.transitStatus(ctx, ConnOpen, ConnReady) {
		return false
	}

	var config ConnConfig
	c.config.Get(ctx, &config)
	config.RootKeyPath = root
	config.ChainID = chainID

	var commit lite.FullCommit
	c.commits.Last(ctx, &commit)

	prt := merkle.DefaultProofRuntime()
	prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)
	cpath, err := config.cpath()
	if err != nil {
		return false
	}
	bz, err := c.cdc.MarshalBinaryBare(remoteconfig)
	if err != nil {
		return false
	}
	err = prt.VerifyValue(proof, commit.SignedHeader.AppHash, cpath.String(), bz)
	if err != nil {
		return false
	}

	// set root-of-trust commit
	c.config.Set(ctx, config)
	c.commits.Set(ctx, config.Height(), config.ROT)

	/*
		// XXX
		// send payload via perconn queue
		c.meta.push(ctx, cdc.MustMarshalBinaryBare(PayloadConnReady{config}))
	*/

	return true
}

func (c conn) establish(ctx sdk.Context, p PayloadConnReady) bool {
	if !c.transitStatus(ctx, ConnReady, ConnEstablished) {
		return false
	}

	return true
}

func (c conn) verify(keypath merkle.KeyPath, proof merkle.Proof, value []byte) bool {
	// TODO
	return false
}
