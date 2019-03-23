package secretConfig

import (
	"encoding/base64"
	"fmt"
	cloudkms "cloud.google.com/go/kms/apiv1"
	"context"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"os"
)

func DecryptSecretConfig(path string, kmsKey string) (*PlaintetSecretConfig, error) {
	encryptedSecretConfig, err := readEncryptedSecretConfigFile(path)
	if err != nil {
		return nil, err
	}

	plaintextKeyVals := make([]PlaintextSecretKeyValue, 0)
	for _, encryptedKeyVal := range encryptedSecretConfig.Secrets {
		plaintextKeyVal, err := encryptedKeyVal.decryptEncryptedKeyVal(kmsKey)
		if err != nil {
			return nil, err
		}
		plaintextKeyVals = append(plaintextKeyVals, *plaintextKeyVal)
	}

	plaintextSecretConfig := PlaintetSecretConfig{
		Namespace: encryptedSecretConfig.Namespace,
		SecretName: encryptedSecretConfig.SecretName,
		PlaintextSecrets: plaintextKeyVals,
	}

	return &plaintextSecretConfig, nil
}

func DecryptSinglefile(path string, kmsKey string) (plaintextSingleFile *PlaintextSingleFile, err error) {

	_, err = os.Stat(path)
	if os.IsNotExist(err) || err != nil {
		return
	}

	readBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	var encryptedSingleFile EncryptedSingleFile
	err = yaml.Unmarshal(readBytes, encryptedSingleFile)

	encryptedB64Decoded, err := b64Decode(encryptedSingleFile.EncryptedValue)
	if err != nil {
		return
	}

	plaintextBytes, err := gcloudDecryptEncryptedBytes(encryptedB64Decoded, kmsKey)

	plaintextSingleFile = &PlaintextSingleFile{
		Key: encryptedSingleFile.Key,
		Namespace: encryptedSingleFile.Namespace,
		PlaintextVaule: plaintextBytes,
	}

	return
}

func gcloudDecryptEncryptedBytes(encrypted []byte, kmsKey string) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       kmsKey,
		Ciphertext: encrypted,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Plaintext, nil
}

func readEncryptedSecretConfigFile(path string) (*EncryptedSecretConfig, error) {
	encryptedSecretConfigBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	var encryptedSecretConfig EncryptedSecretConfig

	err = yaml.Unmarshal(encryptedSecretConfigBytes, &encryptedSecretConfig)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshall yaml: %v", err)
	}

	return &encryptedSecretConfig, nil
}

func (encryptedKeyVal *EncryptedSecretKeyValue) decryptEncryptedKeyVal(kmsKey string) (*PlaintextSecretKeyValue, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	encryptedB64DecodedVal, err := encryptedKeyVal.b64Decode()
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       kmsKey,
		Ciphertext: encryptedB64DecodedVal,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}

	plaintextSecretKeyVal := PlaintextSecretKeyValue{Key: encryptedKeyVal.Key, Value: resp.Plaintext}

	return &plaintextSecretKeyVal, nil
}

func (keyVal *EncryptedSecretKeyValue) b64Decode() (b64Decoded []byte, err error) {
	return base64.StdEncoding.DecodeString(keyVal.B64EncryptedValue)
}

func b64Decode(b64Encoded string) (b64Decoded []byte, err error) {
	return base64.StdEncoding.DecodeString(b64Encoded)
}
