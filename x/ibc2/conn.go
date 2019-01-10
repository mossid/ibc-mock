package ibc

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ConnIdle        Status = iota // not opened
	ConnOpen                      // waiting for request
	ConnReady                     // waiting for response
	ConnEstablished               // ready to send packets
)

type ConnConfig struct {
	ROT         lite.FullCommit // Root-of-trust fullcommit from the other chain
	ChainID     ChainID         // ChainID of this chain registered on the other chain
	RootKeyPath merkle.KeyPath  // Root keypath of ibc module on the other chain
	State       State
}

func (config ConnConfig) Height() uint64 {
	height := config.ROT.SignedHeader.Header.Height
	if height < 0 {
		panic("invalid header")
	}
	return uint64(height)
}

func UniqueID(header *tmtypes.Header, customID string) []byte {
	hasher := tmhash.NewTruncated()
	hasher.Write(header.Hash())
	hasher.Write([]byte(customID))
	return hasher.Sum(nil)
}

func (config ConnConfig) UniqueID(customID string) []byte {
	return UniqueID(config.ROT.SignedHeader.Header, customID)
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
	k := DummyKeeper(config.State)
	return config.extendKeyPath(k.port(portID).queue(config.ChainID).outgoing.Key(index))
}

func (config ConnConfig) cpath() (merkle.KeyPath, error) {
	k := DummyKeeper(config.State)
	return config.extendKeyPath(k.conn(config.ChainID).config.Key())
}

func (config ConnConfig) ppath(portID PortID) (merkle.KeyPath, error) {
	k := DummyKeeper(config.State)
	return config.extendKeyPath(k.port(portID).config.Key())
}

func (c conn) open(ctx sdk.Context, rot lite.FullCommit) bool {
	if !transitStatus(ctx, c.status, ConnIdle, ConnOpen) {
		return false
	}

	c.config.Set(ctx, ConnConfig{ROT: rot})
	c.commits.Set(ctx, uint64(rot.Height()), rot)

	return true
}

func (c conn) update(ctx sdk.Context, source Source) bool {
	if !assertStatus(ctx, c.status, ConnOpen, ConnReady, ConnEstablished) {
		return false
	}

	shdr := source.GetLastSignedHeader()

	verifier := lite.NewDynamicVerifier(
		shdr.ChainID,
		provider{c, ctx, []byte(shdr.ChainID)},
		source,
	)

	err := verifier.Verify(shdr)
	return err == nil
}

func (c conn) ready(ctx sdk.Context,
	root merkle.KeyPath, registeredChainID ChainID, proof *merkle.Proof,
	remoteconfig ConnConfig,
) bool {
	if !transitStatus(ctx, c.status, ConnOpen, ConnReady) {
		return false
	}

	var config ConnConfig
	c.config.Get(ctx, &config)
	config.RootKeyPath = root
	config.ChainID = registeredChainID

	cpath, err := config.cpath()
	if err != nil {
		return false
	}
	bz, err := c.cdc.MarshalBinaryBare(remoteconfig)
	if err != nil {
		return false
	}

	if !c.verify(ctx, cpath, proof, bz) {
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
	if !transitStatus(ctx, c.status, ConnReady, ConnEstablished) {
		return false
	}

	return true
}

func (c conn) latestcommit(ctx sdk.Context) (index uint64, res lite.FullCommit, ok bool) {
	index, ok = c.commits.Last(ctx, &res)
	return
}

func (c conn) verify(ctx sdk.Context, keypath merkle.KeyPath, proof *merkle.Proof, value []byte) bool {
	var commit lite.FullCommit
	c.commits.Last(ctx, &commit)

	prt := merkle.DefaultProofRuntime()
	prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)
	err := prt.VerifyValue(proof, commit.SignedHeader.AppHash, keypath.String(), value)
	return err == nil
}
