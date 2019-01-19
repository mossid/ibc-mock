package ibc

import (
	"errors"
	"fmt"

	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/lite"
	liteErr "github.com/tendermint/tendermint/lite/errors"
	tmtypes "github.com/tendermint/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Temporal object, should not persist
type provider struct {
	conn
	ctx     sdk.Context
	chainID []byte
}

var _ lite.Provider = provider{}

func (p provider) LatestFullCommit(chainID string, minHeight, maxHeight int64) (res lite.FullCommit, err error) {
	if chainID != string(p.chainID) {
		err = errors.New("invalid chainID")
		return
	}

	// minHeight and maxHeight required to be inclusive but Range takes endpoint
	// as exclusive
	_, ok := p.commits.Range(uint64(minHeight), uint64(maxHeight)+1).Last(p.ctx, &res)
	if !ok {
		err = liteErr.ErrCommitNotFound()
	}
	return
}

func (p provider) ValidatorSet(chainID string, height int64) (valset *tmtypes.ValidatorSet, err error) {
	if chainID != string(p.chainID) {
		err = errors.New("invalid chainID")
		return
	}

	var commit lite.FullCommit
	err = p.commits.GetSafe(p.ctx, uint64(height), &commit)
	if err != nil {
		return
	}
	valset = commit.Validators
	return
}

func (p provider) SetLogger(logger log.Logger) {
	// TODO
}

var _ lite.PersistentProvider = provider{}

func (p provider) SaveFullCommit(fc lite.FullCommit) (err error) {
	fmt.Println("save fullcommit", fc)
	if fc.ChainID() != string(p.chainID) {
		err = errors.New("invalid ChainID")
		return
	}
	p.commits.Set(p.ctx, uint64(fc.Height()), fc)
	return
}

type Source struct {
	chainID string
	commits []lite.FullCommit
}

func NewSource(header *tmtypes.Header, commits []lite.FullCommit) Source {
	return Source{
		chainID: header.ChainID,
		commits: commits,
	}
}

var _ lite.Provider = Source{}

func (s Source) LatestFullCommit(chainID string, minHeight, maxHeight int64) (res lite.FullCommit, err error) {
	if s.chainID != chainID {
		err = errors.New("invalid ChainID")
		return
	}

	for i := len(s.commits) - 1; i >= 0; i-- {
		res = s.commits[i]
		height := res.Height()
		if height < minHeight {
			err = errors.New("commit not found")
			return
		}
		if height <= maxHeight {
			return
		}
	}
	err = errors.New("commit not found")
	return
}

func (s Source) SetLogger(logger log.Logger) {
	// TODO
}

func (s Source) ValidatorSet(chainID string, height int64) (valset *tmtypes.ValidatorSet, err error) {
	if chainID != s.chainID {
		err = errors.New("invalid chainID")
		return
	}

	commit, err := s.LatestFullCommit(chainID, height, height)
	if err != nil {
		err = liteErr.ErrUnknownValidators(chainID, height)
		return
	}

	valset = commit.Validators
	return
}

func (s Source) GetLastSignedHeader() tmtypes.SignedHeader {
	return s.commits[len(s.commits)-1].SignedHeader
}

func (s Source) saveAll(p provider) {
	for _, commit := range s.commits {
		p.SaveFullCommit(commit)
	}
}
