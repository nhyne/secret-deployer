package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a string",
	Long: `Encrypt a string to be put into a secret config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		keyVals, err := cmd.Flags().GetStringSlice("key-vals")
		if err != nil {
			fmt.Println("ERROR: invalid key-vals list")
		}

		if len(keyVals) % 2 == 1 {
			fmt.Println("ERROR: odd number of key-vals")
		}


	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("namespace", "n", "", "Namespace for secret config")
	encryptCmd.MarkFlagRequired("namespace")
	encryptCmd.Flags().StringP("secret", "s", "", "Secret name for secret config")
	encryptCmd.MarkFlagRequired("secret")
	encryptCmd.Flags().StringSlice("key-vals", nil, "Key Values. Should follow pattern: 'key1,value1,key2,value2....")
	encryptCmd.MarkFlagRequired("key-vals")
}
