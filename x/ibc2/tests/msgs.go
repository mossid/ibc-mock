package ibctest

import (
	"github.com/tendermint/tendermint/crypto/merkle"
	"github.com/tendermint/tendermint/lite"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc2"
)

const Route = "ibctest"

type MsgSpeak struct {
	ChainID ibc.ChainID `json:"chain-id"`
	Signer  sdk.AccAddress
}

var _ sdk.Msg = MsgSpeak{}

func (MsgSpeak) Route() string {
	return Route
}

func (MsgSpeak) Type() string {
	return "speak"
}

func (msg MsgSpeak) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgSpeak) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg) // TODO
}

func (msg MsgSpeak) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgListen struct {
	ChainID ibc.ChainID
	Config  ibc.ConnConfig
	Signer  sdk.AccAddress `json:"signer"`
}

var _ sdk.Msg = MsgListen{}

func (MsgListen) Route() string {
	return Route
}

func (MsgListen) Type() string {
	return "open"
}

func (msg MsgListen) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgListen) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg) // TODO
}

func (msg MsgListen) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgUpdate struct {
	ChainID ibc.ChainID
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

type MsgInit struct {
	PortID ibc.PortID
	Config ibc.PortConfig
	Signer sdk.AccAddress
}

var _ sdk.Msg = MsgUpdate{}

func (MsgInit) Route() string {
	return Route
}

func (MsgInit) Type() string {
	return "init"
}

func (msg MsgInit) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (msg MsgInit) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgInit) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgSpeakSafe struct {
	ChainID           ibc.ChainID
	RegisteredChainID ibc.ChainID
	Height            uint64
	Proof             *merkle.Proof
	Signer            sdk.AccAddress
	RemoteConfig      ibc.ConnConfig
}

var _ sdk.Msg = MsgSpeakSafe{}

func (MsgSpeakSafe) Route() string {
	return Route
}

func (MsgSpeakSafe) Type() string {
	return "ready"
}

func (msg MsgSpeakSafe) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgSpeakSafe) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgSpeakSafe) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgPushMessage struct {
	Message string
	ChainID ibc.ChainID
	PortID  ibc.PortID
	Signer  sdk.AccAddress
}

func (MsgPushMessage) Route() string {
	return Route
}

func (MsgPushMessage) Type() string {
	return "push_message"
}

func (msg MsgPushMessage) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgPushMessage) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}

func (msg MsgPushMessage) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type PacketPushMessage struct {
	Message string
}

var _ ibc.Packet = PacketPushMessage{}

func (PacketPushMessage) Route() string {
	return Route
}

func (PacketPushMessage) Type() string {
	return "push_message"
}

func (msg PacketPushMessage) ValidateBasic() sdk.Error {
	return nil
}

func (msg PacketPushMessage) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(msg)
}
