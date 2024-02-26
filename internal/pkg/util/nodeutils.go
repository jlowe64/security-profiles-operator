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

package util

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"

	// apparmorapi "sigs.k8s.io/security-profiles-operator/api/apparmorprofile/v1alpha1"

	// selxv1alpha2 "sigs.k8s.io/security-profiles-operator/api/selinuxprofile/v1alpha2"
	statusv1alpha1 "sigs.k8s.io/security-profiles-operator/api/secprofnodestatus/v1alpha1"
)

func GetDynamicClient(client client.Client) (dynamic.Interface, error) {
	clusterConfig, err := rest.InClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("get in-cluster config: %w", err)
	}

	// Create a dynamic client for working with nodes
	dynamicClient, err := dynamic.NewForConfig(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("load dynamic client: %w", err)
	}
	return dynamicClient, nil
}

func GetNodeList(ctx context.Context, client client.Client) ([]string, error) {
	dynamicClient, err := GetDynamicClient(client)
	if err != nil {
		return nil, err
	}
	// Specify the resource (nodes) and namespace
	nodeResource := schema.GroupVersionResource{Group: "", Version: "v1", Resource: "nodes"}

	// List the nodes (using the dynamic client)
	nodeList, err := dynamicClient.Resource(nodeResource).Namespace("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Extract node names
	var nodeNames []string
	for _, item := range nodeList.Items {
		var node map[string]interface{}
		err := runtime.DefaultUnstructuredConverter.FromUnstructured(item.Object, &node)
		if err != nil {
			return nil, err
		}
		nodeName, _, err := unstructured.NestedString(node, "metadata", "name")
		if err != nil {
			return nil, err
		}
		nodeNames = append(nodeNames, nodeName)
	}

	return nodeNames, nil
}

// Compares nodeNames with finalizers from GetFinalizersFromSeccompProfile
// If there are differences it returns true, else false
func CompareFinalizers(ctx context.Context, nodeStatusList *statusv1alpha1.SecurityProfileNodeStatusList, client client.Client) (bool, error) {
	nodeNames, err := GetNodeList(ctx, client)
	if err != nil {
		return false, err
	}
	// Convert lists to sets for efficient comparison
	nodeSet := make(map[string]struct{}, len(nodeNames))
	for _, name := range nodeNames {
		nodeSet[name] = struct{}{}
	}

	nodeStatusListSet := make(map[string]struct{}, len(nodeStatusList.Items))
	for _, finalizer := range nodeStatusList.Items {
		nodeStatusListSet[finalizer.NodeName] = struct{}{}
	}

	// Check for differences using set operations
	for node := range nodeSet {
		if _, exists := nodeStatusListSet[node]; !exists {
			return true, nil // Node exists in nodeNames but not in finalizers
		}
	}

	for finalizer := range nodeStatusListSet {
		if _, exists := nodeSet[finalizer]; !exists {
			return true, nil // Finalizer exists in finalizers but not in nodeNames
		}
	}

	return false, nil // No differences found
}
