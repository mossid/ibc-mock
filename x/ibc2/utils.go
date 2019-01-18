package ibc

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type Status = byte

func assertStatus(ctx sdk.Context, value store.Value, expected ...Status) bool {
	var actual Status
	value.GetIfExists(ctx, &actual)
	fmt.Println("actual", actual)
	for _, exp := range expected {
		if exp == actual {
			return true
		}
	}
	return false
}

func transitStatus(ctx sdk.Context, value store.Value, from Status, to Status) bool {
	fmt.Printf("transit %+v %v %v", value, from, to)
	if !assertStatus(ctx, value, from) {
		return false
	}
	value.Set(ctx, to)
	return true
}

func TMFromABCI(header abci.Header) *tmtypes.Header {
	return &tmtypes.Header{
		Version: version.Consensus{
			Block: version.Protocol(header.Version.Block),
			App:   version.Protocol(header.Version.App),
		},
		ChainID:  header.ChainID,
		Height:   header.Height,
		Time:     header.Time,
		NumTxs:   header.NumTxs,
		TotalTxs: header.TotalTxs,
		LastBlockID: tmtypes.BlockID{
			Hash: header.LastBlockId.Hash,
			PartsHeader: tmtypes.PartSetHeader{
				Total: int(header.LastBlockId.PartsHeader.Total),
				Hash:  header.LastBlockId.PartsHeader.Hash,
			},
		},
		LastCommitHash:     header.LastCommitHash,
		DataHash:           header.DataHash,
		ValidatorsHash:     header.ValidatorsHash,
		NextValidatorsHash: header.NextValidatorsHash,
		ConsensusHash:      header.ConsensusHash,
		AppHash:            header.AppHash,
		LastResultsHash:    header.LastResultsHash,
		EvidenceHash:       header.EvidenceHash,
		ProposerAddress:    header.ProposerAddress,
	}
}
