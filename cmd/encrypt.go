// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

	// Here you will define your flags and configuration settings.
	encryptCmd.Flags().StringP("namespace", "n", "", "Namespace for secret config")
	encryptCmd.Flags().StringP("secret", "s", "", "Secret name for secret config")
	encryptCmd.Flags().StringSlice("key-vals", "", "Key Values. Should follow pattern: 'key1,value1,key2,value2....")

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
