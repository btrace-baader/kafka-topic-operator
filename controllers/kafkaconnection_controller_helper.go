package controllers

import (
	"context"

	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
	"github.com/btrace-baader/kafka-topic-operator/kube"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (r *KafkaConnectionReconciler) manageSecret(log logr.Logger, ctx context.Context, req ctrl.Request, kafkaConnection *kafkav1alpha1.KafkaConnection, namespace string) error {
	// create a new secret object
	secret, err := kube.NewSecret(*kafkaConnection, namespace)
	if err != nil {
		return err
	}
	// add KafkaConnection as the owner of resulting secret
	if err := controllerutil.SetControllerReference(kafkaConnection, secret, r.Scheme); err != nil {
		return err
	}
	if err := r.createOrUpdateSecretInKube(ctx, req, kafkaConnection, secret, namespace); err != nil {
		return err
	}
	return nil
}

func (r *KafkaConnectionReconciler) createOrUpdateSecretInKube(ctx context.Context, req ctrl.Request, kafkaConnection *kafkav1alpha1.KafkaConnection, secret *v1.Secret, namespace string) error {
	// create secret (secret, namespace string, kafkaconnection KafkaConnectionSpec)
	err := r.Client.Get(ctx, types.NamespacedName{Name: kafkaConnection.Name, Namespace: namespace}, &v1.Secret{})
	if err != nil && apierrs.IsNotFound(err) {
		logrus.Infof("creating a new secret %v/%v", namespace, kafkaConnection.Name)
		err = r.Client.Create(ctx, secret)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// secret already exists, update it
		logrus.Infof("updating secret %v/%v", namespace, kafkaConnection.Name)
		err = r.Client.Update(ctx, secret)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *KafkaConnectionReconciler) getNamespaces(ctx context.Context) (*v1.NamespaceList, error) {
	namespace := &v1.NamespaceList{}
	err := r.Client.List(ctx, namespace)
	if err != nil {
		return &v1.NamespaceList{}, err
	}
	return namespace, nil
}

func (r *KafkaConnectionReconciler) updateState(log logr.Logger, ctx context.Context, connection *kafkav1alpha1.KafkaConnection, state kafkav1alpha1.KafkaConnectionState) {
	connection.Status = kafkav1alpha1.KafkaConnectionStatus{State: state}
	if err := r.Status().Update(ctx, connection); err != nil {
		log.Error(err, "failed to update kafkaConnection status")
	}
}
