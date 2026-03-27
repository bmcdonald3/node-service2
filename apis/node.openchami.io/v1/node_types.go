// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
    "context"
    "github.com/openchami/fabrica/pkg/fabrica"
)

type Node struct {
    APIVersion string           `json:"apiVersion"`
    Kind       string           `json:"kind"`
    Metadata   fabrica.Metadata `json:"metadata"`
    Spec       NodeSpec         `json:"spec" validate:"required"`
    Status     NodeStatus       `json:"status,omitempty"`
}

type NodeSpec struct {
    Xname  string            `json:"xname" validate:"required"`
    Role   string            `json:"role,omitempty"`
    Labels map[string]string `json:"labels,omitempty" validate:"dive,keys,labelkey,endkeys,labelvalue"`
}

type NodeStatus struct {
    EffectiveProfile string `json:"effectiveProfile,omitempty"`
    Phase            string `json:"phase,omitempty"`
}

func (n *Node) Validate(ctx context.Context) error {
    return nil
}

// GetKind returns the kind of the resource
func (r *Node) GetKind() string {
	return "Node"
}

// GetName returns the name of the resource
func (r *Node) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *Node) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *Node) IsHub() {}
