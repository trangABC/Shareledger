package subcommands

import (
	"fmt"

	"github.com/spf13/cobra"

	// "github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	cmn "github.com/tendermint/tendermint/libs/common"

	"github.com/sharering/shareledger/types"
)

var (
	showAll bool
)

// ShowPrivKeyCmd dumps node's Private Key to the standard output.
var ShowPrivKeyCmd = &cobra.Command{
	Use:   "show_priv_key",
	Short: "Show this node's private key",
	RunE:  showPrivKey,
}

func init() {
	ShowPrivKeyCmd.Flags().BoolVar(&showAll, "showAll", false, "Show all private, public key and address of ShareLedger and corresponding Tendermint ones.")
}

func showPrivKey(cmd *cobra.Command, args []string) error {
	privValKeyFile := config.PrivValidatorKeyFile()
	privValStateFile := config.PrivValidatorStateFile()

	var pv *privval.FilePV

	if cmn.FileExists(privValKeyFile) {
		pv = privval.LoadFilePV(privValKeyFile, privValStateFile)
		privateKey := types.ConvertToPrivKey(pv.Key.PrivKey)

		fmt.Printf("%x\n", privateKey[:])

		if showAll {
			publicKey := privateKey.PubKey()

			fmt.Printf("Public Key: %x\n", publicKey[:])
			fmt.Printf("Address   : %X\n", publicKey.Address()[:])

			fmt.Printf("\n***TENDERMINT****\n\n")
			fmt.Printf("Public Key: %x\n", pv.Key.PubKey)
			fmt.Printf("Address   : %X\n", pv.Key.Address[:])

		}
		return nil
	}

	fmt.Printf("Private Validator File not found.")
	return nil
}
