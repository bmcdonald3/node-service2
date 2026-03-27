package reconcilers

import (
	"context"
	"time"

	"github.com/openchami/fabrica/pkg/reconcile"
	v1 "github.com/user/node-service/apis/node.openchami.io/v1"
)

type SMDClient interface {
	GetNodesByLabels(labels map[string]string) ([]string, error)
}

type MockSMDClient struct{}

func (m *MockSMDClient) GetNodesByLabels(labels map[string]string) ([]string, error) {
	if labels["role"] == "compute" {
		return []string{"x1000c0s1b0n0", "x1000c0s1b0n1"}, nil
	}
	return []string{}, nil
}

type NodeSetReconciler struct {
	reconcile.BaseReconciler
	SMD SMDClient
}

func (r *NodeSetReconciler) Reconcile(ctx context.Context, resource interface{}) (reconcile.Result, error) {
	if r.SMD == nil {
		r.SMD = &MockSMDClient{}
	}

	nodeSet := resource.(*v1.NodeSet)

	var resolved []string
	resolved = append(resolved, nodeSet.Spec.Xnames...)

	if len(nodeSet.Spec.LabelSelector) > 0 {
		smdNodes, err := r.SMD.GetNodesByLabels(nodeSet.Spec.LabelSelector)
		if err != nil {
			nodeSet.Status.Phase = "Error"
			r.UpdateStatus(ctx, nodeSet)
			return reconcile.Result{RequeueAfter: 1 * time.Minute}, err
		}
		resolved = append(resolved, smdNodes...)
	}

	uniqueMap := make(map[string]bool)
	var finalResolved []string
	for _, xname := range resolved {
		if !uniqueMap[xname] {
			uniqueMap[xname] = true
			finalResolved = append(finalResolved, xname)
		}
	}

	if len(nodeSet.Status.ResolvedXnames) != len(finalResolved) || nodeSet.Status.Phase != "Resolved" {
		nodeSet.Status.ResolvedXnames = finalResolved
		nodeSet.Status.MatchCount = len(finalResolved)
		nodeSet.Status.Phase = "Resolved"
		r.UpdateStatus(ctx, nodeSet)
		r.EmitEvent(ctx, "io.openchami.nodeset.resolved", nodeSet)
	}

	return reconcile.Result{RequeueAfter: 5 * time.Minute}, nil
}

func (r *NodeSetReconciler) GetResourceKind() string {
	return "NodeSet"
}