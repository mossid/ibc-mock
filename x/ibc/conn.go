package ibc

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// not opened
	ConnIdle Status = iota
	// ready to send packets
	ConnSpeak
	// ready to receive packets
	ConnListen
	// counterchain is tracking us
	// ready to safely send packets
	ConnSpeakSafe
	// counterchain approved we are tracking them
	// ready to safely receive packets
	ConnListenSafe
)

type ConnConfig struct {
	ROT         lite.FullCommit // Root-of-trust fullcommit from the other chain
	ChainID     ChainID         // ChainID of this chain registered on the other chain
	RootKeyPath string          // Root keypath of ibc module on the other chain
}

func (config ConnConfig) IsEmpty() bool {
	return len(config.ChainID) == 0
}

func (config ConnConfig) Height() uint64 {
	height := config.ROT.SignedHeader.Header.Height
	if height < 0 {
		panic("invalid header")
	}
	return uint64(height)
}

func UniqueID(header *tmtypes.Header, customID string) []byte {
	return ConcatHash(header.Hash(), []byte(customID))
}

func (config ConnConfig) UniqueID(customID string) []byte {
	return UniqueID(config.ROT.SignedHeader.Header, customID)
}

func (config ConnConfig) extendKeyPath(key []byte) (res merkle.KeyPath, err error) {
	// TODO: optimize
	keys, err := merkle.KeyPathToKeys(config.RootKeyPath)
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

func (c conn) assertOwner(ctx sdk.Context, user User) bool {
	if user == nil {
		return true
	}
	var req User
	c.owner.Get(ctx, &req)
	return user.Satisfy(req)
}

func (c conn) speak(ctx sdk.Context, owner User) bool {
	if !transitStatus(ctx, c.status, ConnIdle, ConnSpeak) {
		return false
	}

	c.owner.Set(ctx, owner)

	return true
}

func (c conn) listen(ctx sdk.Context, config ConnConfig, user User) bool {
	if !c.assertOwner(ctx, user) {
		return false
	}
	if !transitStatus(ctx, c.status, ConnSpeak, ConnListen) {
		return false
	}

	c.config.Set(ctx, config)
	c.commits.Set(ctx, uint64(config.ROT.Height()), config.ROT)

	return true
}

func (c conn) update(ctx sdk.Context, source Source, user User) bool {
	if !c.assertOwner(ctx, user) {
		return false
	}
	if !assertStatus(ctx, c.status, ConnSpeak, ConnListen) {
		return false
	}

	shdr := source.GetLastSignedHeader()

	provider := provider{c, ctx, []byte(shdr.ChainID)}

	verifier := lite.NewDynamicVerifier(
		shdr.ChainID,
		provider,
		source,
	)

	err := verifier.Verify(shdr)
	if err != nil {
		return false
	}

	source.saveAll(provider)
	return true
}

func (c conn) speakSafe(ctx sdk.Context, height uint64, proof *merkle.Proof, remoteconfig ConnConfig, user User) bool {
	if !c.assertOwner(ctx, user) {
		return false
	}
	if !transitStatus(ctx, c.status, ConnListen, ConnSpeakSafe) {
		return false
	}

	var config ConnConfig
	c.config.Get(ctx, &config)

	cpath, err := config.cpath()
	if err != nil {
		return false
	}
	bz, err := c.config.Cdc().MarshalBinaryBare(remoteconfig)
	if err != nil {
		return false
	}

	if !c.verify(ctx, cpath, height, proof, bz) {
		return false
	}

	// TODO: check compat config.rot

	// TODO: push approve payload to perconn

	return true
}

func (c conn) listenSafe(ctx sdk.Context) bool {
	if !transitStatus(ctx, c.status, ConnSpeakSafe, ConnListenSafe) {
		return false
	}

	return true
}

func (c conn) latestcommit(ctx sdk.Context) (index uint64, res lite.FullCommit, ok bool) {
	index, ok = c.commits.Last(ctx, &res)
	return
}

func (c conn) verify(ctx sdk.Context, keypath merkle.KeyPath, height uint64, proof *merkle.Proof, value []byte) bool {
	var commit lite.FullCommit
	err := c.commits.GetSafe(ctx, height, &commit)
	if err != nil {
		return false
	}

	prt := merkle.DefaultProofRuntime()
	prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)
	prt.RegisterOpDecoder(store.ProofOpMultiStore, store.MultiStoreProofOpDecoder)
	err = prt.VerifyValue(proof, commit.SignedHeader.AppHash, keypath.String(), value)
	return err == nil
}
