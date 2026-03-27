// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// Node represents a node resource
type Node struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       NodeSpec   `json:"spec" validate:"required"`
	Status     NodeStatus `json:"status,omitempty"`
}

// NodeSpec defines the desired state of Node
type NodeSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	// Add your spec fields here
}

// NodeStatus defines the observed state of Node
type NodeStatus struct {
	Phase      string `json:"phase,omitempty"`
	Message    string `json:"message,omitempty"`
	Ready      bool   `json:"ready"`
		// Add your status fields here
}

// Validate implements custom validation logic for Node
func (r *Node) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

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
