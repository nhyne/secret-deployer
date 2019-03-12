package cmd

import (
	"fmt"
	"github.com/nhyne/secret-deployer/pkg/secretConfig"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a string",
	Long: `Encrypt a string to be put into a secret config file.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectId := viper.Get("projectId")
		location := viper.Get("location")
		keyringId := viper.Get("keyringId")
		keyId := viper.Get("keyId")

		kmsKeyId := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyringId, keyId)
		inputSilce, err := cmd.Flags().GetStringSlice("key-vals")
		if err != nil {
			fmt.Println("error: invalid key-vals list")
		}

		if len(inputSilce) % 2 == 1 {
			fmt.Println("error: odd number of key-vals")
		}

		outFile, err := cmd.Flags().GetString("out-file")
		if err != nil {
			fmt.Printf("error reading output file flag: %v", err)
		}

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			fmt.Printf("error reading namespace flag: %v", err)
		}

		secretName, err := cmd.Flags().GetString("secret")
		if err != nil {
			fmt.Printf("error reading secret flag: %v", err)
		}

		plaintextKeyVals := generatePlainTextSlice(inputSilce)

		encryptedSecretConfig, err := secretConfig.GenerateSecretConfig(kmsKeyId, namespace, secretName, plaintextKeyVals)
		if err != nil {
			fmt.Printf("error generating secret config: %v", err)
		}

		fileBytes, err := yaml.Marshal(encryptedSecretConfig)
		if err != nil {
			fmt.Printf("error marshalling secret config: %v", err)
		}

		err = ioutil.WriteFile(outFile, fileBytes, 0644)
		if err != nil {
			fmt.Printf("error writing secret config to file: %v", err)
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
	encryptCmd.Flags().StringP("out-file", "o", "", "Target output file. Will delete the file if it exists.")
	encryptCmd.MarkFlagRequired("out-file")
}

func generatePlainTextSlice(inputSlice []string) ([]*secretConfig.PlaintextSecretKeyValue) {
	plaintextKeyValSlice := make([]*secretConfig.PlaintextSecretKeyValue, 0)
	keys := make([]string, 0)
	vals := make([]string, 0)
	for i, val := range inputSlice {
		if i % 2 == 0 {
			keys = append(keys, val)
		} else {
			vals = append(vals, val)
		}
	}

	for i, _ := range keys {
		plaintextKeyValSlice = append(plaintextKeyValSlice, &secretConfig.PlaintextSecretKeyValue{Key: keys[i], Value: vals[i]})
	}
	return plaintextKeyValSlice
}