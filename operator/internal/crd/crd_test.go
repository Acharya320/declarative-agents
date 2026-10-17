// Copyright (c) 2026 Nokia
// SPDX-License-Identifier: BSD-3-Clause

package crd

import (
	"context"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/Nokia-Bell-Labs/declarative-agents/operator/api/v1alpha1"
)

// jsonFields lists a struct's JSON property names.
func jsonFields(t reflect.Type) []string {
	var names []string
	for i := 0; i < t.NumField(); i++ {
		tag := strings.Split(t.Field(i).Tag.Get("json"), ",")[0]
		if tag != "" && tag != "-" {
			names = append(names, tag)
		}
	}
	sort.Strings(names)
	return names
}

func schemaFields(props map[string]apiextensionsv1.JSONSchemaProps) []string {
	var names []string
	for name := range props {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// TestCRDMatchesTypes holds the embedded schema to the Go types (srd009
// R1.1): the spec, status, profile, port, and work-volume properties are
// exactly the types' JSON fields, and the names match the group and kind.
func TestCRDMatchesTypes(t *testing.T) {
	def, err := Definition()
	if err != nil {
		t.Fatal(err)
	}
	if def.Spec.Group != v1alpha1.Group || def.Spec.Names.Kind != v1alpha1.Kind || def.Spec.Names.Plural != v1alpha1.Plural || def.Name != v1alpha1.Plural+"."+v1alpha1.Group {
		t.Fatalf("names = %+v", def.Spec.Names)
	}
	root := def.Spec.Versions[0].Schema.OpenAPIV3Schema.Properties
	spec, status := root["spec"], root["status"]
	for name, c := range map[string]struct {
		got  []string
		want reflect.Type
	}{
		"spec":        {schemaFields(spec.Properties), reflect.TypeOf(v1alpha1.AgentSpec{})},
		"status":      {schemaFields(status.Properties), reflect.TypeOf(v1alpha1.AgentStatus{})},
		"profile":     {schemaFields(spec.Properties["profile"].Properties), reflect.TypeOf(v1alpha1.ProfileSource{})},
		"ports":       {schemaFields(spec.Properties["ports"].Items.Schema.Properties), reflect.TypeOf(v1alpha1.AgentPort{})},
		"workVolumes": {schemaFields(spec.Properties["workVolumes"].Items.Schema.Properties), reflect.TypeOf(v1alpha1.WorkVolume{})},
	} {
		if want := jsonFields(c.want); !reflect.DeepEqual(c.got, want) {
			t.Errorf("%s schema fields = %v, types have %v", name, c.got, want)
		}
	}
	if def.Spec.Versions[0].Subresources == nil || def.Spec.Versions[0].Subresources.Status == nil {
		t.Error("the status subresource is not declared")
	}
}

// TestApplyCreatesAndUpdates covers srd009 R1.3: Apply creates the CRD when
// absent and brings a drifted one back to the embedded definition.
func TestApplyCreatesAndUpdates(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := apiextensionsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	c := fake.NewClientBuilder().WithScheme(scheme).Build()
	ctx := context.Background()
	if err := Apply(ctx, c); err != nil {
		t.Fatal(err)
	}
	var live apiextensionsv1.CustomResourceDefinition
	if err := c.Get(ctx, client.ObjectKey{Name: v1alpha1.Plural + "." + v1alpha1.Group}, &live); err != nil {
		t.Fatal(err)
	}
	live.Spec.Versions[0].Served = false
	if err := c.Update(ctx, &live); err != nil {
		t.Fatal(err)
	}
	if err := Apply(ctx, c); err != nil {
		t.Fatal(err)
	}
	if err := c.Get(ctx, client.ObjectKey{Name: live.Name}, &live); err != nil || !live.Spec.Versions[0].Served {
		t.Fatalf("drift not repaired: served=%v err=%v", live.Spec.Versions[0].Served, err)
	}
}

// TestWaitEstablishedReturnsOnTheCondition waits for the Established
// condition and gives up when the context ends.
func TestWaitEstablishedReturnsOnTheCondition(t *testing.T) {
	scheme := runtime.NewScheme()
	if err := apiextensionsv1.AddToScheme(scheme); err != nil {
		t.Fatal(err)
	}
	def, _ := Definition()
	def.Status.Conditions = []apiextensionsv1.CustomResourceDefinitionCondition{{Type: apiextensionsv1.Established, Status: apiextensionsv1.ConditionTrue}}
	c := fake.NewClientBuilder().WithScheme(scheme).WithObjects(def).Build()
	if err := WaitEstablished(context.Background(), c); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	if err := WaitEstablished(ctx, fake.NewClientBuilder().WithScheme(scheme).Build()); err == nil {
		t.Fatal("an absent CRD counted as established")
	}
}
