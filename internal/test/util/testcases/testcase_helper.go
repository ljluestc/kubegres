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

package testcases

import (
	"fmt"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubegresv1 "reactive-tech.io/kubegres/api/v1"
	"reactive-tech.io/kubegres/internal/test/util"
)

// TestCase represents a test case for Kubegres
type TestCase struct {
	Name              string
	Kubegres          *kubegresv1.Kubegres
	CustomConfig      *util.CustomConfig
	ExpectFailover    bool
	ExpectSTONITH     bool
	Setup             func(t *testing.T, tc *TestCase) error
	TriggerAction     func(t *testing.T, tc *TestCase) error
	Verification      func(t *testing.T, tc *TestCase) error
	Cleanup           func(t *testing.T, tc *TestCase) error
	Client            *util.K8sClientType
	KubegresResources *util.KubegresResources
	PrimaryPod        v1.Pod
	PrimaryPodName    string
	StartTime         time.Time
	FailoverDuration  time.Duration
}

// NewTestCase creates a new test case
func NewTestCase(name string) *TestCase {
	return &TestCase{
		Name:   name,
		Client: &util.K8sClient,
	}
}

// Run executes the test case
func (tc *TestCase) Run(t *testing.T) {
	t.Run(tc.Name, func(t *testing.T) {
		// Skip if marked to skip
		if tc.Kubegres == nil {
			t.Skip("Test case not fully configured, skipping")
			return
		}

		// Apply the custom configuration if provided
		if tc.CustomConfig != nil {
			tc.Kubegres.Spec.CustomConfig = tc.CustomConfig.KubegresConfig

			// Create the ConfigMap
			configMap := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Name:      tc.CustomConfig.ConfigMapName,
					Namespace: tc.Kubegres.Namespace,
				},
				Data: tc.CustomConfig.ConfigMapData,
			}

			err := tc.Client.Client.Create(tc.Client.Ctx, configMap)
			if err != nil {
				t.Fatalf("Failed to create ConfigMap: %v", err)
			}
		}

		// Create the Kubegres resource
		err := tc.Client.Client.Create(tc.Client.Ctx, tc.Kubegres)
		if err != nil {
			t.Fatalf("Failed to create Kubegres: %v", err)
		}

		// Initialize KubegresResources
		tc.KubegresResources = util.NewKubegresResources(tc.Client, tc.Kubegres)

		// Run setup if provided
		if tc.Setup != nil {
			if err := tc.Setup(t, tc); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
		}

		// Run the trigger action if provided
		if tc.TriggerAction != nil {
			tc.StartTime = time.Now()
			if err := tc.TriggerAction(t, tc); err != nil {
				t.Fatalf("Trigger action failed: %v", err)
			}
		}

		// Run verification if provided
		if tc.Verification != nil {
			if err := tc.Verification(t, tc); err != nil {
				t.Fatalf("Verification failed: %v", err)
			}
		}

		// Run cleanup if provided
		if tc.Cleanup != nil {
			if err := tc.Cleanup(t, tc); err != nil {
				t.Logf("Cleanup failed: %v", err)
			}
		} else {
			// Default cleanup - delete the Kubegres resource
			if tc.Kubegres != nil {
				err = tc.Client.Client.Delete(tc.Client.Ctx, tc.Kubegres)
				if err != nil {
					t.Logf("Failed to delete Kubegres: %v", err)
				}
			}
		}
	})
}

// CreateBasicKubegres creates a basic Kubegres instance for testing
func CreateBasicKubegres(name string, replicas int32, enableSTONITH bool) *kubegresv1.Kubegres {
	return &kubegresv1.Kubegres{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: kubegresv1.KubegresSpec{
			Replicas: &replicas,
			Image:    "postgres:13.4",
			Failover: kubegresv1.KubegresFailover{
				EnableSTONITH: &enableSTONITH,
			},
		},
	}
}

// DefaultSetup provides a default setup for tests
func DefaultSetup(t *testing.T, tc *TestCase) error {
	// Wait for StatefulSet and pods to be ready
	var replicas int
	if tc.Kubegres.Spec.Replicas != nil {
		replicas = int(*tc.Kubegres.Spec.Replicas)
	} else {
		replicas = 1 // Default to 1 if not specified
	}

	err := tc.KubegresResources.WaitForStatefulSetAndPodsReady(tc.Kubegres.Name, replicas)
	if err != nil {
		return fmt.Errorf("failed waiting for StatefulSet and pods to be ready: %w", err)
	}

	// Allow time for cluster to stabilize
	time.Sleep(15 * time.Second)

	// Get the primary pod
	tc.PrimaryPod, err = tc.KubegresResources.GetPrimaryPod()
	if err != nil {
		return fmt.Errorf("failed to get primary pod: %w", err)
	}
	tc.PrimaryPodName = tc.PrimaryPod.Name

	return nil
}

// DefaultVerification provides a default verification for failover tests
func DefaultVerification(t *testing.T, tc *TestCase) error {
	// Wait for a new primary to be elected
	err := tc.KubegresResources.WaitForAPrimaryToBeElected(func(pod v1.Pod) bool {
		return pod.Name != tc.PrimaryPodName
	})
	if err != nil {
		return fmt.Errorf("failed waiting for new primary to be elected: %w", err)
	}

	// Get logs and check for expected patterns
	logs, err := tc.KubegresResources.GetLogs()
	if err != nil {
		return fmt.Errorf("failed to get logs: %w", err)
	}

	// Check if failover was triggered
	if tc.ExpectFailover && !checkIfFailoverWasTriggeredInLogs(logs) {
		return fmt.Errorf("expected failover to be triggered but it wasn't")
	}

	// Check if STONITH was triggered
	if tc.ExpectSTONITH && !checkIfSTONITHWasTriggeredInLogs(logs) {
		return fmt.Errorf("expected STONITH to be triggered but it wasn't")
	}

	return nil
}

// Helper function to check if failover was triggered in logs
func checkIfFailoverWasTriggeredInLogs(logs string) bool {
	return containsAny(logs, []string{
		"FailOver: Promoting Replica to Primary",
		"Started failing-over",
		"Failover triggered",
	})
}

// Helper function to check if STONITH was triggered in logs
func checkIfSTONITHWasTriggeredInLogs(logs string) bool {
	return containsAny(logs, []string{
		"STONITH: Verifying primary is fully terminated",
		"STONITHEnabled",
		"STONITH mechanism is enabled",
	})
}

// Helper function to check if any of the strings are contained in the text
func containsAny(text string, substrs []string) bool {
	for _, substr := range substrs {
		if strings.Contains(text, substr) {
			return true
		}
	}
	return false
}
