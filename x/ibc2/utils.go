package ibc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type Status = byte

func assertStatus(ctx sdk.Context, value store.Value, expected Status) bool {
	var actual Status
	value.GetIfExists(ctx, &actual)
	return expected == actual
}

func transitStatus(ctx sdk.Context, value store.Value, from Status, to Status) bool {
	if !assertStatus(ctx, value, from) {
		return false
	}
	value.Set(ctx, to)
	return true
}
