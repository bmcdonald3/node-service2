// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
// SPDX-License-Identifier: MIT

package reconcilers

import (
	"context"

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

var smdClient = &MockSMDClient{}

func (r *NodeSetReconciler) reconcileNodeSet(ctx context.Context, res *v1.NodeSet) error {
	var resolved []string
	resolved = append(resolved, res.Spec.Xnames...)

	if len(res.Spec.LabelSelector) > 0 {
		smdNodes, err := smdClient.GetNodesByLabels(res.Spec.LabelSelector)
		if err != nil {
			res.Status.Phase = "Error"
			r.UpdateStatus(ctx, res)
			return err 
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

	if len(res.Status.ResolvedXnames) != len(finalResolved) || res.Status.Phase != "Resolved" {
		res.Status.ResolvedXnames = finalResolved
		res.Status.MatchCount = len(finalResolved)
		res.Status.Phase = "Resolved"
		r.UpdateStatus(ctx, res)
		r.EmitEvent(ctx, "io.openchami.nodeset.resolved", res)
	}

	return nil
}