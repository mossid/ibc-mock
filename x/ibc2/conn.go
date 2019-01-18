package ibc

import (
	"fmt"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ConnIdle  Status = iota // not opened
	ConnOpen                // ready to send packets
	ConnReady               // ready to receive packets
)

type ConnConfig struct {
	ROT         lite.FullCommit // Root-of-trust fullcommit from the other chain
	ChainID     ChainID         // ChainID of this chain registered on the other chain
	RootKeyPath string          // Root keypath of ibc module on the other chain
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

	fmt.Printf("before\n%+v\n%x\n", config.RootKeyPath, key)
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

	fmt.Printf("after\n%+v\n", res)

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

func (c conn) open(ctx sdk.Context, rot lite.FullCommit, root string) bool {
	if !transitStatus(ctx, c.status, ConnIdle, ConnOpen) {
		return false
	}

	fmt.Println("vvv", root)
	c.config.Set(ctx, ConnConfig{
		ROT:         rot,
		RootKeyPath: root,
	})
	c.commits.Set(ctx, uint64(rot.Height()), rot)

	return true
}

func (c conn) update(ctx sdk.Context, source Source) bool {
	if !assertStatus(ctx, c.status, ConnOpen, ConnReady) {
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

func (c conn) ready(ctx sdk.Context,
	registeredChainID ChainID, proof *merkle.Proof, remoteconfig ConnConfig,
) bool {
	if !transitStatus(ctx, c.status, ConnOpen, ConnReady) {
		fmt.Println("ppp11")
		return false
	}

	var config ConnConfig
	c.config.Get(ctx, &config)
	fmt.Println("fff", config)
	config.ChainID = registeredChainID

	cpath, err := config.cpath()
	if err != nil {
		fmt.Println("aaa11", err)
		return false
	}
	bz, err := c.config.Cdc().MarshalBinaryBare(remoteconfig)
	if err != nil {
		fmt.Println("bbb11", err)
		return false
	}

	if !c.verify(ctx, cpath, proof, bz) {
		fmt.Println("ccc11")
		return false
	}

	// set root-of-trust commit
	c.config.Set(ctx, config)
	c.commits.Set(ctx, config.Height(), config.ROT)

	// TODO: check compat config.rot

	return true
}

func (c conn) latestcommit(ctx sdk.Context) (index uint64, res lite.FullCommit, ok bool) {
	index, ok = c.commits.Last(ctx, &res)
	return
}

func (c conn) verify(ctx sdk.Context, keypath merkle.KeyPath, proof *merkle.Proof, value []byte) bool {
	fmt.Printf("verify, %+v, %+v, %+v\n", keypath, proof, value)

	var commit lite.FullCommit
	c.commits.Last(ctx, &commit)

	fmt.Println("lastcommit", commit)

	prt := merkle.DefaultProofRuntime()
	prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)
	prt.RegisterOpDecoder(store.ProofOpMultiStore, store.MultiStoreProofOpDecoder)
	err := prt.VerifyValue(proof, commit.SignedHeader.AppHash, keypath.String(), value)
	fmt.Println(err)
	return err == nil
}
