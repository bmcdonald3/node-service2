// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"github.com/openchami/fabrica/pkg/fabrica"
)

// ProfileBinding represents a profilebinding resource
type ProfileBinding struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	Spec       ProfileBindingSpec   `json:"spec" validate:"required"`
	Status     ProfileBindingStatus `json:"status,omitempty"`
}

// ProfileBindingSpec defines the desired state of ProfileBinding
type ProfileBindingSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	// Add your spec fields here
}

// ProfileBindingStatus defines the observed state of ProfileBinding
type ProfileBindingStatus struct {
	Phase      string `json:"phase,omitempty"`
	Message    string `json:"message,omitempty"`
	Ready      bool   `json:"ready"`
		// Add your status fields here
}

// Validate implements custom validation logic for ProfileBinding
func (r *ProfileBinding) Validate(ctx context.Context) error {
	// Add custom validation logic here
	// Example:
	// if r.Spec.Description == "forbidden" {
	//     return errors.New("description 'forbidden' is not allowed")
	// }

	return nil
}
// GetKind returns the kind of the resource
func (r *ProfileBinding) GetKind() string {
	return "ProfileBinding"
}

// GetName returns the name of the resource
func (r *ProfileBinding) GetName() string {
	return r.Metadata.Name
}

// GetUID returns the UID of the resource
func (r *ProfileBinding) GetUID() string {
	return r.Metadata.UID
}

// IsHub marks this as the hub/storage version
func (r *ProfileBinding) IsHub() {}
