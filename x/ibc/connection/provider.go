package connection

import (
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Temporal object, should not persist
type provider struct {
	remote
	ctx sdk.Context
}

var _ lite.Provider = provider{}

func (p provider) LatestFullCommit(chainID string, minHeight, maxHeight int64) (res lite.FullCommit, err error) {
	// TODO: verify chainID length, heights isPositive

	_, ok := p.commits.Prefix([]byte(chainID)).Range(uint64(minHeight), uint64(maxHeight)).Last(p.ctx, &res)
	if !ok {
		err = fmt.Errorf("aaa")
	}

	return
}

func (p provider) ValidatorSet(chainID string, height int64) (valset *tmtypes.ValidatorSet, err error) {
	err = p.valsets.Prefix([]byte(chainID)).GetSafe(p.ctx, uint64(height), &valset)
	return
}

func (p provider) SetLogger(logger log.Logger) {
	// TODO
}

var _ lite.PersistentProvider = provider{}

func (p provider) SaveFullCommit(fc lite.FullCommit) error {
	value := p.commits.Prefix([]byte(fc.ChainID())).Value(uint64(fc.Height() + 1))
	if value.Exists(p.ctx) {
		return fmt.Errorf("bbb")
	}
	value.Set(p.ctx, fc)
	return nil
}

type source struct {
	chainID      string
	nextValset   *tmtypes.ValidatorSet
	valsetHeight int64
	commits      []lite.FullCommit
}

var _ lite.Provider = source{}

func (p source) LatestFullCommit(chainID string, minHeight, maxHeight int64) (res lite.FullCommit, err error) {
	if chainID != p.chainID {
		err = fmt.Errorf("ccc")
		return
	}

	// Assume that the commits are already sorted
	// TODO: binary search
	for i := len(p.commits) - 1; i >= 0; i-- {
		height := p.commits[i].Height()
		if height <= maxHeight && !(height < minHeight) {
			return p.commits[i], nil
		}
	}

	err = fmt.Errorf("ddd")
	return
}

func (p source) ValidatorSet(chainID string, height int64) (valset *tmtypes.ValidatorSet, err error) {
	// Currently lite/dynamic_verifier.go only requires the valset for fc.Height()+1
	// So we store the valset only for that height
	if chainID != p.chainID {
		err = fmt.Errorf("eee")
		return
	}

	if height != p.valsetHeight {
		err = fmt.Errorf("ggg")
		return
	}

	return p.nextValset, nil
}

func (p source) SetLogger(logger log.Logger) {
	// TODO
}
