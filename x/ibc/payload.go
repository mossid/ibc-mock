package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PacketQueueOpen struct {
	Config PortConfig
}

var _ Packet = PacketQueueOpen{}

func (p PacketQueueOpen) Route() string {
	return "ibc"
}

func (p PacketQueueOpen) Type() string {
	return "payload_queue_ready"
}

func (p PacketQueueOpen) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (p PacketQueueOpen) GetSignBytes() []byte {
	return msgCdc.MustMarshalJSON(p) // TODO
}

type PacketConnOpen struct {
	Config ConnConfig
}

var _ Packet = PacketConnOpen{}

func (p PacketConnOpen) Route() string {
	return "ibc"
}

func (p PacketConnOpen) Type() string {
	return "payload_conn_open"
}

func (p PacketConnOpen) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (p PacketConnOpen) GetSignBytes() []byte {
	return msgCdc.MustMarshalJSON(p) // TODO
}

type PacketConnReady struct {
	Config ConnConfig
}

var _ Packet = PacketConnReady{}

func (p PacketConnReady) Route() string {
	return "ibc"
}

func (p PacketConnReady) Type() string {
	return "payload_conn_ready"
}

func (p PacketConnReady) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (p PacketConnReady) GetSignBytes() []byte {
	return msgCdc.MustMarshalJSON(p) // TODO
}
