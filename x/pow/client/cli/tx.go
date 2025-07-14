package cli

import (
    "github.com/cosmos/cosmos-sdk/client"
    "github.com/cosmos/cosmos-sdk/client/flags"
    "github.com/cosmos/cosmos-sdk/client/tx"
    "github.com/kavinda/bether/x/pow"
    "github.com/spf13/cobra"
)

func CmdMineBlock() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "mine-block [miner-address]",
        Short: "Mine a block using PoW",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            clientCtx, err := client.GetClientTxContext(cmd)
            if err != nil {
                return err
            }
            minerAddr, err := sdk.AccAddressFromBech32(args[0])
            if err != nil {
                return err
            }
            module := pow.NewPoWModule(pow.NewKeeper(clientCtx.KVStoreKey))
            header, err := module.MineBlock(clientCtx, minerAddr)
            if err != nil {
                return err
            }
            clientCtx.PrintString("Block mined: " + header.PreviousHash.Hex() + "\n")
            return nil
        },
    }
    flags.AddTxFlagsToCmd(cmd)
    return cmd
}
