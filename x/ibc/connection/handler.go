package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return handleMsgCheckpoint(ctx, k, msg)
		case MsgOpenConnection:
			return handleMsgOpenConnection(ctx, k, msg)
		default:
			panic("return error") // TODO
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, k Keeper, msg MsgCheckpoint) sdk.Result {
	k.checkpointer.lastvalset.Snapshot(ctx)
	k.checkpointer.checkpoints.Set(ctx, uint64(ctx.BlockHeight()), ctx.BlockHeader())
	return sdk.Result{}
}

func handleMsgOpenConnection(ctx sdk.Context, k Keeper, msg MsgOpenConnection) sdk.Result {
	id := msg.Config.UniqueID(msg.UserChainID)
	if k.existsChain(ctx, id) {
		return types.ErrConnectionAlreadyOpened(types.DefaultCodespace).Result()
	}
	k.registerChain(ctx, id, msg.Config)

	return sdk.Result{} // TODO: add tags
}

var (
	// TODO: move it to parameter store
	DIFFREQ = sdk.NewDec(5).Quo(sdk.NewDec(6)) // 2/3 + (1/3)/2
)
