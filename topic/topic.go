package topic

import (
	"github.com/Shopify/sarama"
	"github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	"github.com/go-logr/logr"
)

type KafkaClient struct {
	admin sarama.ClusterAdmin
	Log   logr.Logger
}

// Init initialises the the admin client
func (client *KafkaClient) Init(kc *v1alpha1.KafkaConnection) error {
	client.Log.Info("creating a new cluster admin connection")
	admin, err := sarama.NewClusterAdmin(kc.Spec.Brokers, client.connectionConfig(kc))
	client.admin = admin
	if err != nil {
		client.Log.Error(err, "can not init connection ")
		return err
	}
	return nil
}

// Exists checks if a topic already exists on the target
func (client *KafkaClient) Exists(topic string) (bool, error) {
	topics, err := client.admin.ListTopics()
	if err != nil {
		return false, err
	}
	if _, ok := topics[topic]; ok {
		return true, nil
	}
	return false, nil
}

// Create creates a topic on the target
func (client *KafkaClient) Create(kt *v1alpha1.KafkaTopic) error {
	err := client.admin.CreateTopic(kt.Name, topicDetail(kt), false)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the topic from target
func (client *KafkaClient) Delete(topic string) error {
	err := client.admin.DeleteTopic(topic)
	if err != nil {
		return err
	}
	return nil
}

func (client *KafkaClient) DeleteIfExists(topic string) error {
	exists, err := client.Exists(topic)
	if err != nil {
		return err
	}
	if exists {
		client.Log.Info("topic exists, deleting...")
		if err := client.Delete(topic); err != nil {
			return err
		}
	} else {
		client.Log.Info("topic does not exist on the target")
	}
	return nil
}

// Update allows for changes in the topic config
func (client *KafkaClient) Update(kt *v1alpha1.KafkaTopic) error {
	//  TODO: CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error ? assignment
	// change according to config provided under spec.config
	err := client.admin.AlterConfig(sarama.TopicResource, kt.Name, stringToStringPointerMap(kt.Spec.Config), false)
	if err != nil {
		return err
	}
	ok, err := client.okayToCreatePartitions(kt.Name, kt.Spec.Partitions)
	if ok {
		var assign [][]int32
		client.Log.Info("updating partitions")
		err := client.admin.CreatePartitions(kt.Name, kt.Spec.Partitions, assign, false)
		if err != nil {
			return err
		}
	} else {
		client.Log.Info("requested partitions are the same as current. not updating")
		return nil
	}
	return err
}

// Close terminates the connection
func (client *KafkaClient) Close() error {
	client.Log.Info("closing connection to kafka")
	err := client.admin.Close()
	return err
}
