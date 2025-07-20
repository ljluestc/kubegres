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
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// EventRecord represents a recorded event
type EventRecord struct {
	Eventtype string
	Reason    string
	Message   string
}

// TestEventRecorder is a test implementation of the EventRecorder interface
type TestEventRecorder struct {
	Events         []EventRecord
	client         client.Client
	Namespace      string
	RecordedEvents []EventRecord
	mu             sync.Mutex
}

// Event records an event
func (r *TestEventRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	event := EventRecord{
		Eventtype: eventtype,
		Reason:    reason,
		Message:   message,
	}

	r.Events = append(r.Events, event)
	r.RecordedEvents = append(r.RecordedEvents, event)
}

// Eventf records an event with formatted message
func (r *TestEventRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	message := fmt.Sprintf(messageFmt, args...)
	event := EventRecord{
		Eventtype: eventtype,
		Reason:    reason,
		Message:   message,
	}

	r.Events = append(r.Events, event)
	r.RecordedEvents = append(r.RecordedEvents, event)
}

// AnnotatedEventf records an event with annotations
func (r *TestEventRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()

	message := fmt.Sprintf(messageFmt, args...)
	event := EventRecord{
		Eventtype: eventtype,
		Reason:    reason,
		Message:   message,
	}

	r.Events = append(r.Events, event)
	r.RecordedEvents = append(r.RecordedEvents, event)
}

// CheckEventExist checks if an event exists
func (r *TestEventRecorder) CheckEventExist(expectedEvent EventRecord) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check in both Events and RecordedEvents arrays
	for _, event := range r.Events {
		if event.Eventtype == expectedEvent.Eventtype &&
			event.Reason == expectedEvent.Reason &&
			event.Message == expectedEvent.Message {
			return true
		}
	}

	for _, event := range r.RecordedEvents {
		if event.Eventtype == expectedEvent.Eventtype &&
			event.Reason == expectedEvent.Reason &&
			event.Message == expectedEvent.Message {
			return true
		}
	}

	return false
}

// GetEvents returns all recorded events
func (r *TestEventRecorder) GetEvents() []EventRecord {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Combine both event arrays
	allEvents := append([]EventRecord{}, r.Events...)
	allEvents = append(allEvents, r.RecordedEvents...)
	return allEvents
}

// Clear clears all recorded events
func (r *TestEventRecorder) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Events = make([]EventRecord, 0)
	r.RecordedEvents = make([]EventRecord, 0)
}

// GetTestEventRecorder retrieves or creates a test event recorder
func GetTestEventRecorder(client K8sClientType) *TestEventRecorder {
	if client.Recorder == nil {
		client.Recorder = &TestEventRecorder{Events: []EventRecord{}}
	}
	if client.Recorder != nil {
		if client.Recorder == nil {
			client.Recorder = &TestEventRecorder{Events: []EventRecord{}}
		}
		return client.Recorder
	}

	// If the client's recorder is nil, create a new one
	testRecorder := &TestEventRecorder{
		Events:         make([]EventRecord, 0),
		RecordedEvents: make([]EventRecord, 0),
	}

	K8sClient.Recorder = testRecorder

	return testRecorder
}

// CreateTestEventRecorder creates a new TestEventRecorder
func CreateTestEventRecorder(c client.Client, namespace string) *TestEventRecorder {
	return &TestEventRecorder{
		client:         c,
		Namespace:      namespace,
		Events:         make([]EventRecord, 0),
		RecordedEvents: make([]EventRecord, 0),
	}
}
