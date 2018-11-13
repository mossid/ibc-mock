package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CHANNELTYPELEN = 4
	QUEUEIDLEN     = 2
)

type Channel interface {
	ChannelType() [CHANNELTYPELEN]byte
	//	Verify(queueid [QUEUEIDLEN]byte, payload Payload) bool
}

type Payload interface {
	Route() string
	Type() string
	ValidateBasic() sdk.Error
	GetSignBytes() []byte

	QueueID() [QUEUEIDLEN]byte
}
