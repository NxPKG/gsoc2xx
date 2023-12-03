package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/Gsoc2/gsoc2/k8-operator/api/v1alpha1"
	"github.com/Gsoc2/gsoc2/k8-operator/packages/api"
	"github.com/Gsoc2/gsoc2/k8-operator/packages/model"
	"github.com/Gsoc2/gsoc2/k8-operator/packages/util"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const SERVICE_ACCOUNT_ACCESS_KEY = "serviceAccountAccessKey"
const SERVICE_ACCOUNT_PUBLIC_KEY = "serviceAccountPublicKey"
const SERVICE_ACCOUNT_PRIVATE_KEY = "serviceAccountPrivateKey"

const GSOC2_TOKEN_SECRET_KEY_NAME = "gsoc2Token"
const SECRET_VERSION_ANNOTATION = "secrets.gsoc2.com/version" // used to set the version of secrets via Etag
const OPERATOR_SETTINGS_CONFIGMAP_NAME = "gsoc2-config"
const OPERATOR_SETTINGS_CONFIGMAP_NAMESPACE = "gsoc2-operator-system"
const GSOC2_DOMAIN = "https://app.gsoc2.com/api"

func (r *Gsoc2SecretReconciler) GetGsoc2ConfigMap(ctx context.Context) (configMap map[string]string, errToReturn error) {
	// default key values
	defaultConfigMapData := make(map[string]string)
	defaultConfigMapData["hostAPI"] = GSOC2_DOMAIN

	kubeConfigMap := &corev1.ConfigMap{}
	err := r.Client.Get(ctx, types.NamespacedName{
		Namespace: OPERATOR_SETTINGS_CONFIGMAP_NAMESPACE,
		Name:      OPERATOR_SETTINGS_CONFIGMAP_NAME,
	}, kubeConfigMap)

	if err != nil {
		if errors.IsNotFound(err) {
			kubeConfigMap = nil
		} else {
			return nil, fmt.Errorf("GetConfigMapByNamespacedName: unable to fetch config map in [namespacedName=%s] [err=%s]", OPERATOR_SETTINGS_CONFIGMAP_NAMESPACE, err)
		}
	}

	if kubeConfigMap == nil {
		return defaultConfigMapData, nil
	} else {
		for key, value := range defaultConfigMapData {
			_, exists := kubeConfigMap.Data[key]
			if !exists {
				kubeConfigMap.Data[key] = value
			}
		}

		return kubeConfigMap.Data, nil
	}
}

func (r *Gsoc2SecretReconciler) GetKubeSecretByNamespacedName(ctx context.Context, namespacedName types.NamespacedName) (*corev1.Secret, error) {
	kubeSecret := &corev1.Secret{}
	err := r.Client.Get(ctx, namespacedName, kubeSecret)
	if err != nil {
		kubeSecret = nil
	}

	return kubeSecret, err
}

func (r *Gsoc2SecretReconciler) GetGsoc2TokenFromKubeSecret(ctx context.Context, gsoc2Secret v1alpha1.Gsoc2Secret) (string, error) {
	// default to new secret ref structure
	secretName := gsoc2Secret.Spec.Authentication.ServiceToken.ServiceTokenSecretReference.SecretName
	secretNamespace := gsoc2Secret.Spec.Authentication.ServiceToken.ServiceTokenSecretReference.SecretNamespace

	// fall back to previous secret ref
	if secretName == "" {
		secretName = gsoc2Secret.Spec.TokenSecretReference.SecretName
	}

	if secretNamespace == "" {
		secretNamespace = gsoc2Secret.Spec.TokenSecretReference.SecretNamespace
	}

	tokenSecret, err := r.GetKubeSecretByNamespacedName(ctx, types.NamespacedName{
		Namespace: secretNamespace,
		Name:      secretName,
	})

	if errors.IsNotFound(err) {
		return "", nil
	}

	if err != nil {
		return "", fmt.Errorf("failed to read Gsoc2 token secret from secret named [%s] in namespace [%s]: with error [%w]", gsoc2Secret.Spec.TokenSecretReference.SecretName, gsoc2Secret.Spec.TokenSecretReference.SecretNamespace, err)
	}

	gsoc2ServiceToken := tokenSecret.Data[GSOC2_TOKEN_SECRET_KEY_NAME]

	return strings.Replace(string(gsoc2ServiceToken), " ", "", -1), nil
}

// Fetches service account credentials from a Kubernetes secret specified in the gsoc2Secret object, extracts the access key, public key, and private key from the secret, and returns them as a ServiceAccountCredentials object.
// If any keys are missing or an error occurs, returns an empty object or an error object, respectively.
func (r *Gsoc2SecretReconciler) GetGsoc2ServiceAccountCredentialsFromKubeSecret(ctx context.Context, gsoc2Secret v1alpha1.Gsoc2Secret) (serviceAccountDetails model.ServiceAccountDetails, err error) {
	serviceAccountCredsFromKubeSecret, err := r.GetKubeSecretByNamespacedName(ctx, types.NamespacedName{
		Namespace: gsoc2Secret.Spec.Authentication.ServiceAccount.ServiceAccountSecretReference.SecretNamespace,
		Name:      gsoc2Secret.Spec.Authentication.ServiceAccount.ServiceAccountSecretReference.SecretName,
	})

	if errors.IsNotFound(err) {
		return model.ServiceAccountDetails{}, nil
	}

	if err != nil {
		return model.ServiceAccountDetails{}, fmt.Errorf("something went wrong when fetching your service account credentials [err=%s]", err)
	}

	accessKeyFromSecret := serviceAccountCredsFromKubeSecret.Data[SERVICE_ACCOUNT_ACCESS_KEY]
	publicKeyFromSecret := serviceAccountCredsFromKubeSecret.Data[SERVICE_ACCOUNT_PUBLIC_KEY]
	privateKeyFromSecret := serviceAccountCredsFromKubeSecret.Data[SERVICE_ACCOUNT_PRIVATE_KEY]

	if accessKeyFromSecret == nil || publicKeyFromSecret == nil || privateKeyFromSecret == nil {
		return model.ServiceAccountDetails{}, nil
	}

	return model.ServiceAccountDetails{AccessKey: string(accessKeyFromSecret), PrivateKey: string(privateKeyFromSecret), PublicKey: string(publicKeyFromSecret)}, nil
}

func (r *Gsoc2SecretReconciler) CreateGsoc2ManagedKubeSecret(ctx context.Context, gsoc2Secret v1alpha1.Gsoc2Secret, secretsFromAPI []model.SingleEnvironmentVariable, encryptedSecretsResponse api.GetEncryptedSecretsV3Response) error {
	plainProcessedSecrets := make(map[string][]byte)
	for _, secret := range secretsFromAPI {
		plainProcessedSecrets[secret.Key] = []byte(secret.Value) // plain process
	}

	// create a new secret as specified by the managed secret spec of CRD
	newKubeSecretInstance := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gsoc2Secret.Spec.ManagedSecretReference.SecretName,
			Namespace: gsoc2Secret.Spec.ManagedSecretReference.SecretNamespace,
			Annotations: map[string]string{
				SECRET_VERSION_ANNOTATION: encryptedSecretsResponse.ETag,
			},
		},
		Type: "Opaque",
		Data: plainProcessedSecrets,
	}

	err := r.Client.Create(ctx, newKubeSecretInstance)
	if err != nil {
		return fmt.Errorf("unable to create the managed Kubernetes secret : %w", err)
	}

	fmt.Println("Successfully created a managed Kubernetes secret with your Gsoc2 secrets")
	return nil
}

func (r *Gsoc2SecretReconciler) UpdateGsoc2ManagedKubeSecret(ctx context.Context, managedKubeSecret corev1.Secret, secretsFromAPI []model.SingleEnvironmentVariable, encryptedSecretsResponse api.GetEncryptedSecretsV3Response) error {
	plainProcessedSecrets := make(map[string][]byte)
	for _, secret := range secretsFromAPI {
		plainProcessedSecrets[secret.Key] = []byte(secret.Value)
	}

	managedKubeSecret.Data = plainProcessedSecrets
	managedKubeSecret.ObjectMeta.Annotations = map[string]string{
		SECRET_VERSION_ANNOTATION: encryptedSecretsResponse.ETag,
	}

	err := r.Client.Update(ctx, &managedKubeSecret)
	if err != nil {
		return fmt.Errorf("unable to update Kubernetes secret because [%w]", err)
	}

	fmt.Println("successfully updated managed Kubernetes secret")
	return nil
}

func (r *Gsoc2SecretReconciler) ReconcileGsoc2Secret(ctx context.Context, gsoc2Secret v1alpha1.Gsoc2Secret) error {
	gsoc2Token, err := r.GetGsoc2TokenFromKubeSecret(ctx, gsoc2Secret)
	if err != nil {
		return fmt.Errorf("ReconcileGsoc2Secret: unable to get service token from kube secret [err=%s]", err)
	}

	serviceAccountCreds, err := r.GetGsoc2ServiceAccountCredentialsFromKubeSecret(ctx, gsoc2Secret)
	if err != nil {
		return fmt.Errorf("ReconcileGsoc2Secret: unable to get service account creds from kube secret [err=%s]", err)
	}

	r.SetGsoc2TokenLoadCondition(ctx, &gsoc2Secret, err)
	if err != nil {
		return fmt.Errorf("unable to load Gsoc2 Token from the specified Kubernetes secret with error [%w]", err)
	}

	// Look for managed secret by name and namespace
	managedKubeSecret, err := r.GetKubeSecretByNamespacedName(ctx, types.NamespacedName{
		Name:      gsoc2Secret.Spec.ManagedSecretReference.SecretName,
		Namespace: gsoc2Secret.Spec.ManagedSecretReference.SecretNamespace,
	})

	if err != nil && !errors.IsNotFound(err) {
		return fmt.Errorf("something went wrong when fetching the managed Kubernetes secret [%w]", err)
	}

	// Get exiting Etag if exists
	secretVersionBasedOnETag := ""
	if managedKubeSecret != nil {
		secretVersionBasedOnETag = managedKubeSecret.Annotations[SECRET_VERSION_ANNOTATION]
	}

	var plainTextSecretsFromApi []model.SingleEnvironmentVariable
	var fullEncryptedSecretsResponse api.GetEncryptedSecretsV3Response

	if serviceAccountCreds.AccessKey != "" || serviceAccountCreds.PrivateKey != "" || serviceAccountCreds.PublicKey != "" {
		plainTextSecretsFromApi, fullEncryptedSecretsResponse, err = util.GetPlainTextSecretsViaServiceAccount(serviceAccountCreds, gsoc2Secret.Spec.Authentication.ServiceAccount.ProjectId, gsoc2Secret.Spec.Authentication.ServiceAccount.EnvironmentName, secretVersionBasedOnETag)
		if err != nil {
			return fmt.Errorf("\nfailed to get secrets because [err=%v]", err)
		}

		fmt.Println("ReconcileGsoc2Secret: Fetched secrets via service account")

	} else if gsoc2Token != "" {
		envSlug := gsoc2Secret.Spec.Authentication.ServiceToken.SecretsScope.EnvSlug
		secretsPath := gsoc2Secret.Spec.Authentication.ServiceToken.SecretsScope.SecretsPath

		plainTextSecretsFromApi, fullEncryptedSecretsResponse, err = util.GetPlainTextSecretsViaServiceToken(gsoc2Token, secretVersionBasedOnETag, envSlug, secretsPath)
		if err != nil {
			return fmt.Errorf("\nfailed to get secrets because [err=%v]", err)
		}

		fmt.Println("ReconcileGsoc2Secret: Fetched secrets via service token")

	} else {
		return fmt.Errorf("no authentication method provided. You must provide either a valid service token or a service account details to fetch secrets")
	}

	if !fullEncryptedSecretsResponse.Modified {
		fmt.Println("No secrets modified so reconcile not needed", "Etag:", fullEncryptedSecretsResponse.ETag, "Modified:", fullEncryptedSecretsResponse.Modified)
		return nil
	}

	if managedKubeSecret == nil {
		return r.CreateGsoc2ManagedKubeSecret(ctx, gsoc2Secret, plainTextSecretsFromApi, fullEncryptedSecretsResponse)
	} else {
		return r.UpdateGsoc2ManagedKubeSecret(ctx, *managedKubeSecret, plainTextSecretsFromApi, fullEncryptedSecretsResponse)
	}

}
