package ibc

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var cdc *codec.Codec

func init() {
	RegisterCodec(cdc)
}

func RegisterCodec(cdc *codec.Codec) {

}
