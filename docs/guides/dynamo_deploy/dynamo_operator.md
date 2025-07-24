# Working with Dynamo Kubernetes Operator

## Overview

Dynamo operator is a Kubernetes operator that simplifies the deployment, configuration, and lifecycle management of DynamoGraphs. It automates the reconciliation of custom resources to ensure your desired state is always achieved. This operator is ideal for users who want to manage complex deployments using declarative YAML definitions and Kubernetes-native tooling.

## Architecture

- **Operator Deployment:**
  Deployed as a Kubernetes `Deployment` in a specific namespace.

- **Controllers:**
  - `DynamoGraphDeploymentController`: Watches `DynamoGraphDeployment` CRs and orchestrates graph deployments.
  - `DynamoComponentDeploymentController`: Watches `DynamoComponentDeployment` CRs and handles individual component deployments.

- **Workflow:**
  1. A custom resource is created by the user or API server.
  2. The corresponding controller detects the change and runs reconciliation.
  3. Kubernetes resources (Deployments, Services, etc.) are created or updated to match the CR spec.
  4. Status fields are updated to reflect the current state.



## Custom Resource Definitions (CRDs)

### CRD: `DynamoGraphDeployment`


| Field            | Type   | Description                                                                                                                                          | Required | Default |
|------------------|--------|------------------------------------------------------------------------------------------------------------------------------------------------------|----------|---------|
| `services`       | map    | Map of service names to runtime configurations. This allows the user to override the service configuration defined in the DynamoComponentDeployment. | Yes      |         |
| `envs`           | list   | list of global environment variables.                                                                                                                | No       |         |


**API Version:** `nvidia.com/v1alpha1`
**Scope:** Namespaced

#### Example
```yaml
apiVersion: nvidia.com/v1alpha1
kind: DynamoGraphDeployment
metadata:
  name: disagg
spec:
  envs:
  - name: GLOBAL_ENV_VAR
    value: some_global_value
  services:
    Frontend:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
    Processor:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
    VllmWorker:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
    PrefillWorker:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
    Router:
      replicas: 0
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
```

### CRD: `DynamoComponentDeployment`

| Field              | Type     | Description                                                   | Required | Default |
|--------------------|----------|---------------------------------------------------------------|----------|---------|
| `dynamoNamespace` | string   | Namespace of the DynamoComponent                               | Yes      |         |
| `serviceName`     | string   | Logical name of the service being deployed                     | Yes      |         |
| `envs`            | array    | Environment variables for runtime                              | No       | `[]`    |
| `annotations`     | map      | Additional metadata annotations for the pod                    | No       |         |
| `labels`          | map      | Custom labels applied to the deployment and pod                | No       |         |
| `resources`       | object   | Resource limits and requests (CPU, memory, GPU)                | No       |         |
| `autoscaling`     | object   | Autoscaling rules for the deployment                           | No       |         |
| `envFromSecret`   | string   | Reference to a secret for injecting env vars                   | No       |         |
| `pvc`             | object   | Persistent volume claim configuration                          | No       |         |
| `ingress`         | object   | Ingress configuration for exposing the service                 | No       |         |
| `extraPodMetadata`| object   | Additional labels and annotations for the pod                  | No       |         |
| `extraPodSpec`    | object   | Custom PodSpec fields to merge into the generated pod          | No       |         |
| `livenessProbe`   | object   | Kubernetes liveness probe                                      | No       |         |
| `readinessProbe`  | object   | Kubernetes readiness probe                                     | No       |         |
| `replicas`        | int      | Number of replicas to run                                      | No       | `1`     |

**API Version:** `nvidia.com/v1alpha1`
**Scope:** Namespaced

#### Example
```yaml
apiVersion: nvidia.com/v1alpha1
kind: DynamoComponentDeployment
metadata:
  name: test-41fa991-vllmworker
spec:
  dynamoNamespace: dynamo
  envs:
    - name: DYN_DEPLOYMENT_CONFIG
      value: '<long JSON config>'
  externalServices:
    PrefillWorker:
      deploymentSelectorKey: dynamo
      deploymentSelectorValue: PrefillWorker/dynamo
  resources:
    limits:
      cpu: "10"
      gpu: "1"
      memory: 20Gi
    requests:
      cpu: "500m"
      gpu: "1"
      memory: 20Gi
  serviceName: Frontend
```

## Installation

[See installation steps](dynamo_cloud.md#overview)


## GitOps Deployment with FluxCD

This section describes how to use FluxCD for GitOps-based deployment of Dynamo inference graphs. GitOps enables you to manage your Dynamo deployments declaratively using Git as the source of truth. We'll use the [aggregated vLLM example](https://github.com/ai-dynamo/dynamo/blob/main/examples/llm/README.md) to demonstrate the workflow.

### Prerequisites

- A Kubernetes cluster with [Dynamo Cloud](dynamo_cloud.md) installed
- [FluxCD](https://fluxcd.io/flux/installation/) installed in your cluster
- A Git repository to store your deployment configurations

### Workflow Overview

The GitOps workflow for Dynamo deployments consists of three main steps:

1. Build and push a pipeline to the Dynamo API store
2. Create and commit a DynamoGraphDeployment custom resource for initial deployment
3. Update the pipeline by building a new version and updating the CR for subsequent updates

### Step 1: Build and Push Pipeline

First, build and push your pipeline using the Dynamo CLI:

```bash
# Set your project root directory
export PROJECT_ROOT=$(pwd)

# Configure environment variables
export KUBE_NS=dynamo-cloud
export DYNAMO_CLOUD=http://localhost:8080  # If using port-forward
# OR
# export DYNAMO_CLOUD=https://dynamo-cloud.nvidia.com  # If using Ingress/VirtualService

# Build and push the service
cd $PROJECT_ROOT/examples/llm
DYNAMO_TAG=$(dynamo build --push graphs.agg:Frontend | grep "Successfully built" |  awk '{ print $NF }' | sed 's/\.$//')
```

The `--push` flag ensures the pipeline is pushed to the remote API store, making it available for deployment.

### Step 2: Create Initial Deployment

Create a new file in your Git repository (e.g., `deployments/llm-agg.yaml`) with the following content:

```yaml
apiVersion: nvidia.com/v1alpha1
kind: DynamoGraphDeployment
metadata:
  name: llm-agg
spec:
  services:
    Frontend:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
    Processor:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
    VllmWorker:
      replicas: 1
      envs:
      - name: SPECIFIC_ENV_VAR
        value: some_specific_value
      # Add PVC for model storage
      pvc:
        name: vllm-model-storage
        mountPath: /models
        size: 100Gi
```

Commit and push this file to your Git repository. FluxCD will detect the new CR and create the initial deployment in your cluster. The operator will:
- Create the specified PVCs
- Build container images for all components
- Deploy the services with the configured resources

### Step 3: Update Existing Deployment

To update your pipeline, just update the associated DynamoGraphDeployment CRD

The Dynamo operator will automatically reconcile it.

### Monitoring the Deployment

You can monitor the deployment status using:

```bash
# Check the DynamoGraphDeployment status
kubectl get dynamographdeployment llm-agg -n $KUBE_NS

# Check the component deployments
kubectl get dynamocomponentdeployment -n $KUBE_NS
```

---

## Deploying a Dynamo Pipeline using the Operator

[See deployment steps](operator_deployment.md)


## Reconciliation Logic

### DynamoGraphDeployment

- **Actions:**
  - Create a DynamoComponent CR to build the docker image
  - Create a DynamoComponentDeployment CR for each component defined in the Dynamo graph being deployed
- **Status Management:**
  - `.status.conditions`: Reflects readiness, failure, progress states
  - `.status.state`: overall state of the deployment, based on the state of the DynamoComponentDeployments

### DynamoComponentDeployment

- **Actions:**
  - Create a Deployment, Service, and Ingress for the service
- **Status Management:**
  - `.status.conditions`: Reflects readiness, failure, progress states

## Configuration


- **Environment Variables:**

| Name                                               | Description                          | Default                                                |
|----------------------------------------------------|--------------------------------------|--------------------------------------------------------|
| `LOG_LEVEL`                                        | Logging verbosity level              | `info`                                                 |
| `DYNAMO_SYSTEM_NAMESPACE`                          | System namespace                     | `dynamo`                                               |

- **Flags:**
  | Flag                  | Description                                | Default |
  |-----------------------|--------------------------------------------|---------|
  | `--natsAddr`          | Address of NATS server                     | ""      |
  | `--etcdAddr`          | Address of etcd server                     | ""      |



## Troubleshooting

| Symptom                | Possible Cause                | Solution                          |
|------------------------|-------------------------------|-----------------------------------|
| Resource not created   | RBAC missing                  | Ensure correct ClusterRole/Binding|
| Status not updated     | CRD schema mismatch           | Regenerate CRDs with kubebuilder  |
| Image build hangs      | Misconfigured DynamoComponent | Check image build logs            |


## Development

- **Code Structure:**

The operator is built using Kubebuilder and the operator-sdk, with the following structure:

- `controllers/`: Reconciliation logic
- `api/v1alpha1/`: CRD types
- `config/`: Manifests and Helm charts


## References

- [Kubernetes Operator Pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Operator SDK](https://sdk.operatorframework.io/)
- [Helm Best Practices for CRDs](https://helm.sh/docs/chart_best_practices/custom_resource_definitions/)
