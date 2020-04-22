package kube

import (
	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestData(t *testing.T) {
	Convey("Create data for configmap", t, func() {
		Convey("non nil config", func() {
			kt := v1alpha1.KafkaTopic{
				ObjectMeta: v1.ObjectMeta{
					Name: "test-topic-1",
				},
				Spec: v1alpha1.KafkaTopicSpec{
					Partitions:        1,
					ReplicationFactor: 1,
					Config: map[string]string{
						"key1": "value1",
					},
					TargetCluster: v1alpha1.ClusterConnection{
						Name: "test-cluster",
					},
				},
			}
			data, e := data(kt)
			So(e, ShouldEqual, nil)
			So(data["partitions"], ShouldEqual, "1")
			So(data["replicationFactor"], ShouldEqual, "1")
			So(data["key1"], ShouldEqual, "value1")
			So(data["target-cluster"], ShouldEqual, "test-cluster")
			So(data["topic-name"], ShouldEqual, "test-topic-1")

		})
		Convey("nil config", func() {
			kt := v1alpha1.KafkaTopic{
				Spec: v1alpha1.KafkaTopicSpec{
					Partitions:        2,
					ReplicationFactor: 3,
					Config:            nil,
					TargetCluster: v1alpha1.ClusterConnection{
						Name: "test-cluster",
					},
				},
			}
			data, e := data(kt)
			So(e, ShouldEqual, nil)
			So(data["partitions"], ShouldEqual, "2")
			So(data["replicationFactor"], ShouldEqual, "3")
			So(data["target-cluster"], ShouldEqual, "test-cluster")
		})
	})
}

func TestNewConfigmap(t *testing.T) {
	Convey("Creating configmap definition.", t, func() {
		Convey("nil config", func() {
			kt := v1alpha1.KafkaTopic{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-topic",
					Namespace: "test",
				},
				Spec: v1alpha1.KafkaTopicSpec{
					Partitions:        2,
					ReplicationFactor: 3,
					Config:            nil,
					TargetCluster: v1alpha1.ClusterConnection{
						Name: "test-connection",
					},
				},
			}
			configmap, e := NewConfigmap(kt)
			So(e, ShouldEqual, nil)
			So(configmap.Name, ShouldEqual, "test-topic")
			So(configmap.Namespace, ShouldEqual, "test")
			So(configmap.Data["partitions"], ShouldEqual, "2")
			So(configmap.Data["replicationFactor"], ShouldEqual, "3")
			So(configmap.Data["target-cluster"], ShouldEqual, "test-connection")
		})
		Convey("non-nil config", func() {
			kt := v1alpha1.KafkaTopic{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-topic",
					Namespace: "test",
				},
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
			configmap, e := NewConfigmap(kt)
			So(e, ShouldEqual, nil)
			So(configmap.Name, ShouldEqual, "test-topic")
			So(configmap.Namespace, ShouldEqual, "test")
			So(configmap.Data["partitions"], ShouldEqual, "2")
			So(configmap.Data["replicationFactor"], ShouldEqual, "3")
			So(configmap.Data["target-cluster"], ShouldEqual, "test-connection")
			So(configmap.Data["key1"], ShouldEqual, "value1")
			So(configmap.Data["topic-name"], ShouldEqual, "test-topic")
		})
		Convey("non-nil config, extra labels", func() {
			m := make(map[string]string)
			m["test-key"] = "test-value"
			kt := v1alpha1.KafkaTopic{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-topic",
					Namespace: "test",
					Labels:    m,
				},
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
			configmap, e := NewConfigmap(kt)
			So(e, ShouldEqual, nil)
			So(configmap.Name, ShouldEqual, "test-topic")
			So(configmap.Namespace, ShouldEqual, "test")
			So(configmap.Data["partitions"], ShouldEqual, "2")
			So(configmap.Data["replicationFactor"], ShouldEqual, "3")
			So(configmap.Data["target-cluster"], ShouldEqual, "test-connection")
			So(configmap.Data["key1"], ShouldEqual, "value1")
			So(configmap.Data["topic-name"], ShouldEqual, "test-topic")
			So(configmap.Labels["test-key"], ShouldEqual, "test-value")
		})
	})
}
