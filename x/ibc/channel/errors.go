package channel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "ibc"
)

func ErrEmptyMsg(codespace sdk.CodespaceType) sdk.Error {
	return nil
}

func ErrUnmatchingPacket(codespace sdk.CodespaceType, msg string) sdk.Error {
	return nil
}
