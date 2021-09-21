package controllers

import (
	"context"
	"time"

	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("KafkaConnection Controller", func() {
	const timeout = time.Second * 30
	const interval = time.Second * 1

	Context("Kafka Connection create ", func() {
		It("Should create kafka connection successfully", func() {
			testKafkaConnection := types.NamespacedName{
				Namespace: "default",
				Name:      "test-connection",
			}
			spec := kafkav1alpha1.KafkaConnectionSpec{
				Brokers:          []string{"10.2.10.10:9092", "10.2.10.10:9092"},
				Username:         "",
				Password:         "",
				SecurityProtocol: "",
				Config:           nil,
			}

			kafkaconnection := &kafkav1alpha1.KafkaConnection{
				ObjectMeta: v1.ObjectMeta{
					Name:      testKafkaConnection.Name,
					Namespace: testKafkaConnection.Namespace,
				},
				Spec: spec,
			}

			Expect(k8sClient.Create(context.Background(), kafkaconnection)).Should(Succeed())

			By("Expect deletion")
			Eventually(func() error {
				kc := &kafkav1alpha1.KafkaConnection{}
				_ = k8sClient.Get(context.Background(), testKafkaConnection, kc)
				return k8sClient.Delete(context.Background(), kc)
			}, timeout, interval).Should(Succeed())

		})

	})
})
