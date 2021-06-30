package kube

import (
	"strconv"

	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NewConfigmap returns a pointer to ConfigMap created using KafkaTopic config
func NewConfigmap(kt v1alpha1.KafkaTopic) (*v1.ConfigMap, error) {
	labels := map[string]string{
		"managed-by": "kafkaTopic-operator",
	}
	data, err := data(kt)
	if err != nil {
		return &v1.ConfigMap{}, err
	}
	if kt.ObjectMeta.Labels != nil {
		labels, err = mergeMaps(labels, kt.Labels)
		if err != nil {
			return &v1.ConfigMap{}, err
		}
	}
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      kt.Name,
			Namespace: kt.Namespace,
			Labels:    labels,
		},
		Data: data,
	}, nil
}

// data converts the KafkaTopic spec into a configmap spec
func data(kt v1alpha1.KafkaTopic) (map[string]string, error) {
	data := make(map[string]string)
	data["partitions"] = strconv.Itoa(int(kt.Spec.Partitions))
	data["replicationFactor"] = strconv.Itoa(int(kt.Spec.ReplicationFactor))
	data["target-cluster"] = kt.Spec.TargetCluster.Name
	data["topic-name"] = kt.Name
	if len(kt.Spec.Config) == 0 {
		return data, nil
	}
	data, err := mergeMaps(kt.Spec.Config, data)
	if err != nil {
		return data, err
	}
	return data, nil
}
