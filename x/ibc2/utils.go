package ibc

import (
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tendermint/tendermint/version"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/store"
)

type Status = byte

func assertMinStatus(ctx sdk.Context, value store.Value, minexp Status) bool {
	var actual Status
	value.GetIfExists(ctx, &actual)
	return actual >= minexp
}

func assertStatus(ctx sdk.Context, value store.Value, expected ...Status) bool {
	var actual Status
	value.GetIfExists(ctx, &actual)
	for _, exp := range expected {
		if exp == actual {
			return true
		}
	}
	return false
}

func transitStatus(ctx sdk.Context, value store.Value, from Status, to Status) bool {
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

func ConcatHash(bzz ...[]byte) []byte {
	hasher := tmhash.NewTruncated()
	for _, bz := range bzz {
		hasher.Write(bz)
	}
	return hasher.Sum(nil)
}
