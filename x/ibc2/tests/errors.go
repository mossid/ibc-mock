package ibctest

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = "ibc"

	CodeEmptyMsg          CodeType = 101
	CodePacketsNotUniform CodeType = 102
	CodeConnSpeakFailed   CodeType = 103
	CodeConnUpdateFailed  CodeType = 104
	CodeConnListenFailed  CodeType = 105
)

func ErrEmptyMsg(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyMsg, "")
}

func ErrPacketsNotUniform(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodePacketsNotUniform, msg)
}

func ErrConnSpeakFailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeConnSpeakFailed, "")
}

func ErrConnUpdateFailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeConnUpdateFailed, "")
}

func ErrConnListenFailed(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeConnListenFailed, "")
}
