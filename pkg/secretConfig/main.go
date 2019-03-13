package secretConfig

type EncryptedSecretConfig struct {
	Namespace string `yaml:namespace`
	SecretName string `yaml:secretName`
	Secrets []EncryptedSecretKeyValue `yaml:secrets`
}

type EncryptedSecretKeyValue struct {
	Key string `yaml:key`
	B64EncryptedValue string `yaml:b64EncryptedValue`
}

type PlaintetSecretConfig struct {
	Namespace string `yaml:namespace`
	SecretName string `yaml:secretName`
	PlaintextSecrets []PlaintextSecretKeyValue `yaml:plaintextSecrets`
}

type PlaintextSecretKeyValue struct {
	Key string `yaml:key`
	Value []byte `yaml:plaintextVal`
}

