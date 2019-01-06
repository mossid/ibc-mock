# IBC Implementation Spec ver2

// DRAFT

IBC implementation on Cosmos-SDK specification.

## State

The keys have hierarchical structure, ChainID, QueueID, Index, from the top. Let's say the ChainID is `0xABCD`, the QueueID is `0xCAFEBABE`.

The key of ChainID(`0xABCD`) stores the information about that chain. We will call it *Chain Config State* for this Chain.

The key of ChainID|QueueID(`0xABCDCAFEBABE`) stores the information about that queue on the connection with that chain. We will call it *Queue Config State* for this Queue.

Under this subspace, actual interchain messages are stored. The messages are stored under integer index. The index is increased each time a new message is pushed. For example, the first message in this queue will be stored at `0xABCDCAFEBABE0000000000000000`(big endian encoding here)

Examples: 
```
0xR002: root key for the ibc module, stores local IBC configuration
0xR002|ABCD: config key for chain ABCD, stores verified remote IBC configuration
0xR002|ABCD|CAFEBABE: config key for queue CAFEBABE on the connection with chain ABCD, stores queue configuration
0xR002|ABCD|CAFEBABE|0000000000000000: first message in the queue
``` 

The byte length of the chainID and queueID is defined in the local IBC configuration. The other chain can use this information to construct appropriate merkle keypath.

## Connection

Connection(often abbreviated as `conn`) is a mutual liteclient tracking between two chains. To open a conn, it is required to do a handshake, making it sure that the both chain are tracking each other and referring the right position for the queues. 

### Handshaking

To open a conn, an account should send a `MsgOpenConn` to the chain (A). The `MsgOpenConn` contains

* Root-of-trust fullcommit from (B)
* Valid IBC configuration, including perconn queue location, version, etc
* Local chain id for (B)

The chain stores them in the chain config state (S). The fullcommit stored means that the chain is ready to verify future commit from the other chain (B), no matter which chain it is. 

Once `MsgOpenConn`s are executed on the both chain with each other's fullcommit, now they can verify whether the other chain is ready to receive packets from itself. It is done when an account send `MsgReadyConn` to the chains. `MsgReadyConn` contains the value of and the merkle path to (S), which is be used for verification.

Since (B) should be able to verify that (A) has its fullcommit, (B) should store the headers of itself. Storing all headers is definitly inefficient, so we will use `MsgCheckpoint`, which accounts can send to manually store the current header as a valid root-of-trust. (B) will only permit the commits stored in (A) that is checkpointed.

Perconn queue is instantiated while `MsgReadyConn` is executed, using the config. 

`MsgEstablishConn` is the final step to make a connection. A dummy payload is sent over the perconn queue, indicating that the chain has established the connection. After receiving it from the other chain, other tasks, such as opening queue and updating connection can be done.

### Following block headers

Updating a connection can be done by sending `MsgUpdateConn`. 

### Updating a connection

### Closing a connection

## Queue

https://github.com/cosmos/cosmos-sdk/blob/develop/docs/spec/ibc/channels-and-packets.md#323-queue

A queue is a subspace in the state, storing ordered messages. Messages are stored under a `PrefixStore` with bigendian encoded index key. The queue must be opened before any message are sent over it.

A queue can be opened either with handshaking or without. With handshaking, it is ensured that the other chain has a compatible format for the queue and is ready to receive the packet. Without handshaking, it is not ensured, but still able to check timeout and recover from the origin chain.

// TODO: check go-amino varint encoding is lexicographical order

### Perconn queue

Perconn queue exists for each connection. It works as meta queue, sending informations from IBC module itself. When a connection is established between two chains, perconn queue is automatically instantiated, providing reliable message sending between them.

Perconn queue works same with the other queues, except it is not involved with normal queue handshaking process, which requires perconn queue in prior. 

### Handshaking

To open a queue, the chain have to initialize the queue config state. It is usually not done by the accounts, but by the internal logic of modules. 

After the queue config state has been initialized, both chains  

 the chain send `PayloadOpenQueue` to the other chain via the perconn queue. 

To open a queue, the chain should send `PayloadOpenQueue` to the other chain via the perconn queue. It contains `QueueConfig`, which describing the location of the queue in the state, encoding scheme, queue version, etc. 

When a chain receives `PayloadOpenQueue` from another chain, it first checks 

### Message Sending


