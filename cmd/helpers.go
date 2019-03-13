package cmd

import (
	"fmt"
	"github.com/spf13/viper"
)

func getGoogleKMSId() (string, error) {
	projectId := viper.Get("projectId")
	if projectId == "" {
		return "", fmt.Errorf("projectId cannot be empty")
	}
	location := viper.Get("location")
	if location == "" {
		return "", fmt.Errorf("location cannot be empty")
	}
	keyringId := viper.Get("keyringId")
	if keyringId == "" {
		return "", fmt.Errorf("keyringId cannot be empty")
	}
	keyId := viper.Get("keyId")
	if keyId == "" {
		return "", fmt.Errorf("keyId cannot be empty")
	}

	kmsKeyId := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", projectId, location, keyringId, keyId)
	return kmsKeyId, nil
}