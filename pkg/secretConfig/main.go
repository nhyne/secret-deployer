package secretConfig

import "encoding/base64"

type EncryptedSecretConfig struct {
	Namespace string `yaml:namespace`
	SecretName string `yaml:secretName`
	Secrets []EncryptedSecretKeyValue `yaml:secrets`
}

type EncryptedSecretKeyValue struct {
	Key string `yaml:key`
	B64EncryptedValue string `yaml:b64EncryptedValue`
}

func (keyVal *EncryptedSecretKeyValue) B64Decode() (b64Decoded []byte, error error) {
	return base64.StdEncoding.DecodeString(keyVal.B64EncryptedValue)
}
