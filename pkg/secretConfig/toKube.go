package secretConfig

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConvertSecretConfigToSecretObject(configPath string, kmsKey string) (secretObject *corev1.Secret, namespace string, err error) {
	plaintextSecretConfig, err := DecryptSecretConfig(configPath, kmsKey)
	if err != nil {
		return nil, "", err
	}

	secretKeyVals := make(map[string][]byte)
	// generate map[string][]byte
	for _, plaintextKeyVal := range plaintextSecretConfig.PlaintextSecrets {

		secretKeyVals[plaintextKeyVal.Key] = []byte(plaintextKeyVal.Value)
	}

	secretObject = createSecretObject(plaintextSecretConfig.SecretName, plaintextSecretConfig.Namespace, make(map[string]string, 0), secretKeyVals)

	return secretObject, plaintextSecretConfig.Namespace, nil
}

func createSecretObject(secretName string, secretNamespace string, secretLabels map[string]string, secretKeyVals map[string][]byte) (*corev1.Secret) {
	return &corev1.Secret{
		Type: corev1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
			Labels: secretLabels,
			Namespace: secretNamespace,
		},
		Data: secretKeyVals,
	}
}