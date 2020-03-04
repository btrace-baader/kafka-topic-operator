package controllers

import (
	"context"
	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"
)

var _ = Describe("KafkaTopic Controller", func() {
	const timeout = time.Second * 30
	const interval = time.Second * 1

	Context("Kafka Topic create ", func() {
		It("Should create kafkatopic successfully", func() {
			testKafkaTopic := types.NamespacedName{
				Namespace: "default",
				Name:      "test-topic",
			}
			spec := kafkav1alpha1.KafkaTopicSpec{
				Partitions:        1,
				ReplicationFactor: 3,
				Config:            nil,
				ClusterRef: kafkav1alpha1.ClusterConnection{
					Name:      testKafkaTopic.Name,
					Namespace: testKafkaTopic.Namespace,
				},
			}

			kafkatopic := &kafkav1alpha1.KafkaTopic{
				ObjectMeta: v1.ObjectMeta{
					Name:      testKafkaTopic.Name,
					Namespace: testKafkaTopic.Namespace,
				},
				Spec: spec,
			}
			Expect(k8sClient.Create(context.Background(), kafkatopic)).Should(Succeed())

			By("Expect deletion")
			Eventually(func() error {
				kt := &kafkav1alpha1.KafkaTopic{}
				_ = k8sClient.Get(context.Background(), testKafkaTopic, kt)
				return k8sClient.Delete(context.Background(), kt)

			}, timeout, interval).Should(Succeed())
		})

	})
})
