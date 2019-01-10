package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Header struct {
	Source      []byte
	Destination []byte
}

type Payload interface {
	Route() string
	Type() string
	ValidateBasic() sdk.Error
	GetSignBytes() []byte
}

type Proof struct {
	Key   []byte
	Value []byte
	Proof *merkle.Proof
}

type Packet struct {
	Header
	Payload
}
