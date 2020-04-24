package controllers

import (
	"context"
	"github.com/btrace-baader/kafka-topic-operator/kube"
	"github.com/btrace-baader/kafka-topic-operator/topic"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	kafkav1alpha1 "github.com/btrace-baader/kafka-topic-operator/api/v1alpha1"
)

func (r *KafkaTopicReconciler) manageConfigmap(ctx context.Context, req ctrl.Request, kafkaTopic *kafkav1alpha1.KafkaTopic) error {
	log := r.Log.WithValues("configmap", req.NamespacedName)
	// create configmap object
	configMap, err := kube.NewConfigmap(*kafkaTopic)
	if err != nil {
		return err
	}
	// add KafkaTopic as the owner of the configmap
	if err := controllerutil.SetControllerReference(kafkaTopic, configMap, r.Scheme); err != nil {
		return err
	}
	err = r.createOrUpdateConfigmapInKube(log, ctx, req, kafkaTopic, configMap)
	if err != nil {
		return err
	}
	return nil
}

// create a new configmap from the provided definition, update if it already exists
func (r *KafkaTopicReconciler) createOrUpdateConfigmapInKube(log logr.Logger, ctx context.Context, req ctrl.Request, kafkaTopic *kafkav1alpha1.KafkaTopic, configMap *v1.ConfigMap) error {
	err := r.Client.Get(ctx, types.NamespacedName{Name: kafkaTopic.Name, Namespace: req.Namespace}, &v1.ConfigMap{})
	if err != nil && apierrs.IsNotFound(err) {
		log.Info("creating a new configmap")
		err = r.Client.Create(ctx, configMap)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	} else {
		// configmap already exists, update it
		log.Info("updating configmap in namespace ", kafkaTopic.Name, kafkaTopic.Namespace)
		err = r.Client.Update(ctx, configMap)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *KafkaTopicReconciler) createOrUpdateTopicOnTarget(log logr.Logger, kclient *topic.KafkaClient, kafkaTopic *kafkav1alpha1.KafkaTopic) error {
	// check if topic exists
	topicExists, err := kclient.Exists(kafkaTopic.Name)
	if err != nil {
		return err
	}
	// create topic on target if it doesn't exist
	if !topicExists {
		log.Info("creating topic")
		err = kclient.Create(kafkaTopic)
		if err != nil {
			return err
		}
	} else {
		// update the topic
		err = kclient.Update(kafkaTopic)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *KafkaTopicReconciler) checkAndRunFinalizers(log logr.Logger, ctx context.Context, kclient *topic.KafkaClient, kafkaTopic *kafkav1alpha1.KafkaTopic) error {
	// check if KafkaTopic is under deletion
	if kafkaTopic.ObjectMeta.DeletionTimestamp.IsZero() {
		// add our finalizer if it doesn't exist
		if !containsString(kafkaTopic.ObjectMeta.Finalizers, kafkaTopicFinalizer) {
			log.Info("adding finalizer to the kafkaTopic")
			kafkaTopic.ObjectMeta.Finalizers = append(kafkaTopic.ObjectMeta.Finalizers, kafkaTopicFinalizer)
			if err := r.Update(ctx, kafkaTopic); err != nil {
				return err
			}
		}
	} else {
		if containsString(kafkaTopic.ObjectMeta.Finalizers, kafkaTopicFinalizer) {
			// finalizer is present so let's delete the topic
			log.Info("deleting the topic from target")
			if err := kclient.DeleteIfExists(kafkaTopic.Name); err != nil {
				return err
			}
			// remove finalizer from the list
			kafkaTopic.ObjectMeta.Finalizers = removeString(kafkaTopic.ObjectMeta.Finalizers, kafkaTopicFinalizer)
			// update the object
			if err := r.Update(ctx, kafkaTopic); err != nil {
				return err
			}
		}
	}
	return nil
}

func getKafkaConnection(ctx context.Context, client client.Client, clusterName string) (*kafkav1alpha1.KafkaConnection, error) {
	cluster := &kafkav1alpha1.KafkaConnection{}
	err := client.Get(ctx, types.NamespacedName{Name: clusterName}, cluster)
	if err != nil {
		return cluster, err
	}
	return cluster, nil
}

func (r *KafkaTopicReconciler) updateState(log logr.Logger, ctx context.Context, topic *kafkav1alpha1.KafkaTopic, state kafkav1alpha1.KafkaTopicState) {
	topic.Status = kafkav1alpha1.KafkaTopicStatus{State: state}
	if err := r.Status().Update(ctx, topic); err != nil {
		log.Error(err, "failed to update KafkaTopic status")
	}
}
