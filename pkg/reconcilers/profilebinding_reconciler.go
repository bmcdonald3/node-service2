package reconcilers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/openchami/fabrica/pkg/reconcile"
	v1 "github.com/user/node-service/apis/node.openchami.io/v1"
)

type DownstreamClient interface {
	ApplyProfile(ctx context.Context, targetName, profile string) error
}

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

type ProfileBindingReconciler struct {
	reconcile.BaseReconciler
	MetadataClient DownstreamClient
	BootClient     DownstreamClient
}

func (r *ProfileBindingReconciler) Reconcile(ctx context.Context, resource interface{}) (reconcile.Result, error) {
	if r.MetadataClient == nil {
		r.MetadataClient = &HTTPDownstreamClient{BaseURL: "http://localhost:8081", HTTPClient: &http.Client{}}
	}
	if r.BootClient == nil {
		r.BootClient = &HTTPDownstreamClient{BaseURL: "http://localhost:8081", HTTPClient: &http.Client{}}
	}

	binding := resource.(*v1.ProfileBinding)

	if !binding.Status.MaterializedMetadata {
		err := r.MetadataClient.ApplyProfile(ctx, binding.Spec.TargetName, binding.Spec.Profile)
		if err != nil {
			binding.Status.Phase = "MetadataError"
			r.UpdateStatus(ctx, binding)
			return reconcile.Result{RequeueAfter: 10 * time.Second}, err
		}
		binding.Status.MaterializedMetadata = true
	}

	if !binding.Status.MaterializedBoot {
		err := r.BootClient.ApplyProfile(ctx, binding.Spec.TargetName, binding.Spec.Profile)
		if err != nil {
			binding.Status.Phase = "BootError"
			r.UpdateStatus(ctx, binding)
			return reconcile.Result{RequeueAfter: 10 * time.Second}, err
		}
		binding.Status.MaterializedBoot = true
	}

	if binding.Status.MaterializedBoot && binding.Status.MaterializedMetadata {
		binding.Status.Phase = "Bound"
	}

	r.UpdateStatus(ctx, binding)
	return reconcile.Result{}, nil
}

func (r *ProfileBindingReconciler) GetResourceKind() string {
	return "ProfileBinding"
}