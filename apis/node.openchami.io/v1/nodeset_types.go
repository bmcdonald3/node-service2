// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
    "context"
    "github.com/openchami/fabrica/pkg/fabrica"
)

type NodeSet struct {
    APIVersion string           `json:"apiVersion"`
    Kind       string           `json:"kind"`
    Metadata   fabrica.Metadata `json:"metadata"`
    Spec       NodeSetSpec      `json:"spec" validate:"required"`
    Status     NodeSetStatus    `json:"status,omitempty"`
}

type NodeSetSpec struct {
    Xnames        []string          `json:"xnames,omitempty"`
    LabelSelector map[string]string `json:"labelSelector,omitempty"`
}

type NodeSetStatus struct {
    ResolvedXnames []string `json:"resolvedXnames,omitempty"`
    MatchCount     int      `json:"matchCount"`
    Phase          string   `json:"phase,omitempty"`
}

func (ns *NodeSet) Validate(ctx context.Context) error {
    return nil
}
// GetKind returns the kind of the resource
func (r *NodeSet) GetKind() string {
	return "NodeSet"
}

// GetName returns the name of the resource
func (r *NodeSet) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *NodeSet) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *NodeSet) IsHub() {}
