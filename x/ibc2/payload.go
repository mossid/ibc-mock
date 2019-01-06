package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type PayloadQueueOpen struct {
	Config PortConfig
}

var _ Payload = PayloadQueueOpen{}

func (p PayloadQueueOpen) Route() string {
	return "ibc"
}

func (p PayloadQueueOpen) Type() string {
	return "payload_queue_ready"
}

func (p PayloadQueueOpen) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (p PayloadQueueOpen) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(p) // TODO
}

type PayloadConnOpen struct {
	Config ConnConfig
}

var _ Payload = PayloadConnOpen{}

func (p PayloadConnOpen) Route() string {
	return "ibc"
}

func (p PayloadConnOpen) Type() string {
	return "payload_conn_open"
}

func (p PayloadConnOpen) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (p PayloadConnOpen) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(p) // TODO
}

type PayloadConnReady struct {
	Config ConnConfig
}

var _ Payload = PayloadConnReady{}

func (p PayloadConnReady) Route() string {
	return "ibc"
}

func (p PayloadConnReady) Type() string {
	return "payload_conn_ready"
}

func (p PayloadConnReady) ValidateBasic() sdk.Error {
	return nil // TODO
}

func (p PayloadConnReady) GetSignBytes() []byte {
	return cdc.MustMarshalJSON(p) // TODO
}
