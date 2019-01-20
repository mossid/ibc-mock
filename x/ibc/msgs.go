package ibc

import (
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgPull struct {
	ChainID ChainID
	Height  uint64
	PortID  PortID
	Packets []Packet       `json:"packets"`
	Proof   *merkle.Proof  `json:"proof"`
	Relayer sdk.AccAddress `json:"relayer"`
}

func (msg MsgPull) Route() string {
	return msg.Packets[0].Route()
}

func (msg MsgPull) Type() string {
	return msg.Packets[0].Type()
}

func (msg MsgPull) ValidateBasic() sdk.Error {
	if len(msg.Packets) == 0 {
		return ErrEmptyMsg(DefaultCodespace)
	}

	p0 := msg.Packets[0]
	route0, ty0 := p0.Route(), p0.Type()

	for _, p := range msg.Packets[1:] {
		if p.Route() != route0 {
			return ErrPacketsNotUniform(DefaultCodespace, "invalid route")
		}
		if p.Type() != ty0 {
			return ErrPacketsNotUniform(DefaultCodespace, "invalid type")
		}
	}

	return nil
}

func (msg MsgPull) GetSignBytes() []byte {
	return msgCdc.MustMarshalJSON(msg)
}

func (msg MsgPull) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Relayer}
}
