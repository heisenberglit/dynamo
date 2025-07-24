/*
 * SPDX-FileCopyrightText: Copyright (c) 2025 NVIDIA CORPORATION & AFFILIATES. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"github.com/ai-dynamo/dynamo/deploy/cloud/operator/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func constructPVC(crd metav1.Object, pvcConfig v1alpha1.PVC) *corev1.PersistentVolumeClaim {
	storageClassName := pvcConfig.StorageClass
	return &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      getPvcName(crd, pvcConfig.Name),
			Namespace: crd.GetNamespace(),
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{pvcConfig.VolumeAccessMode},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: pvcConfig.Size,
				},
			},
			StorageClassName: &storageClassName,
		},
	}
}

func getPvcName(crd metav1.Object, defaultName *string) string {
	if defaultName != nil {
		return *defaultName
	}
	return crd.GetName()
}

type dockerSecretRetriever interface {
	// returns a list of secret names associated with the docker registry
	GetSecrets(namespace, registry string) ([]string, error)
}
