<!--
SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES.
All rights reserved.
SPDX-License-Identifier: Apache-2.0
-->

# Dynamo Support Matrix

This document provides the support matrix for Dynamo, including hardware, software and build instructions.

## Hardware Compatibility

| **CPU Architecture** | **Status**   |
| :------------------- | :----------- |
| **x86_64**           | Supported    |
| **ARM64**            | Experimental |

> [!Warning]
> While **x86_64** architecture is supported on systems with a minimum of 32 GB RAM and at least 4 CPU cores,
> the **ARM64** support is experimental and may have limitations.

### GPU Compatibility

If you are using a **GPU**, the following GPU models and architectures are supported:

| **GPU Architecture**                 | **Status** |
| :----------------------------------- | :--------- |
| **NVIDIA Blackwell Architecture**    | Supported  |
| **NVIDIA Hopper Architecture**       | Supported  |
| **NVIDIA Ada Lovelace Architecture** | Supported  |
| **NVIDIA Ampere Architecture**       | Supported  |


## Platform Architecture Compatibility

**Dynamo** is compatible with the following platforms:

| **Operating System** | **Version** | **Architecture** | **Status**   |
| :------------------- | :---------- | :--------------- | :----------- |
| **Ubuntu**           | 22.04       | x86_64           | Supported    |
| **Ubuntu**           | 24.04       | x86_64           | Supported    |
| **Ubuntu**           | 24.04       | ARM64            | Experimental |
| **CentOS Stream**    | 9           | x86_64           | Experimental |

> [!Note]
> For **Linux**, the **ARM64** support is experimental and may have limitations.
> Wheels are built using a manylinux_2_28-compatible environment and they have been validated on CentOS 9 and Ubuntu (22.04, 24.04).
>
> Compatibility with other Linux distributions is expected but has not been officially verified yet.

> [!Caution]
> KV Block Manager is supported only with Python 3.12. Python 3.12 support is currently limited to Ubuntu 24.04.


## Software Compatibility

### Runtime Dependency

| **Python Package** | **Version**   | glibc version                        | CUDA Version |
| :----------------- | :------------ | :----------------------------------- | :----------- |
| ai-dynamo          | 0.3.2         | >=2.28                               |              |
| ai-dynamo-runtime  | 0.3.2         | >=2.28 (Python 3.12 has known issues)|              |
| NIXL               | 0.4.0         | >=2.27                               | >=11.8       |

### Build Dependency

| **Build Dependency** | **Version**                                                                      |
| :------------------- | :------------------------------------------------------------------------------- |
| **Base Container**   | [25.03](https://catalog.ngc.nvidia.com/orgs/nvidia/containers/cuda-dl-base/tags) |
| **TensorRT-LLM**     | 1.0.0rc²                                                                         |
| **NIXL**             | 0.4.0                                                                            |

> [!Important]
> ² Specific versions of TensorRT-LLM supported by Dynamo are subject to change.

## Build Support

**Dynamo** currently provides build support in the following ways:

- **Wheels**: Pre-built Python wheels are only available for **x86_64 Linux**.
   No wheels are available for other platforms at this time.

- **Container Images**: We distribute only the source code for container images, **x86_64 Linux** and **ARM64** are supported for these.
   Users must build the container image from source if they require it.

Once you've confirmed that your platform and architecture are compatible, you can install **Dynamo** by following the instructions in the [Quick Start Guide](https://github.com/ai-dynamo/dynamo/blob/main/README.md#installation).
