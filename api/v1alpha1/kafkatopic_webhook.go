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

package v1alpha1

import (
	"errors"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var kafkatopiclog = logf.Log.WithName("kafkatopic-resource")

func (r *KafkaTopic) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-kafka-btrace-com-v1alpha1-kafkatopic,mutating=true,failurePolicy=fail,sideEffects=None,groups=kafka.btrace.com,resources=kafkatopics,verbs=create;update,versions=v1alpha1,name=mkafkatopic.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &KafkaTopic{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *KafkaTopic) Default() {
	kafkatopiclog.Info("default", "name", r.Name)

	// the default termination policy is 'NotDeletable'
	if r.Spec.TerminationPolicy == "" {
		r.Spec.TerminationPolicy = NOT_DELETABLE
	}
}

//+kubebuilder:webhook:path=/validate-kafka-btrace-com-v1alpha1-kafkatopic,mutating=false,failurePolicy=fail,sideEffects=None,groups=kafka.btrace.com,resources=kafkatopics,verbs=create;update,versions=v1alpha1,name=vkafkatopic.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &KafkaTopic{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaTopic) ValidateCreate() error {
	kafkatopiclog.Info("validate create", "name", r.Name)
	if err := r.validateKafkaTopic(); err != nil {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaTopic) ValidateUpdate(old runtime.Object) error {
	kafkatopiclog.Info("validate update", "name", r.Name)
	if err := r.validateKafkaTopic(); err != nil {
		return err
	}
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaTopic) ValidateDelete() error {
	kafkatopiclog.Info("validate delete", "name", r.Name)
	if r.isNotDeletable() {
		return errors.New("topic has spec.terminationPolicy set to NotDeletable")
	}
	return nil
}

func (r *KafkaTopic) validateKafkaTopic() error {
	var allErrs field.ErrorList
	if err := r.validateKafkaTopicName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateKafkaTopicSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}
	return apierrors.NewInvalid(
		schema.GroupKind{Group: "kafka.btrace.com", Kind: "kafkatopic"},
		r.Name, allErrs)
}

func (r *KafkaTopic) validateKafkaTopicName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength {
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 63 characters")
	}
	return nil
}

func (r *KafkaTopic) validateKafkaTopicSpec() *field.Error {
	spec := field.NewPath("spec")
	if !r.replicationFactorOkay() {
		return field.Invalid(spec.Child("replicationFactor"), "replication-factor", "must be greater than 3")
	}
	if !r.partitionsOkay() {
		return field.Invalid(spec.Child("partitions"), "partitions", "must be greater than 0")
	}
	if !r.targetClusterOkay() {
		return field.Invalid(spec.Child("target-cluster"), "target-cluster", "target-cluster/kafkaconnection must be defined")
	}
	if !r.terminationPolicyOkay() {
		return field.Invalid(spec.Child("terminationPolicy"), "termnationPolicy", "possible values 'KeepTopic', 'DeleteAll', 'NotDeletable'")
	}
	return nil
}

func (r *KafkaTopic) replicationFactorOkay() bool {
	if r.Spec.ReplicationFactor >= 3 {
		return true
	}
	return false

}

func (r *KafkaTopic) partitionsOkay() bool {
	if r.Spec.Partitions > 0 {
		return true
	}
	return false
}

func (r *KafkaTopic) targetClusterOkay() bool {
	if len(r.Spec.TargetCluster.Name) > 0 {
		return true
	}
	return false
}

func (r *KafkaTopic) terminationPolicyOkay() bool {
	if r.Spec.TerminationPolicy == KEEP_TOPIC || r.Spec.TerminationPolicy == NOT_DELETABLE || r.Spec.TerminationPolicy == DELETE_ALL {
		return true
	}
	return false
}

func (r *KafkaTopic) isNotDeletable() bool {
	if r.Spec.TerminationPolicy == NOT_DELETABLE {
		return true
	}
	return false
}
