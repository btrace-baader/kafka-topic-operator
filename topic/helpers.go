package topic

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/Shopify/sarama"
	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	"io/ioutil"
	"strconv"
)

func (client *KafkaClient) connectionConfig(kc *v1alpha1.KafkaConnection) *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_0_1_0
	if kc.Spec.SecurityProtocol == "SASL" || kc.Spec.SecurityProtocol == "SASL_SSL" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = kc.Spec.Username
		config.Net.SASL.Password = kc.Spec.Password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		config.Net.SASL.Version = 1
		config.Net.SASL.Handshake = true
		config.Net.TLS.Enable = true
		config.Net.TLS.Config = client.tlsConfig()
	}
	return config
}

func topicDetail(kt *v1alpha1.KafkaTopic) *sarama.TopicDetail {
	return &sarama.TopicDetail{
		NumPartitions:     kt.Spec.Partitions,
		ReplicationFactor: kt.Spec.ReplicationFactor,
		ReplicaAssignment: nil,
		ConfigEntries:     stringToStringPointerMap(kt.Spec.Config),
	}
}

func (client *KafkaClient) tlsConfig() *tls.Config {
	caCert, err := ioutil.ReadFile("/etc/ssl/cert.pem")
	if err != nil {
		client.Log.Error(err, "failed to fetch cert.pem")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	return &tls.Config{
		RootCAs: caCertPool,
	}
}

func stringToStringPointerMap(in map[string]string) map[string]*string {
	out := make(map[string]*string, 0)
	for k, v := range in {
		out[k] = stringPointer(v)
	}
	return out
}

func stringPointer(str string) *string {
	return &str
}

// okayToCreatePartitions ensures that the number of desired partitions is valid
func (client *KafkaClient) okayToCreatePartitions(topic string, requested int32) (bool, error) {
	data, err := client.admin.DescribeTopics([]string{topic})
	if err != nil {
		return false, err
	}
	current := int32(len(data[0].Partitions))
	client.Log.Info("current:requested no of partitions", strconv.Itoa(int(current)), strconv.Itoa(int(requested)))
	if requested > current {
		client.Log.Info("ok to create partitions")
		return true, nil
	}
	return false, nil
}
