package controllers

import (
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/Gsoc2/gsoc2/k8-operator/api/v1alpha1"
	secretsv1alpha1 "github.com/Gsoc2/gsoc2/k8-operator/api/v1alpha1"
	"github.com/Gsoc2/gsoc2/k8-operator/packages/api"
)

// Gsoc2SecretReconciler reconciles a Gsoc2Secret object
type Gsoc2SecretReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=secrets.gsoc2.com,resources=gsoc2secrets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=secrets.gsoc2.com,resources=gsoc2secrets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=secrets.gsoc2.com,resources=gsoc2secrets/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=secrets,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;delete
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=list;watch;get;update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *Gsoc2SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var gsoc2SecretCR v1alpha1.Gsoc2Secret
	requeueTime := time.Minute // seconds

	err := r.Get(ctx, req.NamespacedName, &gsoc2SecretCR)
	if err != nil {
		if errors.IsNotFound(err) {
			fmt.Printf("Gsoc2 Secret CRD not found [err=%v]", err)
			return ctrl.Result{
				Requeue: false,
			}, nil
		} else {
			fmt.Printf("Unable to fetch Gsoc2 Secret CRD from cluster because [err=%v]", err)
			return ctrl.Result{
				RequeueAfter: requeueTime,
			}, nil
		}
	}

	if gsoc2SecretCR.Spec.ResyncInterval != 0 {
		requeueTime = time.Second * time.Duration(gsoc2SecretCR.Spec.ResyncInterval)
		fmt.Println("Manual re-sync interval set", "requeueAfter", requeueTime)
	}

	fmt.Println("Requeue duration set", "requeueAfter", requeueTime)

	// Check if the resource is already marked for deletion
	if gsoc2SecretCR.GetDeletionTimestamp() != nil {
		return ctrl.Result{
			Requeue: false,
		}, nil
	}

	// Get modified/default config
	gsoc2Config, err := r.GetGsoc2ConfigMap(ctx)
	if err != nil {
		fmt.Printf("unable to fetch gsoc2-config [err=%s]. Will requeue after [requeueTime=%v]\n", err, requeueTime)
		return ctrl.Result{
			RequeueAfter: requeueTime,
		}, nil
	}

	if gsoc2SecretCR.Spec.HostAPI == "" {
		api.API_HOST_URL = gsoc2Config["hostAPI"]
	} else {
		api.API_HOST_URL = gsoc2SecretCR.Spec.HostAPI
	}

	err = r.ReconcileGsoc2Secret(ctx, gsoc2SecretCR)
	r.SetReadyToSyncSecretsConditions(ctx, &gsoc2SecretCR, err)

	if err != nil {
		fmt.Printf("unable to reconcile Gsoc2 Secret because [err=%v]. Will requeue after [requeueTime=%v]\n", err, requeueTime)
		return ctrl.Result{
			RequeueAfter: requeueTime,
		}, nil
	}

	numDeployments, err := r.ReconcileDeploymentsWithManagedSecrets(ctx, gsoc2SecretCR)
	r.SetGsoc2AutoRedeploymentReady(ctx, &gsoc2SecretCR, numDeployments, err)
	if err != nil {
		fmt.Printf("unable to reconcile auto redeployment because [err=%v]", err)
		return ctrl.Result{
			RequeueAfter: requeueTime,
		}, nil
	}

	// Sync again after the specified time
	fmt.Printf("Operator will requeue after [%v] \n", requeueTime)
	return ctrl.Result{
		RequeueAfter: requeueTime,
	}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *Gsoc2SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&secretsv1alpha1.Gsoc2Secret{}).
		Complete(r)
}
