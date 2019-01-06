package connection

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return handleMsgCheckpoint(ctx, k, msg)
		case MsgOpenConnection:
			return handleMsgOpenConnection(ctx, k, msg)
		case channel.MsgReceive:
			switch p := msg.Packets[0].Payload.(type) { // TODO: extend to range
			case channel.PayloadConnectionListening:
				return handlePayloadConnectionListening(ctx, k, msg.Packets[0].Source(), p)
			default:
				panic("return error") // TODO
			}
		default:
			panic("return error") // TODO
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, k Keeper, msg MsgCheckpoint) sdk.Result {
	k.checkpoint(ctx, uint64(ctx.BlockHeight()), ctx.BlockHeader())
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

func handlePayloadConnectionListening(ctx sdk.Context, k Keeper, source []byte, p channel.PayloadConnectionListening) sdk.Result {
	if !k.existsChain(ctx, source) {
		return types.ErrConnectionNotOpened(types.DefaultCodespace).Result()
	}

	if k.establishedChain(ctx, p.ChainID) {
		return types.ErrConnectionAlreadyOpened(types.DefaultCodespace).Result()
	}

	if !k.checkpointedConfig(ctx, p.Config.ROT.SignedHeader.Header) {
		return types.ErrHeaderNotCheckpointed(types.DefaultCodespace).Result()
	}

	// TODO: check encoding compatability

	return sdk.Result{}
}
