/*
Copyright 2020 The Kubernetes Authors.

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

package eventlistener

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"sigs.k8s.io/security-profiles-operator/internal/pkg/eventwatcher"
)

func TestNodeDeleted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock dependencies
	mockClientset := mock_kubernetes_clientset.NewMockClientset(ctrl)
	mockObj := mock_unstructured.NewMockUnstructured(ctrl)

	// Define mock behavior
	mockClientset.EXPECT().ForResource("nodes").Return(nil)
	mockObj.EXPECT().GetName().Return("node1")

	// Create EventWatcher with mock callback
	callbackCalled := false
	callback := func(obj *unstructured.Unstructured) {
		callbackCalled = true
		assert.Equal(t, "node1", obj.GetName())
	}
	watcher := eventwatcher.NewEventWatcher(mockClientset, "nodes", "Deleted", callback)

	// Simulate node deletion
	mockClientset.EXPECT().Actions(gomock.AssignableToTypeOf(&v1.Node{})).DoAndReturn(func(action interface{}) interface{} {
		if action == cache.Delete {
			callback(mockObj)
		}
		return nil
	})

	// Run watcher (mock doesn't actually run)
	watcher.Run(context.Background())

	// Assertions
	assert.True(t, callbackCalled, "Callback should be called for node deletion")
	mock.VerifyAndUnmock(ctrl) // Check all mock expectations met
}
