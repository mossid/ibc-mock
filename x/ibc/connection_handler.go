package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgOpenConnection:
			return handleMsgOpenConnection(ctx, k, msg)
		default:
			panic("aaa") // TODO
		}
	}
}

func handleMsgOpenConnection(ctx sdk.Context, k Keeper, msg MsgOpenConnection) sdk.Result {
	id := msg.Config.UniqueID(msg.UserChainID)
	if k.existsChain(ctx, id) {
		return sdk.Result{} // TODO: error
	}
	k.registerChain(ctx, id, msg.Config)

	return sdk.Result{} // TODO: add tags
}
