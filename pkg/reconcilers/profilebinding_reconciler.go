// Copyright © 2025 OpenCHAMI a Series of LF Projects, LLC
// SPDX-License-Identifier: MIT

package reconcilers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	v1 "github.com/user/node-service/apis/node.openchami.io/v1"
)

type HTTPDownstreamClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func (c *HTTPDownstreamClient) ApplyProfile(ctx context.Context, targetName, profile string) error {
	payload := map[string]string{"profile": profile}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, fmt.Sprintf("%s/bindings/%s", c.BaseURL, targetName), bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("downstream returned status %d", resp.StatusCode)
	}
	return nil
}

var (
	metadataClient = &HTTPDownstreamClient{BaseURL: "http://localhost:8081", HTTPClient: &http.Client{}}
	bootClient     = &HTTPDownstreamClient{BaseURL: "http://localhost:8081", HTTPClient: &http.Client{}}
)

func (r *ProfileBindingReconciler) reconcileProfileBinding(ctx context.Context, res *v1.ProfileBinding) error {
	if !res.Status.MaterializedMetadata {
		err := metadataClient.ApplyProfile(ctx, res.Spec.TargetName, res.Spec.Profile)
		if err != nil {
			res.Status.Phase = "MetadataError"
			r.UpdateStatus(ctx, res)
			return err 
		}
		res.Status.MaterializedMetadata = true
	}

	if !res.Status.MaterializedBoot {
		err := bootClient.ApplyProfile(ctx, res.Spec.TargetName, res.Spec.Profile)
		if err != nil {
			res.Status.Phase = "BootError"
			r.UpdateStatus(ctx, res)
			return err 
		}
		res.Status.MaterializedBoot = true
	}

	if res.Status.MaterializedBoot && res.Status.MaterializedMetadata {
		res.Status.Phase = "Bound"
	}

	r.UpdateStatus(ctx, res)
	return nil
}