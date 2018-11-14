package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	return Route
}

func (pay PayloadConnectionListening) Type() string {
	return "payload_connection_listening"
}

func (pay PayloadConnectionListening) ValidateBasic() sdk.Error {
	return nil
}

func (pay PayloadConnectionListening) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(pay) //TODO
}

func (pay PayloadConnectionListening) QueueID() QueueID {
	return PacketQueueID()
}
