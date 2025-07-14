package pow

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "math/big"
    "time"
    "github.com/cosmos/cosmos-sdk/types"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/ethereum/go-ethereum/common"
)

type Miner struct {
    Address    sdk.AccAddress
    Nonce      uint64
    Difficulty *big.Int
}

type BlockHeader struct {
    PreviousHash common.Hash
    Timestamp    int64
    Nonce        uint64
    Difficulty   *big.Int
    Reward       sdk.Coin
}

type Keeper interface {
    GetLastBlockHash(ctx sdk.Context) common.Hash
    DistributeReward(ctx sdk.Context, miner sdk.AccAddress, reward sdk.Coin)
    SaveBlock(ctx sdk.Context, header BlockHeader)
}

type PoWModule struct {
    keeper Keeper
}

func NewPoWModule(keeper Keeper) PoWModule {
    return PoWModule{keeper: keeper}
}

func (m PoWModule) MineBlock(ctx sdk.Context, minerAddr sdk.AccAddress) (BlockHeader, error) {
    prevHash := m.keeper.GetLastBlockHash(ctx)
    difficulty := big.NewInt(1e6)
    header := BlockHeader{
        PreviousHash: prevHash,
        Timestamp:    time.Now().Unix(),
        Difficulty:   difficulty,
        Reward:       sdk.NewCoin("tbtr", sdk.NewInt(1_000_000)),
    }
    target := new(big.Int).Div(big.NewInt(1).Lsh(big.NewInt(1), 256), difficulty)
    for nonce := uint64(0); nonce < 1e9; nonce++ {
        header.Nonce = nonce
        hash := calculateHash(header)
        hashInt := new(big.Int)
        hashInt.SetString(hex.EncodeToString(hash[:]), 16)
        if hashInt.Cmp(target) < 0 {
            m.keeper.SaveBlock(ctx, header)
            m.keeper.DistributeReward(ctx, minerAddr, header.Reward)
            return header, nil
        }
        if time.Now().Unix()-header.Timestamp > 10 {
            return BlockHeader{}, fmt.Errorf("mining timeout")
        }
    }
    return BlockHeader{}, fmt.Errorf("no valid nonce found")
}

func calculateHash(header BlockHeader) common.Hash {
    data := fmt.Sprintf("%x%d%d", header.PreviousHash, header.Timestamp, header.Nonce)
    hash := sha256.Sum256([]byte(data))
    return common.BytesToHash(hash[:])
}

func (m PoWModule) ValidateBlock(ctx sdk.Context, header BlockHeader) bool {
    hash := calculateHash(header)
    target := new(big.Int).Div(big.NewInt(1).Lsh(big.NewInt(1), 256), header.Difficulty)
    hashInt := new(big.Int)
    hashInt.SetString(hex.EncodeToString(hash[:]), 16)
    return hashInt.Cmp(target) < 0
}