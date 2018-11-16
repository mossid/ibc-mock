package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = 100 //TODO
)

func ErrEmptyMsg(codespace sdk.CodespaceType) sdk.Error {
	return nil
}

func ErrUnmatchingPacket(codespace sdk.CodespaceType, msg string) sdk.Error {
	return nil
}
