package cmd

import (
	"fmt"
	"github.com/nhyne/secret-deployer/pkg/secretConfig"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

// applyCmd represents the apply command
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a single encrypted secret file.",
	RunE: func(cmd *cobra.Command, args []string) (error) {
		configPath, err := cmd.Flags().GetString("file")
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
		err = doesInputFileExist(configPath)
		if err != nil {
			return err
		}

		kmsKeyId, err := getGoogleKMSId()
		if err != nil {
			return err
		}
		err = applySecret(configPath, kmsKeyId)
		if err != nil {
			return fmt.Errorf("coult not apply secret: %v", err)
		}
		return nil
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

	clientset, err := generateClientsetFromFile("/Users/adamjohnson/.kube/config")
	if err != nil {
		return err
	}
	secretExists, err := doesSecretExist(secretObject.Name, secretObject.Namespace, clientset)
	if err != nil {
		return fmt.Errorf("could not check if secret exists: %v", err)
	}
	if secretExists {
		return fmt.Errorf("secret already exists, use update to modify it.")
	}

	_, err = clientset.CoreV1().Secrets(namespace).Create(secretObject)
	if err != nil {
		return fmt.Errorf("could not create secret: %v", err)
	}

	return nil
}

//// returns true if the secret exists
func doesSecretExist(secretName string, namespace string, clientset *kubernetes.Clientset) (bool, error) {
	_, err := clientset.CoreV1().Secrets(namespace).Get(secretName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	return true, nil
}

func generateClientsetFromFile(kubeconfigPath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, fmt.Errorf("coulf not generate config: %v", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("could not create clientset: %v", err)
	}
	return clientset, nil
}

func doesInputFileExist(secretConfigPath string) (error) {
	if _, err := os.Stat(secretConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("secret file \"%s\" does not exist", secretConfigPath)
	} else if err != nil {
		return fmt.Errorf("error stating secret config file: %v", err)
	}
	return nil
}
