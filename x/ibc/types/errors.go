package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = 100 //TODO

	CodeEmptyMsg                CodeType = 101
	CodeUnmatchingPacket        CodeType = 102
	CodeConnectionAlreadyOpened CodeType = 103
)

func ErrEmptyMsg(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyMsg, "")
}

func ErrUnmatchingPacket(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeUnmatchingPacket, msg)
}

func ErrConnectionAlreadyOpened(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeConnectionAlreadyOpened, "")
}
