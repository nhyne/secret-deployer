package cmd

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/nhyne/secret-deployer/pkg/secretConfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a single encrypted secret file.",
	Run: func(cmd *cobra.Command, args []string) {
		projectId := viper.Get("projectId")
		location := viper.Get("location")
		keyringId := viper.Get("keyringId")
		keyId := viper.Get("keyId")

		configPath, err := cmd.Flags().GetString("file")
		if err != nil {
			fmt.Printf("error: %v", err)
		}

		kmsKeyId := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyringId, keyId)
		applySecret(configPath, kmsKeyId)
	},
}

func init() {
	rootCmd.AddCommand(applyCmd)

	applyCmd.Flags().Bool("dryrun", true, "When true will show secrets to change.")
	applyCmd.Flags().StringP("file", "f", "", "Secret config to apply.")
	applyCmd.MarkFlagRequired("file")
}

func applySecret(path string, kmsKey string) (error) {

	secretObject, namespace, err := secretConfig.ConvertSecretConfigToSecretObject(path, kmsKey)
	if err != nil {
		return err
	}

	config, err := generateConfigFromFile("/Users/adamjohnson/.kube/config")
	if err != nil {
		fmt.Printf("%v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	secretOut, err := clientset.CoreV1().Secrets(namespace).Create(secretObject)
	if err != nil {
		fmt.Printf("could not create secret: %v", err)
	}

	return nil
}

//func applySecretObject(secretObject *corev1.Secret) (error) {
//
//}
//
//// returns true if the secret exists
//func doesSecretExist(secretName string, namespace string) (bool, error) {
//
//}

func generateConfigFromFile(kubeconfigPath string) (*restclient.Config, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("coulf not generate config: %v", err)
	}
	return config, nil
}
