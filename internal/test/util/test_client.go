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
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

// SetupTestClient sets up the global K8sClient with a fake client for testing
func SetupTestClient() {
	s := scheme.Scheme
	_ = kubegresv1.AddToScheme(s)
	_ = v1.AddToScheme(s)

	client := fake.NewClientBuilder().WithScheme(s).Build()

	K8sClient = K8sClientType{
		Client:   client,
		Ctx:      context.Background(),
		Recorder: &TestEventRecorder{Events: []EventRecord{}},
	}
}

// K8sClientType represents a Kubernetes client for testing
type K8sClientType struct {
	Client   client.Client
	Ctx      context.Context
	Recorder *TestEventRecorder
}

// K8sClient is available for tests
var K8sClient K8sClientType
