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
	"context"
	"fmt"
	"os"
	"testing"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/test/util"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// TestMain is the main entry point for tests in this package
func TestMain(m *testing.M) {
	// Setup test environment
	testEnv := &envtest.Environment{
		CRDDirectoryPaths:     []string{"../../config/crd/bases"},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	if err != nil {
		fmt.Printf("Failed to start test environment: %v\n", err)
		os.Exit(1)
	}

	// Setup scheme
	s := runtime.NewScheme()
	err = kubegresv1.AddToScheme(s)
	if err != nil {
		fmt.Printf("Failed to add Kubegres to scheme: %v\n", err)
		os.Exit(1)
	}

	err = v1.AddToScheme(s)
	if err != nil {
		fmt.Printf("Failed to add v1 to scheme: %v\n", err)
		os.Exit(1)
	}

	// Setup client
	c, err := client.New(cfg, client.Options{Scheme: s})
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Initialize global client
	util.K8sClient = util.K8sClientType{
		Client:   c,
		Ctx:      context.Background(),
		Config:   cfg,
		TestEnv:  testEnv,
		Recorder: record.NewFakeRecorder(100),
	}

	// Run tests
	code := m.Run()

	// Teardown
	err = testEnv.Stop()
	if err != nil {
		fmt.Printf("Failed to stop test environment: %v\n", err)
	}

	os.Exit(code)
}
