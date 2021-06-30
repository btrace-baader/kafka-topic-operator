module github.com/btrace-baader/kafka-topic-operator

go 1.16

require (
	github.com/Shopify/sarama v1.24.1
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/prometheus/common v0.4.1
	github.com/sirupsen/logrus v1.2.0
	github.com/smartystreets/goconvey v1.6.4
	golang.org/x/tools v0.0.0-20190328211700-ab21143f2384
	k8s.io/api v0.0.0-20190918195907-bd6ac527cfd2
	k8s.io/apimachinery v0.0.0-20190817020851-f2f3a405f61d
	k8s.io/client-go v0.0.0-20190918200256-06eb1244587a
	sigs.k8s.io/controller-runtime v0.3.0
)
