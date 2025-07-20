// Package test contains integration tests for Kubegres
package test

import (
	"fmt"
	"os"
	"testing"

	"reactive-tech.io/kubegres/internal/test/util"
)

// TestMain handles test setup and teardown
func TestMain(m *testing.M) {
	fmt.Println("Setting up test environment")

	// Setup the test client
	util.SetupTestClient()

	fmt.Println("Using fake client for testing")

	// Run tests
	code := m.Run()

	fmt.Println("Tearing down test environment")

	// Exit with appropriate code
	os.Exit(code)
}
