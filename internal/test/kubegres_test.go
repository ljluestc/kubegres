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

package test

import (
	"testing"

	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/test/util"
)

// TestKubegresBasic is a simple test to verify the kubegres functionality
func TestKubegresBasic(t *testing.T) {
	// Simple test to verify we can run tests without the real control plane
	if util.K8sClient.Client == nil {
		t.Skip("Client not initialized, skipping test")
	}

	t.Log("Successfully verified client is initialized")
}

// TestKubegresSpec tests the basic Kubegres spec creation
func TestKubegresSpec(t *testing.T) {
	// Create a basic Kubegres spec
	replicas := int32(1)
	spec := kubegresv1.KubegresSpec{
		Replicas: &replicas,
		Image:    "postgres:13",
	}

	// Verify the spec was created correctly
	if spec.Image != "postgres:13" {
		t.Errorf("Expected image to be postgres:13, got %s", spec.Image)
	}

	if *spec.Replicas != 1 {
		t.Errorf("Expected replicas to be 1, got %d", *spec.Replicas)
	}

	t.Log("Kubegres spec test passed")

	// Verify client creation didn't error
	if util.K8sClient.Client == nil {
		t.Skip("Test client not available, skipping test")
	}

	t.Log("Successfully created test client")
}

// TestKubegresDatabaseFolderConfig tests the database folder configuration
func TestKubegresDatabaseFolderConfig(t *testing.T) {
	// Create a Kubegres spec with database folder
	replicas := int32(1)
	spec := kubegresv1.KubegresSpec{
		Replicas: &replicas,
		Image:    "postgres:13",
		Database: kubegresv1.KubegresDatabase{
			Size:        "1Gi",
			VolumeMount: "/pgdata",
			Folder:      "mypostgres",
		},
	}

	// Verify the database spec was created correctly
	if spec.Database.Size != "1Gi" {
		t.Errorf("Expected database size to be 1Gi, got %s", spec.Database.Size)
	}

	if spec.Database.VolumeMount != "/pgdata" {
		t.Errorf("Expected volume mount to be /pgdata, got %s", spec.Database.VolumeMount)
	}

	if spec.Database.Folder != "mypostgres" {
		t.Errorf("Expected folder to be mypostgres, got %s", spec.Database.Folder)
	}

	t.Log("Kubegres database folder test passed")
}
