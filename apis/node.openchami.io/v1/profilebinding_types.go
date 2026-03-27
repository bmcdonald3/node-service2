// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
    "context"
    "github.com/openchami/fabrica/pkg/fabrica"
)

type ProfileBinding struct {
    APIVersion string               `json:"apiVersion"`
    Kind       string               `json:"kind"`
    Metadata   fabrica.Metadata     `json:"metadata"`
    Spec       ProfileBindingSpec   `json:"spec" validate:"required"`
    Status     ProfileBindingStatus `json:"status,omitempty"`
}

type ProfileBindingSpec struct {
    TargetKind          string `json:"targetKind" validate:"required,oneof=Node NodeSet"`
    TargetName          string `json:"targetName" validate:"required"`
    Profile             string `json:"profile" validate:"required"`
    BootProfileOverride string `json:"bootProfileOverride,omitempty"`
}

type ProfileBindingStatus struct {
    MaterializedBoot     bool   `json:"materializedBoot"`
    MaterializedMetadata bool   `json:"materializedMetadata"`
    Phase                string `json:"phase,omitempty"`
}

func (pb *ProfileBinding) Validate(ctx context.Context) error {
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
