package controllers

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

import (
	"context"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/controller"
)

// KubegresReconciler reconciles a Kubegres object
type KubegresReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=kubegres.reactive-tech.io,resources=kubegres,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=kubegres.reactive-tech.io,resources=kubegres/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=kubegres.reactive-tech.io,resources=kubegres/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=statefulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=apps,resources=statefulsets/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=configmaps,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=configmaps/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get
//+kubebuilder:rbac:groups=batch,resources=cronjobs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch,resources=cronjobs/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
func (r *KubegresReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := ctrl.Log.WithName("controllers").WithName("Kubegres")
	log.Info("Reconciling Kubegres", "namespace", req.Namespace, "name", req.Name)

	// Use the original KubegresReconciler from the internal controller package
	kubegresController := controller.NewKubegresController(r.Client, ctx, req)
	kubegresCtrl, ok := kubegresController.(*controller.KubegresReconciler)
	if !ok {
		log.Error(nil, "Error casting controller interface to KubegresReconciler")
		return ctrl.Result{}, nil
	}
	return kubegresCtrl.Reconcile(ctx, req)
}

// SetupWithManager sets up the controller with the Manager.
func (r *KubegresReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kubegresv1.Kubegres{}).
		Complete(r)
}

func SetupLogger() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))
}
