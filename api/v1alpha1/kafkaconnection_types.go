package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// KafkaConnectionSpec defines the desired state of KafkaConnection
type KafkaConnectionSpec struct {
	Broker     string            `json:"broker"`
	Username   string            `json:"username,omitempty"`
	Password   string            `json:"password,omitempty"`
	AuthMethod string            `json:"auth-method,omitempty"`
	Config     map[string]string `json:"config,omitempty"`
}

// KafkaConnectionStatus defines the observed state of KafkaConnection
type KafkaConnectionStatus struct {
	State KafkaConnectionState `json:"status"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// KafkaConnection is the Schema for the kafkaconnections API
type KafkaConnection struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   KafkaConnectionSpec   `json:"spec,omitempty"`
	Status KafkaConnectionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// KafkaConnectionList contains a list of KafkaConnection
type KafkaConnectionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KafkaConnection `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KafkaConnection{}, &KafkaConnectionList{})
}
