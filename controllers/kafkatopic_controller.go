package controllers

import (
	"context"

	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	"github.com/btrace-baader/kafka-topic-operator/topic"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	kafkaTopicFinalizer = "kafka.btrace.com"
)

// KafkaTopicReconciler reconciles a KafkaTopic object
type KafkaTopicReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkatopics,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkatopics/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;update;list;watch;create;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;update;list;watch

func (r *KafkaTopicReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("kafkatopic", req.NamespacedName)

	kafkaTopic := &kafkav1alpha1.KafkaTopic{}
	// Check if a KafkaTopic type resource exists, it might be that it's deleted
	// in between reconciliation so error will be ignored.
	if err := r.Get(ctx, req.NamespacedName, kafkaTopic); err != nil {
		log.Info("unable to fetch kafkaTopic", "error", err)
		return ctrl.Result{}, ignoreNotFound(err)
	}

	// Get credentials config from KafkaConnection
	kafkaConnection, err := getKafkaConnection(ctx, r.Client, kafkaTopic.Spec.TargetCluster.Name)
	if err != nil {
		return requeueWithError(log, "unable to get KafkaConnection object", err)
	}

	kclient := &topic.KafkaClient{}
	// set the logger for topic for consistent log format
	kclient.Log = log

	err = kclient.Init(kafkaConnection)
	if err != nil {
		r.updateState(log, ctx, kafkaTopic, kafkav1alpha1.TOPIC_CONNECTION_ERROR)
		return requeueWithError(log, "failed to create connection to kafka", err)
	}

	// Close connection to kafka client
	defer func() {
		err := kclient.Close()
		if err != nil {
			log.Error(err, "failed to close connection")
		}
	}()

	err = r.createOrUpdateTopicOnTarget(log, kclient, kafkaTopic)
	if err != nil {
		r.updateState(log, ctx, kafkaTopic, kafkav1alpha1.TOPIC_CREATION_ERROR)
		return requeueWithError(log, "failed to create/update topic", err)
	}

	// check and create/update configmap on kube cluster
	if err := r.manageConfigmap(ctx, req, kafkaTopic); err != nil {
		r.updateState(log, ctx, kafkaTopic, kafkav1alpha1.CONFIGMAP_CREATION_ERROR)
		return requeueWithError(log, "failure in configmap creation/update", err)
	}

	// set finalizers only if termination policy is DeleteAll
	// finalizers are used here for deleting topic on kafka cluster
	if kafkaTopic.Spec.TerminationPolicy == kafkav1alpha1.DELETE_ALL {
		if err := r.checkAndRunFinalizers(log, ctx, kclient, kafkaTopic); err != nil {
			r.updateState(log, ctx, kafkaTopic, kafkav1alpha1.TOPIC_DELETE_ERROR)
			return requeueWithError(log, "failure in running finalizers", err)
		}
	}

	r.updateState(log, ctx, kafkaTopic, kafkav1alpha1.TOPIC_CREATED)
	// rerun the reconcile loop after 120 seconds
	return requeueWithTimeout(120)
}

func (r *KafkaTopicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kafkav1alpha1.KafkaTopic{}).
		Complete(r)
}
