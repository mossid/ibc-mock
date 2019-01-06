package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Payload interface {
	Route() string
	Type() string
	ValidateBasic() sdk.Error
	GetSignBytes() []byte
}
