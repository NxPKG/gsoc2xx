package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Authentication struct {
	// +kubebuilder:validation:Optional
	ServiceAccount ServiceAccountDetails `json:"serviceAccount"`
	// +kubebuilder:validation:Optional
	ServiceToken ServiceTokenDetails `json:"serviceToken"`
}

type ServiceTokenDetails struct {
	// +kubebuilder:validation:Required
	ServiceTokenSecretReference KubeSecretReference `json:"serviceTokenSecretReference"`

	// +kubebuilder:validation:Required
	SecretsScope SecretScopeInWorkspace `json:"secretsScope"`
}

type ServiceAccountDetails struct {
	ServiceAccountSecretReference KubeSecretReference `json:"serviceAccountSecretReference"`
	ProjectId                     string              `json:"projectId"`
	EnvironmentName               string              `json:"environmentName"`
}

type SecretScopeInWorkspace struct {
	// +kubebuilder:validation:Required
	SecretsPath string `json:"secretsPath"`

	// +kubebuilder:validation:Required
	EnvSlug string `json:"envSlug"`
}

type KubeSecretReference struct {
	// The name of the Kubernetes Secret
	// +kubebuilder:validation:Required
	SecretName string `json:"secretName"`

	// The name space where the Kubernetes Secret is located
	// +kubebuilder:validation:Required
	SecretNamespace string `json:"secretNamespace"`
}

// Gsoc2SecretSpec defines the desired state of Gsoc2Secret
type Gsoc2SecretSpec struct {
	// +kubebuilder:validation:Optional
	TokenSecretReference KubeSecretReference `json:"tokenSecretReference"`

	// +kubebuilder:validation:Optional
	Authentication Authentication `json:"authentication"`

	// +kubebuilder:validation:Required
	ManagedSecretReference KubeSecretReference `json:"managedSecretReference"`

	// +kubebuilder:default:=60
	ResyncInterval int `json:"resyncInterval"`

	// Gsoc2 host to pull secrets from
	// +kubebuilder:validation:Optional
	HostAPI string `json:"hostAPI"`
}

// Gsoc2SecretStatus defines the observed state of Gsoc2Secret
type Gsoc2SecretStatus struct {
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Gsoc2Secret is the Schema for the gsoc2secrets API
type Gsoc2Secret struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   Gsoc2SecretSpec   `json:"spec,omitempty"`
	Status Gsoc2SecretStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// Gsoc2SecretList contains a list of Gsoc2Secret
type Gsoc2SecretList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Gsoc2Secret `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Gsoc2Secret{}, &Gsoc2SecretList{})
}
