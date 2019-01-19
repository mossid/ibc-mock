package ibctest

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/mossid/ibc-mock/x/ibc2"
)

func NewHandler(k keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgSpeak:
			return handleMsgSpeak(ctx, k, msg)
		case MsgListen:
			return handleMsgListen(ctx, k, msg)
		case MsgSpeakSafe:
			return handleMsgSpeakSafe(ctx, k, msg)
		case MsgUpdate:
			return handleMsgUpdate(ctx, k, msg)
		case MsgInit:
			return handleMsgInit(ctx, k, msg)
		case MsgPushMessage:
			return handleMsgPushMessage(ctx, k, msg)
		case ibc.MsgPull:
			var data []byte
			var tags sdk.Tags
			for _, packet := range msg.Packets {
				var result sdk.Result
				switch packet := packet.(type) {
				case PacketPushMessage:
					result = handlePacketPushMessage(ctx, k, packet)
				}
				data = append(data, result.Data...)
				tags = append(tags, result.Tags...)
			}
			return sdk.Result{
				Code: sdk.CodeOK,
				Data: data,
				Tags: tags,
				// TODO: log
			}
			/*
				case MsgListenSafe:
					return handleMsgUpdate(ctx, k, msg)
			*/
		default:
			return sdk.ErrUnknownRequest("aaa").Result()
		}
	}
}

func handleMsgSpeak(ctx sdk.Context, k keeper, msg MsgSpeak) sdk.Result {
	if !k.ibc.Speak(ctx, msg.ChainID, nil) {
		return ErrConnSpeakFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handleMsgListen(ctx sdk.Context, k keeper, msg MsgListen) sdk.Result {
	if !k.ibc.Listen(ctx, msg.ChainID, msg.Config, nil) {
		return ErrConnSpeakFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handleMsgSpeakSafe(ctx sdk.Context, k keeper, msg MsgSpeakSafe) sdk.Result {
	fmt.Printf("readyid %x\n", msg.ChainID)

	if !k.ibc.SpeakSafe(ctx, msg.ChainID,
		msg.Height, msg.Proof, msg.RemoteConfig, nil,
	) {
		return ErrConnListenFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

/*
func handleMsgListenSafe(ctx sdk.Context, k ibc.Keeper, msg MsgListenSafe) sdk.Result {
	// TODO?
	return sdk.Result{}
}
*/
func handleMsgUpdate(ctx sdk.Context, k keeper, msg MsgUpdate) sdk.Result {
	source := ibc.NewSource(msg.Commits[0].SignedHeader.Header, msg.Commits)

	fmt.Printf("id %x\n", msg.ChainID)

	if !k.ibc.Update(ctx, msg.ChainID, source, nil) {
		return ErrConnUpdateFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handleMsgInit(ctx sdk.Context, k keeper, msg MsgInit) sdk.Result {
	if !k.ibc.Init(ctx, msg.PortID, msg.Config) {
		return ErrConnUpdateFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handleMsgPushMessage(ctx sdk.Context, k keeper, msg MsgPushMessage) sdk.Result {
	if !k.ibc.Push(ctx, msg.PortID, msg.ChainID, PacketPushMessage{msg.Message}) {
		return ErrConnUpdateFailed(DefaultCodespace).Result()
	}

	return sdk.Result{}
}

func handlePacketPushMessage(ctx sdk.Context, k keeper, packet PacketPushMessage) sdk.Result {
	return sdk.Result{Tags: sdk.NewTags("packet-message", []byte(packet.Message))}
}
