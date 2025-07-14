package pow

import (
    "cosmossdk.io/store"
    "github.com/cosmos/cosmos-sdk/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ethereum/go-ethereum/common"
    "github.com/cometbft/cometbft/libs/log"
)

type Keeper struct {
    storeKey sdk.StoreKey
}

func NewKeeper(storeKey sdk.StoreKey) Keeper {
    return Keeper{storeKey: storeKey}
}

func (k Keeper) GetLastBlockHash(ctx sdk.Context) common.Hash {
    store := ctx.KVStore(k.storeKey)
    hash := store.Get([]byte("last_block_hash"))
    if hash == nil {
        return common.Hash{}
    }
    return common.BytesToHash(hash)
}

func (k Keeper) DistributeReward(ctx sdk.Context, minerAddr sdk.AccAddress, reward sdk.Coin) {
    ctx.Logger().Info("Reward distributed", "miner", minerAddr.String(), "amount", reward.String())
}

func (k Keeper) SaveBlock(ctx sdk.Context, header BlockHeader) {
    store := ctx.KVStore(k.storeKey)
    hash := calculateHash(header)
    store.Set([]byte("last_block_hash"), hash[:])
    ctx.Logger().Info("Block saved", "hash", hash.Hex())
}