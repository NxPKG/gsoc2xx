# Gsoc2 Helm Chart

This is the Gsoc2 Secrets Operator Helm chart. Find the integration documentation [here](https://gsoc2.com/docs/integrations/platforms/kubernetes)

## Installation

To install the chart, run the following :

```sh
# Add the Gsoc2 repository
helm repo add gsoc2 'https://dl.cloudsmith.io/public/gsoc2/helm-charts/helm/charts/' && helm repo update

# Install Gsoc2 Secrets Operator (with default values)
helm upgrade --install --atomic \
  -n gsoc2-dev --create-namespace \
  gsoc2-secrets-operator gsoc2/secrets-operator

# Install Gsoc2 Secrets Operator (with custom inline values, replace with your own values)
helm upgrade --install --atomic \
  -n gsoc2-dev --create-namespace \
  --set controllerManager.replicas=3 \
  gsoc2-secrets-operator gsoc2/secrets-operator

# Install Gsoc2 Secrets Operator (with custom values file, replace with your own values file)
helm upgrade --install --atomic \
  -n gsoc2-dev --create-namespace \
  -f custom-values.yaml \
  gsoc2-secrets-operator gsoc2/secrets-operator
```

## Synchronization

To sync your secrets from Gsoc2 (or from your own instance), create the below resources :

```sh
# Create the tokenSecretReference (replace with your own token)
kubectl create secret generic gsoc2-example-service-token \
  --from-literal=gsoc2Token="<gsoc2-token-here>"

# Create the Gsoc2Secret
cat <<EOF | kubectl apply -f -
apiVersion: secrets.gsoc2.com/v1alpha1
kind: Gsoc2Secret
metadata:
  # Name of of this Gsoc2Secret resource
  name: gsoc2secret-example
spec:
  # The host that should be used to pull secrets from. The default value is https://app.gsoc2.com/api.
  hostAPI: https://app.gsoc2.com/api

  # The Kubernetes secret the stores the Gsoc2 token
  tokenSecretReference:
    # Kubernetes secret name
    secretName: gsoc2-example-service-token
    # The secret namespace
    secretNamespace: default

  # The Kubernetes secret that Gsoc2 Operator will create and populate with secrets from the above project
  managedSecretReference:
    # The name of managed Kubernetes secret that should be created
    secretName: gsoc2-managed-secret
    # The namespace the managed secret should be installed in
    secretNamespace: default
EOF
```

### Managed secrets

#### Methods

To use the above created manage secrets, you can use the below methods :
- `env`
- `envFrom`
- `volumes`

Check the [docs](https://gsoc2.com/docs/integrations/platforms/kubernetes#using-managed-secret-in-your-deployment) to learn more about their implementation within your k8s resources

#### Auto-reload

And if you want to [auto-reload](https://gsoc2.com/docs/integrations/platforms/kubernetes#auto-redeployment) your deployments, add this annotation where the managed secret is consumed :

```yaml
annotations:
  secrets.gsoc2.com/auto-reload: "true"
```

## Parameters

*Coming soon*

## Local development

*Coming soon*

## Upgrading

### 0.1.2

Latest stable version, no breaking changes