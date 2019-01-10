package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/lite"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const Route = "ibc"

type MsgOpen struct {
	CustomChainID string          `json:"custom-chain-id"`
	ROT           lite.FullCommit `json:"rot"`
	Signer        sdk.AccAddress  `json:"signer"`
}

var _ sdk.Msg = MsgOpen{}

func (MsgOpen) Route() string {
	return Route
}

func (MsgOpen) Type() string {
	return "open"
}

func (msg MsgOpen) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgOpen) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg) // TODO
}

func (msg MsgOpen) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgUpdate struct {
	ChainID ChainID
	Commits []lite.FullCommit
	Signer  sdk.AccAddress
}

var _ sdk.Msg = MsgUpdate{}

func (MsgUpdate) Route() string {
	return Route
}

func (MsgUpdate) Type() string {
	return "update"
}

func (msg MsgUpdate) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (msg MsgUpdate) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgUpdate) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgReady struct {
	ChainID     ChainID
	RootKeyPath merkle.KeyPath
	Proof       *merkle.Proof
	Signer      sdk.AccAddress
}

var _ sdk.Msg = MsgReady{}

func (MsgReady) Route() string {
	return Route
}

func (MsgReady) Type() string {
	return "ready"
}

func (msg MsgReady) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgReady) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgReady) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgEstablish struct {
	Signer sdk.AccAddress
}

var _ sdk.Msg = MsgEstablish{}

func (MsgEstablish) Route() string {
	return Route
}

func (MsgEstablish) Type() string {
	return "establish"
}

func (msg MsgEstablish) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgEstablish) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgEstablish) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgReceive struct {
	Packets []Packet       `json:"packets"`
	Proof   Proof          `json:"proof"`
	Relayer sdk.AccAddress `json:"relayer"`
}

func (msg MsgReceive) Route() string {
	return msg.Packets[0].Payload.Route()
}

func (msg MsgReceive) Type() string {
	return msg.Packets[0].Payload.Type()
}

func (msg MsgReceive) ValidateBasic() sdk.Error {
	if len(msg.Packets) == 0 {
		return ErrEmptyMsg(DefaultCodespace)
	}

	p0 := msg.Packets[0]
	route0, ty0 := p0.Payload.Route(), p0.Type()

	for _, p := range msg.Packets[1:] {
		if p.Payload.Route() != route0 {
			return ErrPacketsNotUniform(DefaultCodespace, "invalid route")
		}
		if p.Type() != ty0 {
			return ErrPacketsNotUniform(DefaultCodespace, "invalid type")
		}
	}

	return nil
}

func (msg MsgReceive) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgReceive) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}
