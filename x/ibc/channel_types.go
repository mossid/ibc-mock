package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
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

type Payload interface {
	Route() string
	Type() string
	ValidateBasic() sdk.Error
	GetSignBytes() []byte

	QueueID() QueueID
}
