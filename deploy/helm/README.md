<!--
SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
SPDX-License-Identifier: Apache-2.0

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
-->

# Manual Helm Deployment

This directory contains Helm charts for manually deploying Dynamo inference graphs to Kubernetes.
This approach allows you to install Dynamo directly using a DynamoGraphDeploymentCRD values file, which is useful for quick deployments or testing specific configurations.

### Prerequisites

- Helm 3.0+
- Kubernetes 1.16+
- ETCD v3.5+ (without auth)
- NATS v2.10+ (with jetstream enabled)

### Basic Installation

```bash
helm upgrade --install dynamo-graph ./deploy/helm/chart -n dynamo-cloud -f ./examples/vllm/deploy/agg.yaml
```

### Customizable Properties

You can override the default configuration by setting the following properties:

```bash
helm upgrade --install dynamo-graph ./deploy/helm/chart -n dynamo-cloud \
  -f ./examples/vllm/deploy/agg.yaml \
  --set "imagePullSecrets[0].name=docker-secret-1" \
  --set etcdAddr="my-etcd-service:2379" \
  --set natsAddr="nats://my-nats-service:4222"
```

#### Available Properties

| Property | Description | Example |
|----------|-------------|---------|
| `imagePullSecrets` | Array of image pull secrets for accessing private registries | `imagePullSecrets[0].name=docker-secret-1` |
| `etcdAddr` | Address of the etcd service | `dynamo-platform-etcd:2379` |
| `natsAddr` | Address of the NATS messaging service | `nats://dynamo-platform-nats:4222` |



