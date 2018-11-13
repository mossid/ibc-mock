package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	Route = "ibc"
)

type MsgOpenConnection struct {
	UserChainID string           `json:"user-chain-id"`
	Config      ConnectionConfig `json:"config"`
	Signer      sdk.AccAddress   `json:"signer"`
}

var _ sdk.Msg = MsgOpenConnection{}

func (MsgOpenConnection) Route() string {
	return Route
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

type PayloadConnectionListening struct {
	Config  ConnectionConfig
	ChainID []byte
}

var _ Payload = PayloadConnectionListening{}

func (pay PayloadConnectionListening) Route() string {
	return "ibc"
}
