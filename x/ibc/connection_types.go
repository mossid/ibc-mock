package ibc

import (
	"github.com/tendermint/tendermint/crypto/tmhash"
	lite "github.com/tendermint/tendermint/lite"
)

type EncodingScheme byte

const (
	Amino EncodingScheme = iota
	JSON
	ABI
)

type ConnectionConfig struct {
	ROT      lite.FullCommit
	Encoding EncodingScheme
}

func (config ConnectionConfig) UniqueID(userID string) []byte {
	hasher := tmhash.NewTruncated()
	hasher.Write(config.ROT.SignedHeader.Header.Hash()) // TODO: check safety
	hasher.Write([]byte(userID))
	return hasher.Sum(nil)
}

func (config ConnectionConfig) Height() uint64 {
	height := config.ROT.SignedHeader.Header.Height
	if height < 0 {
		panic("invalid header")
	}
	return uint64(height)
}

type PerConnectionChannel struct {
	core ChannelCore
}

func (PerConnectionChannel) ChannelType() [CHANNELTYPELEN]byte {
	return [CHANNELTYPELEN]byte{0x00, 0x00, 0x00, 0x00}
}
