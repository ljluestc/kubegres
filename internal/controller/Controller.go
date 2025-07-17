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

package controller

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateKubegresControllerV1 creates a new Kubegres controller instance
func CreateKubegresControllerV1(c client.Client, ctx context.Context, req ctrl.Request) interface{} {
	return &KubegresController{
		Client:  c,
		Ctx:     ctx,
		Request: req,
	}
}

// KubegresController handles reconciliation of Kubegres resources
type KubegresController struct {
	Client  client.Client
	Ctx     context.Context
	Request ctrl.Request
}

// Reconcile performs reconciliation for a Kubegres resource
func (r *KubegresController) Reconcile() (ctrl.Result, error) {
	// Temporary stub implementation
	return ctrl.Result{}, nil
}
