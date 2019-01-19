package ibctest

import (
	"github.com/mossid/ibc-mock/x/ibc2"
)

type keeper struct {
	ibc ibc.Channel
}

func newKeeper(ch ibc.Channel) keeper {
	return keeper{ch}
}
