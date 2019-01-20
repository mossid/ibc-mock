package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = "ibc"

	CodeEmptyMsg          CodeType = 101
	CodePacketsNotUniform CodeType = 102
	CodeAnteFailed        CodeType = 103
)

func ErrEmptyMsg(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyMsg, "")
}

func ErrPacketsNotUniform(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodePacketsNotUniform, msg)
}

func ErrAnteFailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAnteFailed, "")
}
