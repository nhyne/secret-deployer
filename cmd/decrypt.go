package cmd

import (
	"fmt"
	"github.com/nhyne/secret-deployer/pkg/secretConfig"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a secret config file to plain text.",
	Run: func(cmd *cobra.Command, args []string) {

		secretConfigPath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("error reading file path flag: %v", err)
		}

		outputPath, err := cmd.Flags().GetString("out-file")
		if err != nil {
			fmt.Printf("error reading out file flag: %v", err)
		}

		projectId := viper.Get("projectId")
		location := viper.Get("location")
		keyringId := viper.Get("keyringId")
		keyId := viper.Get("keyId")

		kmsKeyId := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyringId, keyId)
		decryptedSecretConfig, err := secretConfig.DecryptSecretConfig(secretConfigPath, kmsKeyId)
		if err != nil {
			fmt.Errorf("error decrypting secret config: %v", err)
		}

		err = writeDecryptedSecretConfig(decryptedSecretConfig, outputPath)
		if err != nil {
			fmt.Printf("could not write output: %v", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)
	decryptCmd.Flags().StringP("file", "f", "", "Path to secret config file")
	decryptCmd.MarkFlagRequired("file")
	decryptCmd.Flags().StringP("out-file", "o", "", "Target output file. Will delete the file if it exists.")
	decryptCmd.MarkFlagRequired("out-file")
}

func writeDecryptedSecretConfig(plaintextSecretConfig *secretConfig.PlaintetSecretConfig, path string) (error) {
	decryptedYamlBytes, err := yaml.Marshal(plaintextSecretConfig)
	if err != nil {
		return fmt.Errorf("could not marshal yaml: %v", err)
	}

	err = ioutil.WriteFile(path, decryptedYamlBytes, 0644)
	if err != nil {
		return fmt.Errorf("could not write to file: %v", err)
	}

	return nil
}