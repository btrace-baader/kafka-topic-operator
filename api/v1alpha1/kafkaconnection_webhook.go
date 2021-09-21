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
var kafkaconnectionlog = logf.Log.WithName("kafkaconnection-resource")

func (r *KafkaConnection) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(saniazara01): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-kafka-btrace-com-v1alpha1-kafkaconnection,mutating=false,failurePolicy=fail,sideEffects=None,groups=kafka.btrace.com,resources=kafkaconnections,verbs=create;update,versions=v1alpha1,name=vkafkaconnection.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &KafkaConnection{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaConnection) ValidateCreate() error {
	kafkaconnectionlog.Info("validate create", "name", r.Name)
	if err := r.validateKafkaConnection(); err != nil {
		return err
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaConnection) ValidateUpdate(old runtime.Object) error {
	kafkaconnectionlog.Info("validate update", "name", r.Name)
	if err := r.validateKafkaConnection(); err != nil {
		return err
	}

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaConnection) ValidateDelete() error {
	kafkaconnectionlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *KafkaConnection) validateKafkaConnection() error {
	var allErrs field.ErrorList
	if err := r.validateKafkaConnectionName(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateKafkaConnectionSpec(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "kafka.btrace.com", Kind: "KafkaConnection"},
		r.Name, allErrs)
}

func (r *KafkaConnection) validateKafkaConnectionName() *field.Error {
	if len(r.ObjectMeta.Name) > validationutils.DNS1035LabelMaxLength {
		return field.Invalid(field.NewPath("metadata").Child("name"), r.Name, "must be no more than 63 characters")
	}
	return nil
}

func (r *KafkaConnection) validateKafkaConnectionSpec() *field.Error {
	spec := field.NewPath("spec")
	if !r.brokerDefined() {
		return field.Invalid(spec.Child("brokers"), r.Name, "must be defined")
	}
	if r.Spec.SecurityProtocol == "SASL" || r.Spec.SecurityProtocol == "SASL_SSL" {
		if !r.usernameDefined() {
			return field.Invalid(spec.Child("username"), "username", "must be defined for SASL_SSL authentication")
		}
		if !r.passwordDefined() {
			return field.Invalid(spec.Child("password"), "password", "must be defined for SASL_SSL authentication")
		}
	}
	return nil
}

func (r *KafkaConnection) brokerDefined() bool {
	if len(r.Spec.Brokers) > 0 {
		return true
	}
	return false
}

func (r *KafkaConnection) usernameDefined() bool {
	if len(r.Spec.Username) > 0 {
		return true
	}
	return false
}

func (r *KafkaConnection) passwordDefined() bool {
	if len(r.Spec.Password) > 0 {
		return true
	}
	return false
}
