package connection

import (
	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/channel"
	"github.com/mossid/ibc-mock/x/ibc/types"
)

func NewAnteHandler(k Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, result sdk.Result, abort bool) {
		prt := merkle.DefaultProofRuntime()
		prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)

		for _, msg := range tx.GetMsgs() {
			if msg2, ok := msg.(channel.MsgReceive); ok {
				// Each MsgReceive contains packets in a single channel
				// (a packet in this mock, since no rangeproof proofop impl)
				for _, packet := range msg2.Packets {
					var config types.ConnectionConfig
					k.remote(packet.Source).config.Get(ctx, &config)
					keypath := config.KeyPath
					msg2.Proof.Verify(prt, k.lastHeader(ctx, packet.Source()).AppHash, keypath)
				}
			}
		}

		newCtx = ctx
		return
	}
}
