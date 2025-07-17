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

	v1 "k8s.io/api/core/v1"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/test/util"
)

// TestSetup is a simple test function to verify the test setup is working
func TestSetup(t *testing.T) {
	// Skip this test for now as it's just a placeholder
	t.Skip("This is a placeholder test to verify the test setup")
}

// Helper function to setup a test environment for Kubegres tests
func SetupTestEnvironment(t *testing.T) (*util.K8sClientType, *kubegresv1.Kubegres, *v1.ConfigMap) {
	// This is a helper function that could be used by various tests
	// instead of having the setup code in TestMain
	return &util.K8sClient, nil, nil
}
