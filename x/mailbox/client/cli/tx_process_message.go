package cli

import (
	"github.com/KYVENetwork/hyperlane-cosmos/x/mailbox/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdProcessMessage() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "process [metadata] [message]",
		Short: "Process a Hyperlane message",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			metadata := args[0]
			message := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.MsgProcessMessage{
				Metadata: metadata,
				Message:  message,
				Sender:   clientCtx.GetFromAddress().String(),
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
