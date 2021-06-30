package kube

import (
	"testing"

	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	. "github.com/smartystreets/goconvey/convey"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestStringData(t *testing.T) {
	Convey("Create data for secret", t, func() {
		Convey("non nil config", func() {
			var kc = v1alpha1.KafkaConnection{
				Spec: v1alpha1.KafkaConnectionSpec{
					Brokers:          []string{"10.130.67.52:9092", "10.130.67.52:9092"},
					Username:         "user-1",
					Password:         "password-1",
					SecurityProtocol: "testMethod",
					Config: map[string]string{
						"key1": "value1",
					},
				},
			}
			stringData, e := stringData(kc)
			So(e, ShouldEqual, nil)
			So(stringData["brokers"], ShouldEqual, "10.130.67.52:9092,10.130.67.52:9092")
			So(stringData["security-protocol"], ShouldEqual, "testMethod")
			So(stringData["username"], ShouldEqual, "user-1")
			So(stringData["password"], ShouldEqual, "password-1")
			So(stringData["key1"], ShouldEqual, "value1")

		})
		Convey("nil config", func() {
			var kc = v1alpha1.KafkaConnection{
				Spec: v1alpha1.KafkaConnectionSpec{
					Brokers:          []string{"10.130.67.52:9092", "10.130.67.52:9092"},
					Username:         "user-1",
					Password:         "password-1",
					SecurityProtocol: "testMethod",
					Config:           nil,
				},
			}
			stringData, e := stringData(kc)
			So(e, ShouldEqual, nil)
			So(stringData["brokers"], ShouldEqual, "10.130.67.52:9092,10.130.67.52:9092")
			So(stringData["security-protocol"], ShouldEqual, "testMethod")
			So(stringData["username"], ShouldEqual, "user-1")
			So(stringData["password"], ShouldEqual, "password-1")
		})
	})
}

func TestNewSecret(t *testing.T) {
	Convey("Creating configmap definition.", t, func() {
		Convey("nil config", func() {
			namespace := "not-test-ns"
			var kc = v1alpha1.KafkaConnection{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "test-ns",
				},
				Spec: v1alpha1.KafkaConnectionSpec{
					Brokers:          []string{"10.130.67.52:9092", "10.130.67.52:9092"},
					Username:         "user-1",
					Password:         "password-1",
					SecurityProtocol: "testMethod",
					Config:           nil,
				},
			}
			secret, e := NewSecret(kc, namespace)
			So(e, ShouldEqual, nil)
			So(secret.Name, ShouldEqual, "test-secret")
			So(secret.Namespace, ShouldEqual, "not-test-ns")
			So(secret.StringData["brokers"], ShouldEqual, "10.130.67.52:9092,10.130.67.52:9092")
			So(secret.StringData["security-protocol"], ShouldEqual, "testMethod")
			So(secret.StringData["username"], ShouldEqual, "user-1")
			So(secret.StringData["password"], ShouldEqual, "password-1")
		})
		Convey("non-nil config", func() {
			namespace := "not-test-ns"
			var kc = v1alpha1.KafkaConnection{
				ObjectMeta: v1.ObjectMeta{
					Name:      "test-secret",
					Namespace: "test-ns",
				},
				Spec: v1alpha1.KafkaConnectionSpec{
					Brokers:          []string{"10.130.67.52:9092", "10.130.67.52:9092"},
					Username:         "user-1",
					Password:         "password-1",
					SecurityProtocol: "testMethod",
					Config: map[string]string{
						"key1": "value1",
					},
				},
			}
			secret, e := NewSecret(kc, namespace)
			So(e, ShouldEqual, nil)
			So(secret.Name, ShouldEqual, "test-secret")
			So(secret.Namespace, ShouldEqual, "not-test-ns")
			So(secret.StringData["brokers"], ShouldEqual, "10.130.67.52:9092,10.130.67.52:9092")
			So(secret.StringData["security-protocol"], ShouldEqual, "testMethod")
			So(secret.StringData["username"], ShouldEqual, "user-1")
			So(secret.StringData["password"], ShouldEqual, "password-1")
			So(secret.StringData["key1"], ShouldEqual, "value1")
		})
	})
}
