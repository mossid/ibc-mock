package channel

import (
	"github.com/mossid/ibc-mock/x/ibc/types"
)

type StdChannel struct {
	core ChannelCore
}

func (ch StdChannel) ChannelType() types.ChannelType {
	return types.ChannelType{0xDE, 0xAD, 0xBE, 0xEF}
}

var _ types.Channel = StdChannel{}
