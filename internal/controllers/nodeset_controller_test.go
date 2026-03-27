package controllers_test

import (
	"context"
	"testing"

	"github.com/openchami/fabrica/pkg/events"
	"github.com/openchami/fabrica/pkg/fabrica"
	"github.com/openchami/fabrica/pkg/reconcile"
	v1 "github.com/user/node-service/apis/node.openchami.io/v1"
	"github.com/user/node-service/internal/controllers"
)

func TestNodeSetReconciler(t *testing.T) {
	eventBus := events.NewInMemoryEventBus(10, 1)
	reconciler := &controllers.NodeSetReconciler{
		BaseReconciler: reconcile.BaseReconciler{
			EventBus: eventBus,
			Logger:   reconcile.NewDefaultLogger(),
		},
		SMD: &controllers.MockSMDClient{},
	}

	nodeSet := &v1.NodeSet{
		Kind: "NodeSet",
		Metadata: fabrica.Metadata{
			UID: "ns-123",
		},
		Spec: v1.NodeSetSpec{
			Xnames: []string{"x1000c1s1b0n0"},
			LabelSelector: map[string]string{
				"role": "compute",
			},
		},
	}

	_, err := reconciler.Reconcile(context.Background(), nodeSet)
	if err != nil {
		t.Fatalf("Reconciliation failed: %v", err)
	}

	if nodeSet.Status.Phase != "Resolved" {
		t.Errorf("Expected phase Resolved, got %s", nodeSet.Status.Phase)
	}

	if nodeSet.Status.MatchCount != 3 {
		t.Errorf("Expected MatchCount 3, got %d", nodeSet.Status.MatchCount)
	}

	expected := map[string]bool{"x1000c1s1b0n0": true, "x1000c0s1b0n0": true, "x1000c0s1b0n1": true}
	for _, xname := range nodeSet.Status.ResolvedXnames {
		if !expected[xname] {
			t.Errorf("Unexpected xname in resolved list: %s", xname)
		}
	}
}