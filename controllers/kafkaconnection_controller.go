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
	"github.com/go-logr/logr"
)

// KafkaConnectionReconciler reconciles a KafkaConnection object
type KafkaConnectionReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;update;list;watch;create;delete
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
//+kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkaconnections,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkaconnections/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kafka.btrace.com,resources=kafkaconnections/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KafkaConnection object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *KafkaConnectionReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("kafkaconnection", req.NamespacedName)

	kafkaConnection := &kafkav1alpha1.KafkaConnection{}
	// Check if a resource of KafkaConnection type exists, it can be removed between
	// reconciliation loops so error is ignored and object reconciled.
	if err := r.Get(ctx, req.NamespacedName, kafkaConnection); err != nil {
		log.Info("unable to fetch KafkaConnection", "error", err)
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

// SetupWithManager sets up the controller with the Manager.
func (r *KafkaConnectionReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kafkav1alpha1.KafkaConnection{}).
		Complete(r)
}
