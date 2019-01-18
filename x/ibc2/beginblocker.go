package ibc

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const interval = 1 // TODO: parameterize

func BeginBlock(ctx sdk.Context, k Keeper) {
	if ctx.BlockHeight() == 0 {
		return // no header at the first block
	}
	header := ctx.BlockHeader()
	if header.Height%interval == 0 {
		fmt.Println("aaa", TMFromABCI(header), "bbb", TMFromABCI(header).Hash())
		k.checkpoints.Set(ctx, TMFromABCI(header).Hash(), []byte{})
	}
}
