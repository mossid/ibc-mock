package ibc

import (
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
)

// reimplementing tmtypes.MockPV to make it marshallable

type mockPV struct {
	PrivKey crypto.PrivKey
}

var _ tmtypes.PrivValidator = (*mockPV)(nil)

func newMockPV() *mockPV {
	return &mockPV{ed25519.GenPrivKey()}
}

func (pv *mockPV) GetAddress() tmtypes.Address {
	return pv.PrivKey.PubKey().Address()
}

func (pv *mockPV) GetPubKey() crypto.PubKey {
	return pv.PrivKey.PubKey()
}

func (pv *mockPV) SignVote(chainID string, vote *tmtypes.Vote) error {
	signBytes := vote.SignBytes(chainID)
	sig, err := pv.PrivKey.Sign(signBytes)
	if err != nil {
		return err
	}
	vote.Signature = sig
	return nil
}

func (pv *mockPV) SignProposal(string, *tmtypes.Proposal) error {
	panic("not needed")
}

func (pv *mockPV) SignHeartbeat(string, *tmtypes.Heartbeat) error {
	panic("not needed")
}
