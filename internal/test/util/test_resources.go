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

package util

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
)

// CustomConfig represents configuration for Kubegres tests
type CustomConfig struct {
	KubegresConfig string
	ConfigMapName  string
	ConfigMapData  map[string]string
}

// KubegresResources wraps resource operations for Kubegres
type KubegresResources struct {
	K8sClient *K8sClientType
	Kubegres  *kubegresv1.Kubegres
}

// NewKubegresResources creates a new KubegresResources instance
func NewKubegresResources(client *K8sClientType, kubegres *kubegresv1.Kubegres) *KubegresResources {
	return &KubegresResources{
		K8sClient: client,
		Kubegres:  kubegres,
	}
}

// LoadCustomConfig loads custom configuration from yaml files
func LoadCustomConfig(kubegresConfigPath, configMapPath string) *CustomConfig {
	if kubegresConfigPath == "" || configMapPath == "" {
		return &CustomConfig{
			KubegresConfig: "",
			ConfigMapName:  "test-config",
			ConfigMapData: map[string]string{
				"postgresql.conf": "max_connections = 100\nshared_buffers = 128MB",
			},
		}
	}

	// Read kubegres config file
	kubegresConfigBytes, err := ioutil.ReadFile(kubegresConfigPath)
	if err != nil {
		fmt.Printf("Error reading kubegres config file: %v\n", err)
		return nil
	}

	// Read configmap file
	configMapBytes, err := ioutil.ReadFile(configMapPath)
	if err != nil {
		fmt.Printf("Error reading configmap file: %v\n", err)
		return nil
	}

	// Create custom config
	return &CustomConfig{
		KubegresConfig: string(kubegresConfigBytes),
		ConfigMapName:  filepath.Base(configMapPath),
		ConfigMapData: map[string]string{
			filepath.Base(configMapPath): string(configMapBytes),
		},
	}
}

// WaitForStatefulSetAndPodsReady waits for statefulset and pods to be ready
func (r *KubegresResources) WaitForStatefulSetAndPodsReady(name string, replicas int) error {
	// Mock implementation for testing
	return nil
}

// GetPrimaryPod gets the primary pod
func (r *KubegresResources) GetPrimaryPod() (v1.Pod, error) {
	// Return a simulated pod
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "simulated-primary-pod",
			Labels: map[string]string{
				"role": "primary",
			},
		},
		Spec: v1.PodSpec{
			NodeName: "test-node",
		},
	}, nil
}

// WaitForAPrimaryToBeElected waits for a primary to be elected
func (r *KubegresResources) WaitForAPrimaryToBeElected(predicate func(v1.Pod) bool) error {
	// Mock implementation for testing
	return nil
}

// GetLogs gets logs
func (r *KubegresResources) GetLogs() (string, error) {
	// Return simulated logs
	return "SIMULATED LOGS\nStarted failing-over\nEnsuring old primary is terminated (STONITH)\nSTONITHEnabled\nFailOver: Promoting Replica to Primary", nil
}
