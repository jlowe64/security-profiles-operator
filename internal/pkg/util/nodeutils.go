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
	seccompprofileapi "sigs.k8s.io/security-profiles-operator/api/seccompprofile/v1beta1"
	// selxv1alpha2 "sigs.k8s.io/security-profiles-operator/api/selinuxprofile/v1alpha2"
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

func GetFinalizersFromSeccompProfile(profile *seccompprofileapi.SeccompProfile) []string {
	return profile.GetFinalizers()
}

// Compares nodeNames with finalizers from GetFinalizersFromSeccompProfile
// If there are differences it returns true, else false
func CompareFinalizers(nodeNames []string, finalizers []string) bool {
	// Convert lists to sets for efficient comparison
	nodeSet := make(map[string]struct{}, len(nodeNames))
	for _, name := range nodeNames {
		nodeSet[name] = struct{}{}
	}

	finalizerSet := make(map[string]struct{}, len(finalizers))
	for _, finalizer := range finalizers {
		finalizerSet[finalizer] = struct{}{}
	}

	// Check for differences using set operations
	for node := range nodeSet {
		if _, exists := finalizerSet[node]; !exists {
			return true // Node exists in nodeNames but not in finalizers
		}
	}

	for finalizer := range finalizerSet {
		if _, exists := nodeSet[finalizer]; !exists {
			return true // Finalizer exists in finalizers but not in nodeNames
		}
	}

	return false // No differences found
}

// GetSeccompProfiles retrieves all SeccompProfile resources and returns them as a list
func GetSeccompProfiles(ctx context.Context, client client.Client) (*seccompprofileapi.SeccompProfileList, error) {
	profilesList := &seccompprofileapi.SeccompProfileList{}
	err := client.List(ctx, profilesList, nil)
	if err != nil {
		return nil, err
	}
	return profilesList, nil
}

func GetProfilesObjects(ctx context.Context, client client.Client, namespace, name, kind string) (metav1.Object, error) {
	// Get profiles for seccomp, apparmor, and selinux
	seccompprofile, err := GetSeccompProfiles(ctx, client)
	if err != nil {
		return nil, err
	}
	fmt.Print(seccompprofile)
	return &seccompprofile.Items[0], nil
}

/*func ReconcileFinalizers(ctx context.Context, client client.Client, nodeNames []string, objNames []string, objNamespace, objKind string) error {
	for _, objName := range objNames {
		obj, err := GetKubernetesObject(ctx, client, objNamespace, objName, objKind)
		if err != nil {
			return err
		}

		currentFinalizers := GetFinalizers(ctx, client)
		// ... Logic to compare nodeNames with currentFinalizers
		// ... Update finalizers if necessary

		err = UpdateKubernetesObject(clientset, obj)
		if err != nil {
			return err
		}
	}
	return nil
}*/
