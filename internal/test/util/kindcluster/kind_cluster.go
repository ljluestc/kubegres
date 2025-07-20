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

package kindcluster

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// KindCluster represents a Kind Kubernetes cluster for testing
type KindCluster struct {
	Name           string
	KubeConfigPath string
}

// NewKindCluster creates a new Kind cluster instance
func NewKindCluster(name string) *KindCluster {
	return &KindCluster{
		Name:           name,
		KubeConfigPath: "",
	}
}

// Create creates a new Kind cluster
func (k *KindCluster) Create() error {
	// Check if cluster already exists
	if k.Exists() {
		return fmt.Errorf("kind cluster %s already exists", k.Name)
	}

	// Create temp file for kubeconfig
	tmpFile, err := os.CreateTemp("", "kubeconfig-*.yaml")
	if err != nil {
		return fmt.Errorf("failed to create temp file for kubeconfig: %w", err)
	}
	k.KubeConfigPath = tmpFile.Name()
	tmpFile.Close()

	// Create the cluster
	cmd := exec.Command("kind", "create", "cluster",
		"--name", k.Name,
		"--kubeconfig", k.KubeConfigPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create kind cluster: %w, output: %s", err, string(output))
	}

	return nil
}

// Delete deletes the Kind cluster
func (k *KindCluster) Delete() error {
	if !k.Exists() {
		return nil // Already deleted
	}

	cmd := exec.Command("kind", "delete", "cluster", "--name", k.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete kind cluster: %w, output: %s", err, string(output))
	}

	// Clean up kubeconfig file
	if k.KubeConfigPath != "" {
		os.Remove(k.KubeConfigPath)
	}

	return nil
}

// Exists checks if the Kind cluster exists
func (k *KindCluster) Exists() bool {
	cmd := exec.Command("kind", "get", "clusters")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	clusters := strings.Split(string(output), "\n")
	for _, cluster := range clusters {
		if strings.TrimSpace(cluster) == k.Name {
			return true
		}
	}

	return false
}

// LoadImage loads a Docker image into the Kind cluster
func (k *KindCluster) LoadImage(imageName string) error {
	if !k.Exists() {
		return errors.New("kind cluster does not exist")
	}

	cmd := exec.Command("kind", "load", "docker-image", imageName, "--name", k.Name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to load image into kind cluster: %w, output: %s", err, string(output))
	}

	return nil
}

// ApplyManifest applies a Kubernetes manifest to the cluster
func (k *KindCluster) ApplyManifest(manifestPath string) error {
	if !k.Exists() {
		return errors.New("kind cluster does not exist")
	}

	cmd := exec.Command("kubectl", "apply", "-f", manifestPath, "--kubeconfig", k.KubeConfigPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to apply manifest: %w, output: %s", err, string(output))
	}

	return nil
}

// GetClusterInfo returns information about the Kind cluster
func (k *KindCluster) GetClusterInfo() (string, error) {
	if !k.Exists() {
		return "", errors.New("kind cluster does not exist")
	}

	cmd := exec.Command("kubectl", "cluster-info", "--context", fmt.Sprintf("kind-%s", k.Name), "--kubeconfig", k.KubeConfigPath)
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get cluster info: %w", err)
	}

	return string(output), nil
}
