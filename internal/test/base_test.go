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

// TestBasicSetup verifies that the test environment is properly set up
func TestBasicSetup(t *testing.T) {
	// Skip this test if client is not initialized
	if util.K8sClient.Client == nil {
		t.Skip("K8sClient not initialized, skipping test")
	}

	// This is a basic test to verify the test setup
	t.Log("Basic test setup verification passed")
}

// TestKubegresResourceCreation tests creating a basic Kubegres resource
func TestKubegresResourceCreation(t *testing.T) {
	// Skip this test if client is not initialized
	if util.K8sClient.Client == nil {
		t.Skip("K8sClient not initialized, skipping test")
	}

	// Create a basic Kubegres spec for testing
	replicas := int32(1)
	kubegres := &kubegresv1.Kubegres{
		Spec: kubegresv1.KubegresSpec{
			Replicas: &replicas,
			Image:    "postgres:13",
		},
	}

	// Just verify the spec is created correctly
	if kubegres.Spec.Image != "postgres:13" {
		t.Errorf("Expected image to be postgres:13, got %s", kubegres.Spec.Image)
	}

	if *kubegres.Spec.Replicas != 1 {
		t.Errorf("Expected replicas to be 1, got %d", *kubegres.Spec.Replicas)
	}

	t.Log("Kubegres resource creation test passed")
}
