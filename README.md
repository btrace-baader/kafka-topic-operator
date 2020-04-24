# kafka-topic-operator

This project contains two custom resource definitions and their respective resources: KafkaConnection and KafkaTopic. 
## Motivation
Please check out [this](https://medium.com/@btracebaader/kafka-topic-operator-how-to-manage-kafka-topics-in-kubernetes-d0290c8341dd) post 
on medium where we introduce this operator.

## Overview
KafkaConnection stores the connection information for the cluster such as broker address, authentication method and additional 
(optional) connection configurations. 
Following are the fields supported by KafkaConnection object:
```cassandraql
	Broker              string             Required
	Username            string             Optional
	Password            string             Optional
	SecurityProtocol    string             Optional
	Config              map[string]string  Optional
```
KafkaConnection also creates corresponding secrets in all namespaces which contain the 
information specified in KafkaConnection. Purpose of the secret is to provide uniform credentials for kafka cluster
that can be used by different services.

KafkaTopic defines a topic that will be created on the target cluster (can be in-cluster or an external service). It also creates
a corresponding ConfigMap with all the configuration. KafkaTopic also manages the lifecycle of the ConfigMap.
Definition of the KafkaTopic has the following fields:
```cassandraql
	Partitions           int32              Required
	ReplicationFactor    int16              Required
	Config               map[string]string  Optional
	TargetCluster        ClusterConnection  Required
```
KafkaTopic creates a configmap with all the configurations which can be used by services deployed
in the same namespace.

Once the CRD and controller are deployed on a cluster, objects of type KafkaConnection and KafkaTopic can be created. Following is a sample configuration of
KafkaConnection:
```cassandraql
apiVersion: kafka.btrace.com/v1alpha1
kind: KafkaConnection
metadata:
  name: test-connection-1
spec:
  broker: "10.130.67.52:9092"
  security-protocol: "SASL_SSL"
  username: "user-1"
  password: "weakP0ssw@rd"
```

KafkaTopic can be deployed on the cluster defined by KafkaConnection (test-connection-1) using the following sample config:
```cassandraql
apiVersion: kafka.btrace.com/v1alpha1
kind: KafkaTopic
metadata:
  name: test-topic-1
  namespace: test
spec:
  partitions: 1
  replication-factor: 3
  target-cluster:
    name: test-connection-1
  config:
    "cleanup.policy": "compact"
    "min.compaction.lag.ms": "86400000"
    "max.compaction.lag.ms": "432000000"
```

These sample configurations can also be found under config/samples. 

### Secrets and ConfigMaps
Secrets created by a KafkaConnection definition are present in all namespaces and hence the cluster configurations
can be accessed from any namespace. However, the ConfigMap created by a KafkaTopic is present only in the namespace where KafkaTopic is deployed. 
The reason being that KafkaConnection is logically a cluster wide resource but KafkaTopic is specific to a
set of services/namespace. 

Secrets have the same name as the KafkaConnection object, following is a sample secret created by applying the above config for KafkaConnection:
```cassandraql
apiVersion: v1
data:
  security-protocol: U0FTTA==
  broker: MTAuMTMwLjY3LjUyOjkwOTI=
  password: d2Vha1Awc3N3QHJk
  username: dXNlci0x
kind: Secret
metadata:
  labels:
    managed-by: kafkaConnection-operator
  name: test-connection-1
  namespace: default
  ownerReferences:
  - apiVersion: kafka.btrace.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: KafkaConnection
    name: test-connection-1
type: Opaque
```
Note that the data is base64 encoded. The data in the secret can be referenced like any regular secret.
For example, the password can be used to set the value of an env var in a deployment. Snippet of an example
deployment would look like this:
```cassandraql
env:
- name: KAFKA_PASSWORD
  valueFrom:
    secretKeyRef:
      key: password
      name: test-connection-1
```

Similarly following ConfigMap would be created in test namespace once we apply the sample KafkaTopic
config given in the earlier section:
```cassandraql
apiVersion: v1
data:
  cleanup.policy: compact
  target-cluster: test/test-connection-1
  min.compaction.lag.ms: "86400000"
  max.compaction.lag.ms: "432000000"
  partitions: "1"
  replication-factor: "3"
  topic-name: test-topic-1
kind: ConfigMap
metadata:
  labels:
    managed-by: kafkaTopic-operator
  name: test-topic-1
  namespace: spaghetti-staging
  ownerReferences:
  - apiVersion: kafka.btrace.com/v1alpha1
    blockOwnerDeletion: true
    controller: true
    kind: KafkaTopic
    name: test-topic-1
```

The data in this ConfigMap can be used in other services. An example of using the replication-factor as env
var in a deployment is given below:

 ```cassandraql
env:
 - name: REPLICATION_FACTOR
   valueFrom:
     configMapKeyRef:
       key: replication-factor
       name: test-topic-1
 ```

### Termination Policy
In order to prevent accidental deletion of topics, a spec.terminationPolicy field is added. 
It decides what can be deleted. There are three possible values for termination policy: 

| Field   | Description |
| ------------- | ------------- |
| NotDeletable  | KafkaTopic is not deletable. |
| KeepTopic | KafkaTopic object is deletable and so is the resulting configmap, however the topic on kafka cluster is not deletable.|
| DeleteAll| KafkaTopic, configmap and topic on cluster are all deletable. |

By default the policy is set to NotDeletable but the object can be edited to set to a different policy based on 
requirements.

## Getting Started
This operator works with version v0.11.0 and above. Use instructions [here](https://cert-manager.io/docs/installation/kubernetes/) to set up the cert-manager.
Once cert manager is deployed, use the following steps to set up the operator:
```
git clone https://github.com/btrace-baader/kafka-topic-operator.git
cd kafka-topic-operator
make install
make deploy
```
If you want to build your own changes and test them locally, [kind](https://github.com/kubernetes-sigs/kind) is one of the choices for a lightweight Kubernetes setup. After you clone the repo and have made your changes, perform the following steps:
```
make install
make docker-build
kind load docker-image <image name:tag>
make deploy
```
After the successful setup, you should be able to create resources of the types KafkaConnection and KafkaTopic, as
described in the overview section.
