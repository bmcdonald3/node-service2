package controllers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/openchami/fabrica/pkg/events"
	"github.com/openchami/fabrica/pkg/fabrica"
	"github.com/openchami/fabrica/pkg/reconcile"
	v1 "github.com/user/node-service/apis/node.openchami.io/v1"
	"github.com/user/node-service/internal/controllers"
)

type MockDownstreamClient struct {
	ShouldFail bool
	Called     bool
}

func (m *MockDownstreamClient) ApplyProfile(ctx context.Context, targetName, profile string) error {
	m.Called = true
	if m.ShouldFail {
		return errors.New("mock network error")
	}
	return nil
}

func TestProfileBindingReconciler(t *testing.T) {
	eventBus := events.NewInMemoryEventBus(10, 1)

	mockMeta := &MockDownstreamClient{}
	mockBoot := &MockDownstreamClient{}

	reconciler := &controllers.ProfileBindingReconciler{
		BaseReconciler: reconcile.BaseReconciler{
			EventBus: eventBus,
			Logger:   reconcile.NewDefaultLogger(),
		},
		MetadataClient: mockMeta,
		BootClient:     mockBoot,
	}

	binding := &v1.ProfileBinding{
		Kind: "ProfileBinding",
		Metadata: fabrica.Metadata{
			UID: "pb-123",
		},
		Spec: v1.ProfileBindingSpec{
			TargetKind: "Node",
			TargetName: "x1000c0s1b0n0",
			Profile:    "compute-new",
		},
	}

	_, err := reconciler.Reconcile(context.Background(), binding)
	if err != nil {
		t.Fatalf("Reconciliation failed: %v", err)
	}

	if !mockMeta.Called || !mockBoot.Called {
		t.Errorf("Expected both downstream clients to be called")
	}

	if binding.Status.Phase != "Bound" {
		t.Errorf("Expected phase Bound, got %s", binding.Status.Phase)
	}

	mockMeta.Called = false
	mockMeta.ShouldFail = true
	binding.Status.MaterializedMetadata = false
	binding.Status.Phase = ""

	result, err := reconciler.Reconcile(context.Background(), binding)
	if err == nil {
		t.Errorf("Expected error from failing metadata client")
	}

	if result.RequeueAfter == 0 {
		t.Errorf("Expected requeue after failure")
	}

	if binding.Status.Phase != "MetadataError" {
		t.Errorf("Expected phase MetadataError, got %s", binding.Status.Phase)
	}
}