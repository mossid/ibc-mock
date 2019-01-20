package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Packet interface {
	Route() string
	Type() string
	ValidateBasic() sdk.Error
	GetSignBytes() []byte
}
