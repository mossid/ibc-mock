package ibc

/*
import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)


func NewAnteHandler(k Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, result sdk.Result, abort bool) {
		for _, msg := range tx.GetMsgs() {
			if msg, ok := msg.(MsgReceive); ok { // TODO: MsgRelay
				if !k.VerifyProof(ctx, msg.Proof) {
					return ctx, result, true
				}
			}
		}
		return ctx, result, false
	}
}
*/
