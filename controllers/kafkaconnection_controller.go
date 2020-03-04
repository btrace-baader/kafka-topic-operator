package controllers

import (
	"context"
	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KafkaConnectionReconciler reconciles a KafkaConnection object
type KafkaConnectionReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkaconnections,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkaconnections/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;update;list;watch;create;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch

func (r *KafkaConnectionReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("kafkaconnection", req.NamespacedName)

	kafkaConnection := &kafkav1alpha1.KafkaConnection{}
	// Check if a resource of KafkaConnection type exists, it can be removed between
	// reconciliation loops so error is ignored and object reconciled.
	if err := r.Get(ctx, req.NamespacedName, kafkaConnection); err != nil {
		log.Error(err, "unable to fetch KafkaConnection")
		return ctrl.Result{}, ignoreNotFound(err)
	}
	namespaces, err := r.getNamespaces(ctx)
	if err != nil {
		log.Error(err, "failed to fetch namespaces")
	}

	for _, ns := range namespaces.Items {
		log.Info(ns.Name)
		if err := r.manageSecret(log, ctx, req, kafkaConnection, ns.Name); err != nil {
			r.updateState(log, ctx, kafkaConnection, kafkav1alpha1.CONNECTION_ERROR)
			return requeueWithError(log, "failed to create/update secret", err)
		}
	}

	r.updateState(log, ctx, kafkaConnection, kafkav1alpha1.CONNECTION_CREATED)

	//TODO(saniazara01): create a watcher for new namespaces
	return requeueWithTimeout(600)
}
