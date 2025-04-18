package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
)

var (
	newOwner          string
	renounceOwnership bool
)

func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        "hooks",
		Short:                      "Hyperlane Core-Hooks commands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		NewIgpCmd(),
		NewMerkleCmd(),
		NewNoopHookCmd(),
	)

	return txCmd
}
