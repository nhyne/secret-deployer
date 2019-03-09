package cmd

import (
	"fmt"
	"github.com/nhyne/secret-deployer/pkg/secretConfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a single encrypted secret file.",
	Run: func(cmd *cobra.Command, args []string) {
		projectId := viper.Get("projectId")
		location := viper.Get("location")

		fmt.Println("projectId: %s, location: %s", projectId, location)

	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().Bool("dryrun", true, "When true will show secrets to change.")
	applyCmd.Flags().StringP("file", "f", "", "Secret config to apply.")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// applyCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// applyCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func readSecretConfig(secretConfigPath string) (secretConfig *secretConfig.EncryptedSecretConfig, err error) {
	yamlFile, err := ioutil.ReadFile(secretConfigPath)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yamlFile, secretConfig)
	if err != nil {
		return
	}

	return
}
