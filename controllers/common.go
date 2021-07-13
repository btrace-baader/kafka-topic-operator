package controllers

import (
	"time"

	"github.com/go-logr/logr"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

func ignoreNotFound(err error) error {
	if apierrs.IsNotFound(err) {
		return nil
	}
	return err
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}

// reconciled without timeout, send empty response
func reconcile() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

// reconcileWithTimeout requeues after specified amount of time
// to automatically reconcile external resources.
func requeueWithTimeout(t int32) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: time.Second * time.Duration(t),
		Requeue: true,
	}, nil
}

// requeueWithError requeues immediately and returns error
func requeueWithError(log logr.Logger, msg string, err error) (ctrl.Result, error) {
	log.Error(err, msg)
	return ctrl.Result{}, err
}
