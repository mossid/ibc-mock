package types

import (
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"github.com/tendermint/tendermint/lite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type (
	ChannelType = [CHANNELTYPELEN]byte
	QueueID     = [QUEUEIDLEN]byte
)

const (
	CHANNELTYPELEN = 4
	QUEUEIDLEN     = 2
)

func PacketQueueID() QueueID {
	return QueueID{0x00, 0x00}
}

func ReceiptQueueID() QueueID {
	return QueueID{0x00, 0x01}
}

type Channel interface {
	ChannelType() ChannelType
	//	Verify(queueid [QUEUEIDLEN]byte, payload Payload) bool
}

type ChannelInfo struct {
	ChannelType ChannelType
	KeyPath     merkle.KeyPath // KeyPath to the message queue
}

type ChannelInfoPair struct {
	Key  []byte
	Info interface{}
}

func (info ChannelInfo) Pairs(infomap store.Mapping) []ChannelInfoPair {
	return []ChannelInfoPair{
		{[]byte{0x00}, info.ChannelType},
		{[]byte{0x01}, info.KeyPath},
	}
}

type Packet struct {
	Header
	Payload
}

type Header struct {
	//	Route [][]byte
}

type Payload interface {
	Route() string
	Type() string
	ValidateBasic() sdk.Error
	GetSignBytes() []byte

	QueueID() QueueID
}

type Proof struct {
	Sequence uint64

	Height      int64
	SrcChain    []byte
	ChannelName string
	ChannelID   []byte
}

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
