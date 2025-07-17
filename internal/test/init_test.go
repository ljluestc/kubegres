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
	"fmt"
	"os"
	"testing"

	"reactive-tech.io/kubegres/internal/test/util"
)

// TestMain is the main entry point for all tests in this package
func TestMain(m *testing.M) {
	// Initialize the test environment
	if err := util.Setup(); err != nil {
		fmt.Printf("Error setting up test environment: %v\n", err)
		os.Exit(1)
	}

	// Run the tests
	code := m.Run()

	// Clean up the test environment
	if err := util.Teardown(); err != nil {
		fmt.Printf("Error tearing down test environment: %v\n", err)
	}

	os.Exit(code)
}
