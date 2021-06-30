package topic

import (
	"testing"

	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	. "github.com/smartystreets/goconvey/convey"
)

func TestConnectionConfig(t *testing.T) {
	Convey("Test connection config", t, func() {
		Convey("Positive test", func() {
			kc := v1alpha1.KafkaConnection{
				Spec: v1alpha1.KafkaConnectionSpec{
					Brokers:          []string{"10.23.43.45:9092", "10.23.43.45:9092"},
					Username:         "user-1",
					Password:         "pass-1",
					SecurityProtocol: "SASL",
					Config:           nil,
				},
			}
			client := KafkaClient{}
			config := client.connectionConfig(&kc)
			So(config.Net.SASL.Enable, ShouldEqual, true)
			So(config.Net.SASL.User, ShouldEqual, "user-1")
			So(config.Net.SASL.Password, ShouldEqual, "pass-1")
			So(config.Net.TLS.Enable, ShouldEqual, true)
		})
		Convey("Positive test sasl_ssl", func() {
			kc := v1alpha1.KafkaConnection{
				Spec: v1alpha1.KafkaConnectionSpec{
					Brokers:          []string{"10.23.43.45:9092", "10.23.43.45:9092"},
					Username:         "user-1",
					Password:         "pass-1",
					SecurityProtocol: "SASL_SSL",
					Config:           nil,
				},
			}
			client := KafkaClient{}
			config := client.connectionConfig(&kc)
			So(config.Net.SASL.Enable, ShouldEqual, true)
			So(config.Net.SASL.User, ShouldEqual, "user-1")
			So(config.Net.SASL.Password, ShouldEqual, "pass-1")
			So(config.Net.TLS.Enable, ShouldEqual, true)
		})
		Convey("Negative test", func() {
			Convey("non-SASL security-protocol", func() {
				kc := v1alpha1.KafkaConnection{
					Spec: v1alpha1.KafkaConnectionSpec{
						Brokers:          []string{"10.23.43.45:9092", "10.23.43.45:9092"},
						Username:         "user-1",
						Password:         "pass-1",
						SecurityProtocol: "NOT-SASL",
						Config:           nil,
					},
				}
				client := KafkaClient{}
				config := client.connectionConfig(&kc)
				So(config.Net.SASL.Enable, ShouldEqual, false)
			})
			Convey("empty security-protocol", func() {
				kc := v1alpha1.KafkaConnection{
					Spec: v1alpha1.KafkaConnectionSpec{
						Brokers:          []string{"10.23.43.45:9092", "10.23.43.45:9092"},
						Username:         "user-1",
						Password:         "pass-1",
						SecurityProtocol: "",
						Config:           nil,
					},
				}
				client := KafkaClient{}
				config := client.connectionConfig(&kc)
				So(config.Net.SASL.Enable, ShouldEqual, false)
			})
		})
	})
}

func TestTopicDetail(t *testing.T) {
	Convey("Test topic detail", t, func() {
		Convey("empty config", func() {
			kt := v1alpha1.KafkaTopic{
				Spec: v1alpha1.KafkaTopicSpec{
					Partitions:        2,
					ReplicationFactor: 3,
					Config:            nil,
					TargetCluster: v1alpha1.ClusterConnection{
						Name: "test-connection",
					},
				},
			}
			td := topicDetail(&kt)
			So(td.ReplicationFactor, ShouldEqual, 3)
			So(td.NumPartitions, ShouldEqual, 2)
		})
		kt := v1alpha1.KafkaTopic{
			Spec: v1alpha1.KafkaTopicSpec{
				Partitions:        2,
				ReplicationFactor: 3,
				Config: map[string]string{
					"key1": "value1",
				},
				TargetCluster: v1alpha1.ClusterConnection{
					Name: "test-connection",
				},
			},
		}
		td := topicDetail(&kt)
		So(td.ReplicationFactor, ShouldEqual, 3)
		So(td.NumPartitions, ShouldEqual, 2)
		So(td.ConfigEntries, ShouldContainKey, "key1")
	})
}
