/*
Copyright 2023 Reactive Tech Limited.
"Reactive Tech Limited" is a company located in England, United Kingdom.
https://www.reactive-tech.io

Lead Developer: Alex Arica

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resourceConfigs

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/test/util"
)

// TestResources represents resources for testing
type TestResources struct {
	Kubegres   kubegresv1.Kubegres
	ConfigMaps []v1.ConfigMap
}

// CreateBasicKubegres creates a basic Kubegres configuration for testing
func CreateBasicKubegres(name string, namespace string, replicas int32) kubegresv1.Kubegres {
	return kubegresv1.Kubegres{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: kubegresv1.KubegresSpec{
			Replicas: &replicas,
			Image:    "postgres:13.4",
		},
	}
}

// CreateKubegresWithSTONITH creates a Kubegres configuration with STONITH enabled
func CreateKubegresWithSTONITH(name string, namespace string, replicas int32) kubegresv1.Kubegres {
	enableSTONITH := true
	kubegres := CreateBasicKubegres(name, namespace, replicas)
	kubegres.Spec.Failover = kubegresv1.KubegresFailover{
		EnableSTONITH: &enableSTONITH,
	}
	return kubegres
}

// CreateKubegresWithBackup creates a Kubegres configuration with backup settings
func CreateKubegresWithBackup(name string, namespace string, replicas int32, schedule string, volumeMount string, pvcName string) kubegresv1.Kubegres {
	kubegres := CreateBasicKubegres(name, namespace, replicas)
	kubegres.Spec.Backup = kubegresv1.KubegresBackUp{
		Schedule:    schedule,
		VolumeMount: volumeMount,
		PvcName:     pvcName,
	}
	return kubegres
}

// CreateTestEnvironment creates the test resources and applies them to the cluster
func CreateTestEnvironment(client *util.K8sClientType, resources TestResources) error {
	// Create ConfigMaps first
	for _, cm := range resources.ConfigMaps {
		if err := client.Client.Create(client.Ctx, &cm); err != nil {
			return err
		}
	}

	// Create Kubegres instance
	if err := client.Client.Create(client.Ctx, &resources.Kubegres); err != nil {
		return err
	}

	return nil
}

// CleanupTestEnvironment cleans up the test resources
func CleanupTestEnvironment(client *util.K8sClientType, resources TestResources) error {
	// Delete Kubegres instance
	if err := client.Client.Delete(client.Ctx, &resources.Kubegres); err != nil {
		return err
	}

	// Delete ConfigMaps
	for _, cm := range resources.ConfigMaps {
		if err := client.Client.Delete(client.Ctx, &cm); err != nil {
			return err
		}
	}

	return nil
}
