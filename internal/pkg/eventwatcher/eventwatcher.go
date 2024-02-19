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

package eventwatcher

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// EventCallback defines the function to be called when an event occurs.
type EventCallback func(object *unstructured.Unstructured)

// EventWatcher watches for events on a specified resource type.
type EventWatcher struct {
	clientset *kubernetes.Clientset
	resource  string
	eventType string
	callback  EventCallback
}

// NewEventWatcher creates a new EventWatcher.
func NewEventWatcher(clientset *kubernetes.Clientset, resource string, eventType string, callback EventCallback) *EventWatcher {
	return &EventWatcher{
		clientset: clientset,
		resource:  resource,
		eventType: eventType,
		callback:  callback,
	}
}

// Run starts the event watcher and blocks until an error occurs.
func (w *EventWatcher) Run(ctx context.Context) error {
	informerFactory := informers.NewSharedInformerFactoryWithOptions(w.clientset, 0, informers.WithNamespace(v1.NamespaceAll))
	resource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: w.resource}
	genericInformer, _ := informerFactory.ForResource(resource)
	informer := genericInformer.Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		w.eventType: func(obj interface{}) {
			object, ok := obj.(*unstructured.Unstructured)
			if !ok {
				return
			}
			if w.callback != nil {
				w.callback(object)
			}
		},
	})

	informerFactory.Start(ctx.Done())
	defer func() {
		for _, informernformer := range informerFactory.WaitForCacheSync(ctx.Done()) {
			informer.Stop()
		}
	}()

	<-ctx.Done()
	return fmt.Errorf("EventWatcher: context canceled")
}
