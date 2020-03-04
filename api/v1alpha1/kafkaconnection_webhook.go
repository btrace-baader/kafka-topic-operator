package v1alpha1

import (
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	validationutils "k8s.io/apimachinery/pkg/util/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
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
// +kubebuilder:webhook:verbs=create;update,path=/validate-kafka-btrace-com-v1alpha1-kafkaconnection,mutating=false,failurePolicy=fail,groups=kafka.btrace.com,resources=kafkaconnections,versions=v1alpha1,name=vkafkaconnection.kb.io

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
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KafkaConnection) ValidateDelete() error {
	kafkaconnectionlog.Info("validate delete", "name", r.Name)

	if err := r.validateKafkaConnection(); err != nil {
		return err
	}
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
		return field.Invalid(spec.Child("broker"), r.Name, "must be defined")
	}
	if r.Spec.AuthMethod == "SASL" || r.Spec.AuthMethod == "SASL_SSL" {
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
	if len(r.Spec.Broker) > 0 {
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
