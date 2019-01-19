package ibc

import (
	"fmt"

	"github.com/tendermint/iavl"
	"github.com/tendermint/tendermint/crypto/merkle"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewAnteHandler(k Keeper) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, result sdk.Result, abort bool) {
		prt := merkle.DefaultProofRuntime()
		prt.RegisterOpDecoder(iavl.ProofOpIAVLValue, iavl.IAVLValueOpDecoder)

		var flag bool
		var msg2 MsgPull
		var ok bool
		for _, msg := range tx.GetMsgs() {
			if msg2, ok = msg.(MsgPull); ok {
				if flag {
					// more than one MsgPull are in a tx
					// should not be happen, when a MsgPull fails
					// later MsgPulls cannot be executed(see baseapp.go L642)
					// but the sequence is still incremented
					fmt.Println("zzzz")
					return ctx, ErrAnteFailed(DefaultCodespace).Result(), true
				}
				flag = true
			}
		}

		if !ok {
			newCtx = ctx
			return
		}

		// Each MsgPull contains packets in a single channel
		// (one packet in this mock, since no rangeproof proofop impl)
		var config ConnConfig
		err := k.conn(msg2.ChainID).config.GetSafe(ctx, &config)
		if err != nil {
			fmt.Println("xxxx")
			return ctx, ErrAnteFailed(DefaultCodespace).Result(), true
		}

		var seq uint64
		k.port(msg2.PortID).queue(msg2.ChainID).incoming.GetIfExists(ctx, &seq)

		for _, packet := range msg2.Packets {
			keypath, err := config.mpath(msg2.PortID, seq)
			if err != nil {
				fmt.Println("vvvv")
				return ctx, ErrAnteFailed(DefaultCodespace).Result(), true
			}
			bz, err := k.cdc.MarshalBinaryBare(packet)
			if err != nil {
				fmt.Println("bbbb")
				return ctx, ErrAnteFailed(DefaultCodespace).Result(), true
			}
			fmt.Printf("antemarshal %+v %x\n", packet, bz)
			if !k.conn(msg2.ChainID).verify(ctx, keypath, msg2.Height, msg2.Proof, bz) {
				fmt.Println("nnnn")
				return ctx, ErrAnteFailed(DefaultCodespace).Result(), true
			}
			seq = seq + 1
		}

		k.port(msg2.PortID).queue(msg2.ChainID).incoming.Set(ctx, seq)

		newCtx = ctx
		return
	}
}
