package v1alpha1

// KafkaConnectionState, possible states : created, error
type KafkaConnectionState string

const (
	CONNECTION_CREATED = "Created"
	CONNECTION_ERROR   = "Error"
)

// KafkaTopicState, possible states : TopicCreated , ConnectionError, TopicCreationError, ConfigMapCreationError, TopicDeleteError
type KafkaTopicState string

const (
	TOPIC_CREATED            = "TopicCreated"
	TOPIC_CONNECTION_ERROR   = "ConnectionError"
	TOPIC_CREATION_ERROR     = "TopicCreationError"
	CONFIGMAP_CREATION_ERROR = "ConfigMapCreationError"
	TOPIC_DELETE_ERROR       = "TopicDeleteError"
)

type ClusterConnection struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

// KafkaTopicTerminationPolicy, possible policies: NotDeletable, KeepTopic, DeleteAll

type KafkaTopicTerminationPolicy string

const (
	NOT_DELETABLE = "NotDeletable"
	KEEP_TOPIC    = "KeepTopic"
	DELETE_ALL    = "DeleteAll"
)
