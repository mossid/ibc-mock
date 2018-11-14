# IBC Bank module implementation spec

`ibcbank` module is a submodule of `bank` that manages interchain asset transfer. It automatically manages the global legder of chains, preventing corruption of the whole network when a chain got byzantine failure. 

## Channel

`ibcbank` uses `BiChannel` to handle failure of asset transfer. 

```go
type PayloadPacketCoins {
    SrcAddr sdk.AccAddress
    DestAddr sdk.AccAddress
    Coins sdk.Coins
}

type PayloadReceiptCoins {
    
}
```
