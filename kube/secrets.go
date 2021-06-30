package kube

import (
	"strings"

	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewSecret returns a pointer to secret created using KafkaConnection config
func NewSecret(kc v1alpha1.KafkaConnection, namespace string) (*v1.Secret, error) {
	labels := map[string]string{
		"managed-by": "kafkaConnection-operator",
	}
	stringData, err := stringData(kc)
	if err != nil {
		return &v1.Secret{}, err
	}
	return &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kc.Name,
			Namespace: namespace,
			Labels:    labels,
		},
		StringData: stringData,
		Type:       "Opaque",
	}, nil
}

// StringData converts the KafkaConnection specs into a secret data
func stringData(kc v1alpha1.KafkaConnection) (map[string]string, error) {
	stringData := make(map[string]string)
	stringData["brokers"] = strings.Join(kc.Spec.Brokers, ",")
	stringData["security-protocol"] = kc.Spec.SecurityProtocol
	stringData["username"] = kc.Spec.Username
	stringData["password"] = kc.Spec.Password
	stringData = removeEmpty(stringData)
	if len(kc.Spec.Config) == 0 {
		return stringData, nil
	}
	stringData, err := mergeMaps(stringData, kc.Spec.Config)
	if err != nil {
		return stringData, err
	}
	return stringData, nil
}
