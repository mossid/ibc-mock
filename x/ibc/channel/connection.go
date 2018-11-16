package channel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc/types"
)

type PerConnectionChannel struct {
	core ChannelCore
}

func PerConnectionChannelType() types.ChannelType {
	return types.ChannelType{0x00, 0x00, 0x00, 0x00}
}

func (ch PerConnectionChannel) OpenChannel(ctx sdk.Context, info types.ChannelInfo, chainID []byte) {
	pay := PayloadChannelListening{info}
	ch.core.send(ctx, pay, types.PacketQueueID(), chainID)
}

func (ch PerConnectionChannel) OpenConnection(ctx sdk.Context, config types.ConnectionConfig, chainID []byte) {
	pay := PayloadConnectionListening{config, chainID}
	ch.core.send(ctx, pay, types.PacketQueueID(), chainID)
}

// Using pointer to avoid recursive definition
type PerConnectionTemplate struct {
	*ChannelTemplate
}

func (tem PerConnectionTemplate) Channel(chainID []byte) PerConnectionChannel {
	return PerConnectionChannel{
		core: tem.ChannelCore(chainID),
	}
}
