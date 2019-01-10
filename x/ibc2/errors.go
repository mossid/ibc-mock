package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = "ibc"

	CodeEmptyMsg          CodeType = 101
	CodePacketsNotUniform CodeType = 102
	CodeConnOpenFailed    CodeType = 103
	CodeConnUpdateFailed  CodeType = 104
)

func ErrEmptyMsg(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyMsg, "")
}

func ErrPacketsNotUniform(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodePacketsNotUniform, msg)
}

func ErrConnOpenFailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeConnOpenFailed, "")
}

func ErrConnUpdateFailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeConnUpdateFailed, "")
}
