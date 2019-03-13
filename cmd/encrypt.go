package cmd

import (
	"fmt"
	"github.com/nhyne/secret-deployer/pkg/secretConfig"
	"gopkg.in/yaml.v2"
	"io/ioutil"

	"github.com/spf13/cobra"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a string",
	Long: `Encrypt a string to be put into a secret config file.`,
	RunE: func(cmd *cobra.Command, args []string) (error) {
		kmsKeyId, err := getGoogleKMSId()
		if err != nil {
			return err
		}
		inputSilce, err := cmd.Flags().GetStringSlice("key-vals")
		if err != nil {
			return fmt.Errorf("error: invalid key-vals list")
		} else if len(inputSilce) % 2 == 1 {
			return fmt.Errorf("error: odd number of key-vals")
		}

		inputKeyFiles, err := cmd.Flags().GetStringSlice("raw-secret-files")
		if err != nil {
			return fmt.Errorf("error reading raw-secret-files flag")
		} else if len(inputKeyFiles) % 2 == 1 {
			return fmt.Errorf("error: off number of raw-secret-files")
		}

		if inputKeyFiles == nil && inputSilce == nil {
			return fmt.Errorf("key-vals and rae-secret-files cannot both be blank")
		}

		outFile, err := cmd.Flags().GetString("out-file")
		if err != nil {
			return fmt.Errorf("error reading output file flag: %v", err)
		}

		namespace, err := cmd.Flags().GetString("namespace")
		if err != nil {
			return fmt.Errorf("error reading namespace flag: %v", err)
		}

		secretName, err := cmd.Flags().GetString("secret")
		if err != nil {
			return fmt.Errorf("error reading secret flag: %v", err)
		}

		plaintextKeyVals := generatePlainTextSlice(inputSilce)
		filePlaintextKeyVals, err := generatePlaintextSliceFromFile(inputKeyFiles)
		if err != nil {
			return fmt.Errorf("could not read input files: %v", err)
		}

		allPlaintextKeyVals := append(plaintextKeyVals, filePlaintextKeyVals...)

		encryptedSecretConfig, err := secretConfig.GenerateSecretConfig(kmsKeyId, namespace, secretName, allPlaintextKeyVals)
		if err != nil {
			return fmt.Errorf("error generating secret config: %v", err)
		}

		fileBytes, err := yaml.Marshal(encryptedSecretConfig)
		if err != nil {
			return fmt.Errorf("error marshalling secret config: %v", err)
		}

		err = ioutil.WriteFile(outFile, fileBytes, 0644)
		if err != nil {
			return fmt.Errorf("error writing secret config to file: %v", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().StringP("namespace", "n", "", "Namespace for secret config")
	encryptCmd.MarkFlagRequired("namespace")
	encryptCmd.Flags().StringP("secret", "s", "", "Secret name for secret config")
	encryptCmd.MarkFlagRequired("secret")
	encryptCmd.Flags().StringP("out-file", "o", "", "Target output file. Will delete the file if it exists.")
	encryptCmd.MarkFlagRequired("out-file")

	encryptCmd.Flags().StringSlice("raw-secret-files", nil, "Path to files with plaintext secret to be encrypted. Should follow pattern: 'key1,path1,key2,path2...'")
	encryptCmd.Flags().StringSlice("key-vals", nil, "Key Values. Should follow pattern: 'key1,value1,key2,value2...'")
}

func generatePlaintextSliceFromFile(inputFileSlice []string) ([]*secretConfig.PlaintextSecretKeyValue, error) {
	plaintextKeyValSlice := make([]*secretConfig.PlaintextSecretKeyValue, 0)
	keys := make([]string, 0)
	files := make([][]byte, 0)

	for i, val := range inputFileSlice {
		if i % 2 == 0 {
			keys = append(keys, val)
		} else {
			bytes, err := ioutil.ReadFile(val)
			if err != nil {
				return nil, fmt.Errorf("could not read file: %v", err)
			}
			files = append(files, bytes)
		}
	}

	for i := range keys {
		plaintextKeyValSlice = append(plaintextKeyValSlice, &secretConfig.PlaintextSecretKeyValue{Key: keys[i], Value: files[i]})
	}
	return plaintextKeyValSlice, nil

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
		plaintextKeyValSlice = append(plaintextKeyValSlice, &secretConfig.PlaintextSecretKeyValue{Key: keys[i], Value: []byte(vals[i])})
	}
	return plaintextKeyValSlice
}
