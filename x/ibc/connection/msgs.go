package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

type MsgOpenConnection struct {
	UserChainID string                 `json:"user-chain-id"`
	Config      types.ConnectionConfig `json:"config"`
	Signer      sdk.AccAddress         `json:"signer"`
}

var _ sdk.Msg = MsgOpenConnection{}

func (MsgOpenConnection) Route() string {
	return "ibc"
}

func (MsgOpenConnection) Type() string {
	return "open_connection"
}

func (msg MsgOpenConnection) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (msg MsgOpenConnection) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg) // TODO
}

func (msg MsgOpenConnection) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}
