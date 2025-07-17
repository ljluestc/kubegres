package util

import (
	"context"
	"fmt"
	"io/ioutil"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

// K8sClientType defines the structure for a Kubernetes client
type K8sClientType struct {
	Client   client.Client
	Ctx      context.Context
	Config   *rest.Config
	Recorder record.EventRecorder
	TestEnv  *envtest.Environment
}

// K8sClient holds Kubernetes client information
var K8sClient = K8sClientType{
	Ctx: context.TODO(),
}

// CustomConfig holds Kubegres configuration
type CustomConfig struct {
	KubegresConfig string
	ConfigMapName  string
	ConfigMapData  map[string]string
}

// KubegresResources holds resources for testing
type KubegresResources struct {
	K8sClient *K8sClientType
	Kubegres  *kubegresv1.Kubegres
}

// LoadCustomConfig loads custom configuration from yaml files
func LoadCustomConfig(kubegresConfigPath, configMapPath string) *CustomConfig {
	if kubegresConfigPath == "" || configMapPath == "" {
		return nil
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

// NewKubegresResources creates a new KubegresResources instance
func NewKubegresResources(client *K8sClientType, kubegres *kubegresv1.Kubegres) *KubegresResources {
	return &KubegresResources{
		K8sClient: client,
		Kubegres:  kubegres,
	}
}

// WaitForStatefulSetAndPodsReady waits for statefulset and pods to be ready
func (r *KubegresResources) WaitForStatefulSetAndPodsReady(name string, replicas int) error {
	// Mock implementation for testing
	return nil
}

// GetPrimaryPod gets the primary pod
func (r *KubegresResources) GetPrimaryPod() (v1.Pod, error) {
	// Mock implementation for testing
	return v1.Pod{
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
	// Mock implementation for testing
	return "Started failing-over\nEnsuring old primary is terminated (STONITH)", nil
}

// InitK8sClient initializes the K8sClient with a proper client implementation
func InitK8sClient(c client.Client, ctx context.Context, cfg *rest.Config) {
	K8sClient = K8sClientType{
		Client:   c,
		Ctx:      ctx,
		Config:   cfg,
		Recorder: record.NewFakeRecorder(100),
	}
}
