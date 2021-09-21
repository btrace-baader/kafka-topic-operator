/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	"github.com/btrace-baader/kafka-topic-operator/topic"
	"github.com/go-logr/logr"
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

// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;update;list;watch;create;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;update;list;watch
//+kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkatopics,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkatopics/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkatopics/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KafkaTopic object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KafkaTopicReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
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

	// check if KafkaTopic is not under deletion
	if kafkaTopic.ObjectMeta.DeletionTimestamp.IsZero() {
		r.updateState(log, ctx, kafkaTopic, kafkav1alpha1.TOPIC_CREATED)
	}
	// rerun the reconcile loop after 120 seconds
	return requeueWithTimeout(120)
}

// SetupWithManager sets up the controller with the Manager.
func (r *KafkaTopicReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kafkav1alpha1.KafkaTopic{}).
		Complete(r)
}
