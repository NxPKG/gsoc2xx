package controllers

import (
	"context"
	"fmt"

	"github.com/Gsoc2/gsoc2/k8-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (r *Gsoc2SecretReconciler) SetReadyToSyncSecretsConditions(ctx context.Context, gsoc2Secret *v1alpha1.Gsoc2Secret, errorToConditionOn error) error {
	if gsoc2Secret.Status.Conditions == nil {
		gsoc2Secret.Status.Conditions = []metav1.Condition{}
	}

	if errorToConditionOn != nil {
		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/ReadyToSyncSecrets",
			Status:  metav1.ConditionFalse,
			Reason:  "Error",
			Message: "Failed to sync secrets. This can be caused by invalid service token or an invalid API host that is set. Check operator logs for more info",
		})

		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/AutoRedeployReady",
			Status:  metav1.ConditionFalse,
			Reason:  "Stopped",
			Message: "Auto redeployment has been stopped because the operator failed to sync secrets",
		})
	} else {
		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/ReadyToSyncSecrets",
			Status:  metav1.ConditionTrue,
			Reason:  "OK",
			Message: "Gsoc2 controller has started syncing your secrets",
		})
	}

	return r.Client.Status().Update(ctx, gsoc2Secret)
}

func (r *Gsoc2SecretReconciler) SetGsoc2TokenLoadCondition(ctx context.Context, gsoc2Secret *v1alpha1.Gsoc2Secret, errorToConditionOn error) {
	if gsoc2Secret.Status.Conditions == nil {
		gsoc2Secret.Status.Conditions = []metav1.Condition{}
	}

	if errorToConditionOn == nil {
		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/LoadedGsoc2Token",
			Status:  metav1.ConditionTrue,
			Reason:  "OK",
			Message: "Gsoc2 controller has located the Gsoc2 token in provided Kubernetes secret",
		})
	} else {
		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/LoadedGsoc2Token",
			Status:  metav1.ConditionFalse,
			Reason:  "Error",
			Message: fmt.Sprintf("Failed to load Gsoc2 Token from the provided Kubernetes secret because: %v", errorToConditionOn),
		})
	}

	err := r.Client.Status().Update(ctx, gsoc2Secret)
	if err != nil {
		fmt.Println("Could not set condition for LoadedGsoc2Token")
	}
}

func (r *Gsoc2SecretReconciler) SetGsoc2AutoRedeploymentReady(ctx context.Context, gsoc2Secret *v1alpha1.Gsoc2Secret, numDeployments int, errorToConditionOn error) {
	if gsoc2Secret.Status.Conditions == nil {
		gsoc2Secret.Status.Conditions = []metav1.Condition{}
	}

	if errorToConditionOn == nil {
		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/AutoRedeployReady",
			Status:  metav1.ConditionTrue,
			Reason:  "OK",
			Message: fmt.Sprintf("Gsoc2 has found %v deployments which are ready to be auto redeployed when secrets change", numDeployments),
		})
	} else {
		meta.SetStatusCondition(&gsoc2Secret.Status.Conditions, metav1.Condition{
			Type:    "secrets.gsoc2.com/AutoRedeployReady",
			Status:  metav1.ConditionFalse,
			Reason:  "Error",
			Message: fmt.Sprintf("Failed reconcile deployments because: %v", errorToConditionOn),
		})
	}

	err := r.Client.Status().Update(ctx, gsoc2Secret)
	if err != nil {
		fmt.Println("Could not set condition for AutoRedeployReady")
	}
}
