package secretConfig

import (
	cloudkms "cloud.google.com/go/kms/apiv1"
	"context"
	"encoding/base64"
	"fmt"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
)

type PlaintextSecretKeyValue struct {
	Key string
	Value string
}

func GenerateSecretConfig(kmsKeyName string, namespace string, secretName string, keyVals []*PlaintextSecretKeyValue) (encryptedConfig EncryptedSecretConfig, err error) {

	if namespace == "" {
		err = fmt.Errorf("namespace cannot be empty")
		return
	}
	encryptedConfig.Namespace = namespace

	if secretName == "" {
		err = fmt.Errorf("secretName cannot be empty")
		return
	}
	encryptedConfig.SecretName = secretName

	for _, plaintextKeyVal := range keyVals {
		encryptedKeyVal, err := plaintextKeyVal.encryptPlaintextKeyValue(kmsKeyName)
		if err != nil {
			return encryptedConfig, fmt.Errorf("could not encrypt key/val: %v", err)
		}
		encryptedConfig.Secrets = append(encryptedConfig.Secrets, encryptedKeyVal)
	}

	return
}

func (plaintext *PlaintextSecretKeyValue) encryptPlaintextKeyValue(kmsKeyName string) (encrypted EncryptedSecretKeyValue, err error) {
	encrypted.Key = plaintext.Key
	encryptedVal, err := gcloudEncryptPlaintext(kmsKeyName, plaintext.Value)
	if err != nil {
		err = fmt.Errorf("encryption request failed: %v", err)
	}
	encrypted.B64EncryptedValue = b64Encode(encryptedVal)
	return
}

func gcloudEncryptPlaintext(kmsKeyName string, plaintext string) (encrypted []byte, err error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		err = fmt.Errorf("could not generate kms client: %v", err)
		return
	}

	// Build the request.
	req := &kmspb.EncryptRequest{
		Name:      kmsKeyName,
		Plaintext: []byte(plaintext),
	}
	// Call the API.
	resp, err := client.Encrypt(ctx, req)
	if err != nil {
		err = fmt.Errorf("API error: %v", err)
		return
	}
	encrypted = resp.Ciphertext
	return
}

func b64Encode(encrypted []byte) (string) {
	return base64.StdEncoding.EncodeToString(encrypted)
}