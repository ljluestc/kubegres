package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/test/util"
)

var K8sClient util.K8sClientType

func init() {
	// Initialize the test environment
	err := util.Setup()
	if err != nil {
		fmt.Printf("Error setting up test environment: %v\n", err)
		os.Exit(1)
	}
	K8sClient = util.K8sClient
}

func TestMain(m *testing.M) {
	// Run tests
	code := m.Run()

	// Tear down the test environment
	if err := util.Teardown(); err != nil {
		fmt.Printf("Error tearing down test environment: %v\n", err)
	}

	os.Exit(code)
}

// LoadCustomConfig uses the utility function to load custom configuration
func LoadCustomConfig(kubegresConfigPath, configMapPath string) *util.CustomConfig {
	return util.LoadCustomConfig(kubegresConfigPath, configMapPath)
}

// NewKubegresResources creates a new resources instance for testing
func NewKubegresResources(client *util.K8sClientType, kubegres *kubegresv1.Kubegres) *util.KubegresResources {
	return util.NewKubegresResources(client, kubegres)
}

func checkIfFailoverWasTriggeredInLogs(logs string) bool {
	return strings.Contains(logs, "Started failing-over") ||
		strings.Contains(logs, "FailOver: Promoting Replica to Primary")
}

func checkIfSTONITHWasTriggeredInLogs(logs string) bool {
	return strings.Contains(logs, "STONITH: Verifying primary is fully terminated") ||
		strings.Contains(logs, "STONITHEnabled")
}

func TestFailOverWithSTONITHEnabled(t *testing.T) {
	t.Skip("Test skipped while fixing infrastructure")

	// Given a Kubegres with replica count set to 3, customConfig and STONITH enabled
	customConfig := LoadCustomConfig(
		"../controller/test/config/primary-replica-connection-check.yaml",
		"../controller/test/config/primary-replica-connection-check-configmap.yaml")

	// Create Kubegres with STONITH enabled
	enableSTONITH := true
	replicas := int32(3)
	kubegres := kubegresv1.Kubegres{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-failover-stonith",
			Namespace: "default",
		},
		Spec: kubegresv1.KubegresSpec{
			Replicas: &replicas,
			Image:    "postgres:12.4",
			Failover: kubegresv1.KubegresFailover{
				EnableSTONITH: &enableSTONITH,
			},
		},
	}

	// Apply the custom configuration
	if customConfig != nil {
		// Create a custom config properly typed
		kubegres.Spec.CustomConfig = customConfig.KubegresConfig
	}

	err := K8sClient.Client.Create(K8sClient.Ctx, &kubegres)
	if err != nil {
		t.Fatal(err)
	}

	kubegresResources := NewKubegresResources(&K8sClient, &kubegres)

	// Given we wait for the statefulSet and all its pods to be ready
	err = kubegresResources.WaitForStatefulSetAndPodsReady("test-failover-stonith", 3)
	if err != nil {
		t.Fatal(err)
	}

	// Given we make sure we have a Primary instance
	time.Sleep(15 * time.Second)

	// Given
	primaryPod, err := kubegresResources.GetPrimaryPod()
	if err != nil {
		t.Fatal(err)
	}

	primaryPodName := primaryPod.Name
	primaryNodeName := primaryPod.Spec.NodeName
	nodePrimary := &v1.Node{}

	// Given we get the node primary
	err = K8sClient.Client.Get(
		K8sClient.Ctx,
		types.NamespacedName{Namespace: "", Name: primaryNodeName},
		nodePrimary)

	if err != nil {
		t.Fatal(fmt.Errorf("we cannot find the node where the primary resides: %w", err))
	}

	// Given we make it not schedulable
	nodePrimary.Spec.Unschedulable = true

	// When we patch the node
	err = K8sClient.Client.Update(K8sClient.Ctx, nodePrimary)
	if err != nil {
		t.Fatal(err)
	}

	// Then a new Primary is elected and the old primary is removed
	err = kubegresResources.WaitForAPrimaryToBeElected(func(pod v1.Pod) bool {
		return pod.Name != primaryPodName
	})
	if err != nil {
		t.Fatal(err)
	}

	logs, err := kubegresResources.GetLogs()
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(logs)

	// Verify both failover and STONITH were triggered
	if !checkIfFailoverWasTriggeredInLogs(logs) {
		t.Fatal("It was expected to have a failover triggered")
	}

	if !checkIfSTONITHWasTriggeredInLogs(logs) {
		t.Fatal("It was expected to have STONITH triggered")
	}

	// Cleanup
	// Given we make the node schedulable again
	nodePrimary.Spec.Unschedulable = false
	err = K8sClient.Client.Update(K8sClient.Ctx, nodePrimary)
	if err != nil {
		t.Fatal(err)
	}

	// Delete the test resources
	err = K8sClient.Client.Delete(K8sClient.Ctx, &kubegres)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFailOverPerformanceWithSTONITH(t *testing.T) {
	t.Skip("Test skipped while fixing infrastructure")

	// Given a Kubegres with replica count set to 3, customConfig and STONITH enabled
	customConfig := LoadCustomConfig(
		"../controller/test/config/primary-replica-connection-check.yaml",
		"../controller/test/config/primary-replica-connection-check-configmap.yaml")

	// Create Kubegres with STONITH enabled
	enableSTONITH := true
	replicas := int32(3)
	kubegres := kubegresv1.Kubegres{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-failover-perf",
			Namespace: "default",
		},
		Spec: kubegresv1.KubegresSpec{
			Replicas: &replicas,
			Image:    "postgres:12.4",
			Failover: kubegresv1.KubegresFailover{
				EnableSTONITH: &enableSTONITH,
			},
		},
	}

	// Apply the custom configuration
	if customConfig != nil {
		kubegres.Spec.CustomConfig = customConfig.KubegresConfig
	}

	err := K8sClient.Client.Create(K8sClient.Ctx, &kubegres)
	if err != nil {
		t.Fatal(err)
	}

	kubegresResources := NewKubegresResources(&K8sClient, &kubegres)

	// Given we wait for the statefulSet and all its pods to be ready
	err = kubegresResources.WaitForStatefulSetAndPodsReady("test-failover-perf", 3)
	if err != nil {
		t.Fatal(err)
	}

	// Given we make sure we have a Primary instance
	time.Sleep(15 * time.Second)

	// Given
	primaryPod, err := kubegresResources.GetPrimaryPod()
	if err != nil {
		t.Fatal(err)
	}

	primaryPodName := primaryPod.Name

	// Delete the primary pod to trigger failover
	err = K8sClient.Client.Delete(K8sClient.Ctx, &primaryPod)
	if err != nil {
		t.Fatal(err)
	}

	startTime := time.Now()

	// Then a new Primary is elected
	err = kubegresResources.WaitForAPrimaryToBeElected(func(pod v1.Pod) bool {
		return pod.Name != primaryPodName
	})
	if err != nil {
		t.Fatal(err)
	}

	failoverDuration := time.Since(startTime)

	// Failover should complete within reasonable time (30 seconds)
	if failoverDuration > 30*time.Second {
		t.Fatalf("Failover with STONITH took too long: %v", failoverDuration)
	}

	fmt.Printf("Failover with STONITH completed in %v\n", failoverDuration)

	// Cleanup
	// Delete the test resources
	err = K8sClient.Client.Delete(K8sClient.Ctx, &kubegres)
	if err != nil {
		t.Fatal(err)
	}
}
