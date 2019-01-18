package ibc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgOpen:
			return handleMsgOpen(ctx, k, msg)
		case MsgReady:
			return handleMsgReady(ctx, k, msg)
		case MsgEstablish:
			return handleMsgEstablish(ctx, k, msg)
		case MsgUpdate:
			return handleMsgUpdate(ctx, k, msg)
		default:
			return sdk.ErrUnknownRequest("aaa").Result()
		}
	}
}

func handleMsgOpen(ctx sdk.Context, k Keeper, msg MsgOpen) sdk.Result {
	id := UniqueID(msg.ROT.SignedHeader.Header, msg.CustomChainID)

	fmt.Printf("id %x\n", id)

	if !k.conn(id).open(ctx, msg.ROT) {
		return ErrConnOpenFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handleMsgReady(ctx sdk.Context, k Keeper, msg MsgReady) sdk.Result {
	fmt.Printf("id %x\n", msg.Config.ChainID)

	if !k.conn(msg.Config.ChainID).ready(ctx,
		msg.RootKeyPath, msg.ChainID, msg.Proof,
		msg.Config,
	) {
		return ErrConnReadyFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handleMsgEstablish(ctx sdk.Context, k Keeper, msg MsgEstablish) sdk.Result {
	// TODO?
	return sdk.Result{}
}

func handleMsgUpdate(ctx sdk.Context, k Keeper, msg MsgUpdate) sdk.Result {
	source := NewSource(msg.Commits[0].SignedHeader.Header, msg.Commits)

	fmt.Printf("id %x\n", msg.ChainID)

	if !k.conn(msg.ChainID).update(ctx, source) {
		return ErrConnUpdateFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}
