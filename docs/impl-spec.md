# IBC Implementation Spec

// DRAFT

IBC on Cosmos-SDK implementation specification. The protocol specification is here: https://github.com/cosmos/cosmos-sdk/tree/develop/docs/spec/ibc



## Connections

`IBC` module manage IBC connections. `ibc.Keeper` stores fullcommits from the chains those have a connection established. The connections are used for verifying whether a packet exists in the egress queue on the source chain. It is done by lightclient proof which guarantees the packet is generated by the valid excution of the application and signed by the supermajority of the valset from the source chain.

### Definitions


```go
func CommitHeightKey(srcChain string) []byte {
    return append([]byte{0x00}, []byte(srcChain)...)
}
func (k Keeper) commitHeight(ctx sdk.Context, srcChain string) store.Value {
    return store.NewValue(k.cdc, ctx.KVStore(k.key), CommitHeightKey(srcChain))
}
```
`commitHeight` returns a `store.Value` indicating the latest stored commit's height. When a new commit has been reported, its height must be greater than the `commitHeight.`

```go
func CommitListPrefix(srcChain string) []byte {
    return append([]byte{0x01}, amino.MarshalBinaryLengthPrefixed(srcChain)...)   
}

func (k Keeper) commitList(ctx sdk.Context, srcChain string) store.List {
    return store.NewList(k.cdc, store.NewPrefixStore(ctx.KVStore(k.key), CommitListPrefix(srcChain)))
}
```
`commitList` returns a `store.List` for the stored commits. This list works as `map[int64]lite.FullCommit`. If a reported commit is valid, it is stored in the list and referred when there is need to prove a packet from that height.

// TODO: figure out how to deal with two chains running on different IBC version

### Connection Lifecycle

#### Opening a connection

A connection can be opened permissionlessly by sending `MsgOpenConnection`. 

```go
type MsgOpenConnection struct {
    ROT/*SYN*/ lite.FullCommit
    Signer sdk.AccAddress
}
```

After the `MsgOpenConnection` is executed, the chain is ready to receive the messages.

`ROT` stands for "Root of Trust", where the sequence of valset update begins. It is recommended to use the genesis block of the chain, but it is possible to use different block. To achieve this, IBC module tracks the valset of its chain, and checkpoints the header when there is `+1/3` valset change. These checkpoints, including the genesis, works as valid `ROT` that other chains who wants to receive messages from the chain can use.

`Signer` is the address sent this message, paying fee for it.

`handleMsgOpenConnection` stores the `ROT` into the `commitList` if the `ROT` is valid. 

// TODO: figure out how to determine ChainID(hashing ROT with userprovided identifier?)

To assure that the both of the chain has opened a connection to each other, the protocol does 4-way handshake. `MsgOpenConnection` sends a `PayloadConnectionListening`, where its `ROT` acts as a `ACK`. Any other message sending/receiving is blocked until the chain receives `PayloadConnectionListening` from the other chain.

```go
type PayloadConnectionListening struct {
    ROT/*ACK*/ lite.FullCommit
    ChainID string // ChainID of the chain receiving this message registered on the sender
}
```

A channel is opened by IBC module after the `MsgOpenConnection` is executed so the chain can verify the incoming `PayloadConnection` from the other. When the `PayloadConnectionListening` is executed, IBC module verifies its `ROT` assuring that is a valid `ROT` checkpointed by this chain, and allowes the full message sending/receiving to the other channels.

The `ChainID`, indicating the ID of this chan registered on the other chain(ID<sub>this</sub>@<sub>that</sub>), is associated with the ID of the other chain registered on this chain(ID<sub>that</sub>@<sub>this</sub>).

#### Following block headers

Updating a connection by submitting new commit can be done by sending `MsgUpdateConnection`.

```go
type MsgUpdateConnection struct {
    SrcChain string
    Commit lite.FullCommit
    Signer sdk.AccAddress
}
```

`SrcChain` is the ChainId of the other chain for indexing `commitList`. The `Commit`'s valset is required to have less than `1/3` of valset change from the `commitHeight`. The transaction will fail if not. Otherwise the `Commit` is verified by `lite.DynamicVerifier` then stored in `commitList`. Lightclient verification will then be able to use the new commit when is needed.

#### Updating a connection

When there is nonbackwardcompatible update on IBC module on a chain, 

#### Closing a connection

Unlike opening a connection, closing it cannot be done permissionlessly. Connection only can be closed when a byzantine behaviour is detected on the other chain or there is a governance consensus.

// TODO

## Channels & Packets

Each outgoing IBC message is stored in a queue. The messages in a queue is guaranteed to be processed in sequence. We call a set of queue "channel". Channel is isolated with each other. One chain-chain pair can have only one connection(by definition), but can have multiple channel. Also channel can be deployed at runtime, e.g. each contracts in Ethermint can have their own channel. /*The queues that forms a channel cannot be modified after the channel has been created.*/

A module takes either `Channel` or `MultiChannel`, depending on its usage. `MultiChannel` is used for modules who need more than one channel, and can use `MultiChannel.Channel` for instanciate channels. If a module need to define its own queue format, they can wrap `ChannelCore` and `MultiChannelCore` manually.

### Channel Lifecycle

#### Opening a channel

### Definitions

```go
type ChannelInfo interface {
    Verify(ctx sdk.Context, queueid []byte) error
}

func (info ChannelInfo) Verify(queueid []byte/*and possibly other args...*/) error {
    if _, ok := info.queueset[string(queueid)]; !ok {
        return errors.New("Invalid queueid")
    }
    return nil
}

type Channel interface {
    ChannelType() [CHANNELTYPELEN]byte
    SetChannelCore(ChannelCore)
}

type ChannelCore struct {
    k Keeper
    key *ChannelKey
    id []byte
    info ChannelInfo
}

func (k Keeper) Channel(ch Channel, key *ChannelKey, info ChannelInfo) {
    core := ChannelCore {
        k: k,
        key: key,
        id: nil,
        info: info,
    }
    ch.SetChannelCore(core)
}

type MultiChannel struct {
    k Keeper
    key *ChannelKey
    info ChannelInfo
}

func (k Keeper) MultiChannel(key *ChannelKey, info ChannelInfo) MultiChannel {
    return MultiChannel{ 
        k: k,
        key: key,
        info: info,
    }
}

func (mcc MultiChannel) Channel(ch Channel, id []byte) {
    if len(id) == 0 {
        panic("Cannot provide nil id for MultiChannel")
    }

    core := ChannelCore {
        k: mcc.k,
        key: mcc.key,
        info: mcc.info,
        id: id,
    }
    
    ch.SetChannelCore(core)
}

```

First there are two types, `ChannelCore` and `MultiChannelCore`. `ChannelCore` is used for accesing the raw channel, exposing all queues declared by `ChannelInfo`. `MultiChannelCore` is a conceptualized multiple channels, implemented as a `ChannelCore` where `key *ChannelKey, info ChannelInfo` is fixed. 

```go
const QUEUEIDLEN = 2

func QueueID(channelid []byte, queueid [QUEUEIDLEN]byte) (res []byte) {
    // Prefixing length to ids for preventing collision
    if len(channelid) == 0 {
        res = append(res, amino.MarshalBinaryLengthPrefixed(channelid)...)
    }
    res = append(res, queueid[:]...)
    return
}

// channelid []byte -> queueid []byte -> queue store.Queue
func QueuePrefix(channelid []byte, queueid []byte) []byte {
    return append([]byte{0x00}, QueueID(channelid, queueid)...)
}

// channelid []byte -> queueid []byte -> incoming uint64
func IncomingKey(channelid []byte, queueid []byte) []byte {
    return append([]byte{0x01}, QueueID(channelid, queueid)...)
}

// 1. Open the channel, sending channel information
func (ch ChannelCore) Open(ctx sdk.Context, chainid string) bool {
    // Send ChannelType through the PerConnection Channel
}

// 2. The payload is sent through IBC channel
// It stores the decoded information to the state
type PayloadChannelOpen struct {
    info []byte
}

// 3. Any usage of the channel checks whether the Payload has been arrived
// If the Queueset informations are asymetrical, 
func (ch ChannelCore) Established(chainid string) bool {
    // Check the IBC module state whether the PayloadChannelType has arrived from the other chain
}

// 4. After the queue is established, ChannelCore.Queue can be called to obtain message queue
func (ch ChannelCore) Queue(ctx sdk.Context, chainid string, queueid []byte) (q store.Queue, err error) {    
    err = ch.info.Verify(queueid)
    if err != nil {
        return
    }
    
    if !ch.Established(chainid) {
        err = errors.New("Channel not established")
        return
    }

    key := QueueKey(ch.id, queueid)
    store := prefix.NewStore(ctx.KVStore(ch.k.store(ch.key)), )
    return store.NewQueue(ch.k.cdc, ctx.KVStore(ch.k.))
}

// 

```

Opening a channel works similar with opening a connection, except it is 2-way handshake, since we assume(and check) that there is already a connection between two chains. Each channel have their own channel format, which has to be symmetrical between a pair of channels. To make it sure that two channel has a same `ChannelInfo`, it is included in `PayloadChannelOpen`. The queue allows

`ChannelCore` and `StdChannelCore` are not intended to be used as it is, but rather in wrapped form. Since they exposes raw access to the queues, it can be potentially harmful.

```go
type DatagramType byte

const (
    Packet DatagramType = iota
    Receipt
)

func (dt DatagramType) Bytes() []byte {
    return []byte{dt}
}

// UniChannel
type UniChannel struct {
    core ChannelCore
}

func (ch *UniChannel) ChannelType() []byte {
    return []byte{0xBA, 0xAD, 0xBE, 0xEF}
}

func (ch *UniChannel) SetChannelCore(core ChannelCore) {
    ch.core = core
}

func (ch UniChannel) queue(ctx sdk.Context, chainid string) store.Queue {
    return ch.core.Queue(ctx, chainid, Packet.Bytes())
}

// BiChannel contains Packet queue and Receipt queue,
// enabling the chains can fail on a payload
type BiChannel struct {
    core ChannelCore
}

func (ch *BiChannel) ChannelType() []byte {
    return []byte{0x00, 0x01, 0x02, 0x03}
}

func (ch *BiChannel) SetChannelCore(core ChannelCore) {
    ch.core = core
}

func (ch BiChannel) queue(ctx sdk.Context, ty DatagramType, chainid string) store.Queue {
    return ch.core.Queue(ctx, chainid, ty.Bytes())
}

```

`{Uni, Bi}Channel` wraps `ChannelCore` and make the users to access on packet and receipt queue in rescricted way. ~~Initially, there are two queues, one for packets, and one for receipts. Outgoing Packet queue, Incoming Packet queue, Outgoing Receipt queue, Incoming Receipt Queue. `DatagramType` can be extended in a future version of the protocol.~~

### Sending a packet

On `{Uni, Bi}Channel`, user can send packets only through `{}Channel.Send()`.

```go
type Packet struct {
    Header
    Payload
}

type Header struct {
    Route []string
}

type Keeper struct {
    cdc *codec.Codec
    key *sdk.MultiStoreKey
}

func (ch UniChannel) Send(ctx sdk.Context, destChain string, payload Payload) (result sdk.Result) {
    
}

func (ch Channel) Send(ctx sdk.Context, destChain string, channelid []byte, payload Payload) {...}
```

`Send` takes a payload and its destination chain. The `Header` is generated internally using the argument `destChain` and attatched to the payload, forming a full `Packet`. It is then inserted into the egress queue. `MultiChannel.Send` additionally takes `channelid []byte` argument for identifying the exact chennel.

Payload is an interface for interblockchain message.

```go
type Payload interface {
    Route() string
    Type() string
    ValidateBasic() sdk.Error
    GetSignBytes() []byte
    
    QueueID() [QUEUEIDLEN]byte
}
```

For the sake of simplicity, `Payload` contains most of the methods in `sdk.Msg`, excepts for `GetSigners()`. Payload When the `Payload` is wrapped by IBC msg types, they inherits the payload and add the missing method. `QueueID()` is used for determining which queue that the payload will be stored. 

### Relaying a packet

A packet can be relayed through mutiple zones without requiring those zones to acknowledge about the packet. For example, there could be two zones sharing an IBC user module, which the hub doesn't know about, but still can be relayed through.

```go
type MsgRelay struct {
    Packets []Packet
    Proof Proof
    Relayer sdk.AccAddress
}
```

The `MsgRelay` is same with `MsgReceive`(see below). However, the IBC user modules can handle the `MsgRelay` when needed. For example, IBC bank module can recalcuate the balances of the chains to preserve the interchain sanity.

### Receiving a packet

In perspective of the receiving chain, packet is same with a transaction that a chain, not an account, is sending.

```go
type Proof struct {
    merkle.Proof
    Height int64
    SrcChain string
    ChannelID []byte
}

type MsgReceive struct {
    Packets []Packet
    Proof Proof
    Relayer sdk.AccAddress
}

func NewAnteHandler(k Keeper) sdk.AnteHandler {
    return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, result Result, abort bool) {
        ...        
    }
}
```

The `ibc.AnteHandler` verifies the blockchain's light client proof for each `MsgReceive` in a transaction, and revert if invalid.

`MsgReceive` contains three fields.`Packets []Packet` is for the IBC packets sent from a single queue of the chain. `Proof Proof` should be a range proof for the messages. `Relayer sdk.AccAddress` is the transaction relayer who pays the fee.

### Handling a receipt

// TODO

While user modules can just call `Channel.Send()` with receipt payloads to send receipts, it can cause some unsafe situation, so IBC module additionally offer a function to handle receipts.

```go
// package ibc

type Handler func(sdk.Context, Payload) (Payload, sdk.Result)

func Handle(ctx sdk.Context, f Handler, payload Payload) sdk.Result {...}
```

When a module wants to handle receipts safely, they can wrap the handlers with `Handle`, which automatically push the result `Payload` of the callback function to the receipt queue. If the provided 

### Packet relayer

// TODO

#### Trustless relaying

The zones can send assets to other zones living under the other hub. For example,

```
    H
  /   \
A      H'
     /   \
   B      C
```

Let's say the `H` works as the trust source for the whole network, and `H'` works as regional hub. `A` trusts `H` but not `H'`, so it wants the packet relayed without trusting `H'`. In other words, `A` is worried about the possibility that `H'` steals its native token so the value goes down.

If `A` trusts `H'`, `H` only holds the total balance of `H'` in denomination of `Atoken`. No matter how much transaction happened between `B` and `C`, `H` does not care and fully trusts `H'`'s packet. In this case, it is possible that the validators steals the assets from `B` and `C`.

In our case, `A` sends a bank IBC packet to `B`, with an option in its preroute to not trust `H'`. When `H` detects the option, it updates the ledger on `B`, not `H'`. When an user want to withdraw `Atoken` from `B`, the packet pass the `H'` chain, but when it arrives on `H`, it does not only verify the lightclient proof of the packet in `H'`, but also in `B`, using the header of `B` stored in `H'`, since it is guaranteed to exists on `H'`. 

## Optimizations



### Timeout

When timeout happens, the chain who received the packet at last automatically revert the packet flow, recursivly sending the timeout message to the src chains. When timeout happens during the reverting process, packet stops  

### Cleanup

## Global Messaging

A chain need to broadcast global informations, such as supported format types, IBC protocol version, etc. These informations are stored in fixed position in the KVStore so the other chain can refer it without sending it on each connection manually. 

// TODO: figure out which global configs are needed

```go
func
```

## Additional Implementations

### Fee precollection

Each channel can set the minimum fee that each IBC packet has to pay. This fee will collected into the `Packet`, and paid to the relayer.

### Fee autoadjustment

The channel can tract the usage of itself, automatically adjust the minimum fee it requires.

### Receipt handling callback

When calling `Channel.Send`, the caller can pass additional argument, `callback []byte`, to specify which action will be done if the fail receipt is returned. For example, an Ethermint contract sending an IBC packet can register one of its function as error handler.

## Helper Types

```go
package store

type Value struct {
    cdc *codec.Codec
    store sdk.KVStore
    key []byte
}

func (v Value) Get(ptr interface{}) {
    v.cdc.MustUnmarshalBinary(store.Get(v.key), ptr)
}

func (v Value) Set(o interface{}) {
    v.store.Set(v.key, v.cdc.MustMarshalBinary(o))
}

type Queue struct {
    cdc *codec.Codec
    store sdk.KVStore
}
//...
```

## Trust model

IBC does not rely on the assumption that the connected chains will be byzantine safe. The modules those using IBC should care about the situation.

## Binary format

Packet has following binary format

```
----------------------------------------
----------------------------------------
----------------------------------------
```

`Encoding` defines how does the actual payload is encoded. Initially `Amino`, `JSON`, `EVMBinary` format will be supported. Each connection will 

### Amino Encoding

Amino Encoding format is `amino.Codec.MarshalBinary`

### JSON Encoding

Json Encoding format is `amino.Codec.MarshalJSON`

### EVMBinary Encoding

EVMBinary format has fixed length for each payload, making it easier to verify it on 