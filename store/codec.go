package store

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

type cdc struct {
	cdc *codec.Codec
}

func (cdc cdc) MarshalBinaryBare(o interface{}) ([]byte, error) {
	return cdc.cdc.MarshalBinaryBare(o)
}

func (cdc cdc) UnmarshalBinaryBare(bz []byte, ptr interface{}) error {
	return cdc.cdc.UnmarshalBinaryBare(bz, ptr)
}
