# node-service (OpenCHAMI Node API Shim)

This service implements the profile-based node configuration and node API shim described in RFD #121. It provides a node-centric API surface that composes inventory, boot parameters, and cloud-init metadata composition, while decoupling configuration intent from inventory group membership.

## Capabilities

The current implementation provides the following operational features:

* **Resource Scaffolding:** API endpoints for `Node`, `NodeSet`, and `ProfileBinding` resources, backed by file storage and JSON validation.
* **NodeSet Resolution:** The `NodeSetReconciler` asynchronously resolves logical group intent (label selectors and explicit xnames) into static lists of node identifiers.
* **Profile Binding Materialization:** The `ProfileBindingReconciler` listens for binding events and executes a write-through strategy, automatically pushing profile assignments to the downstream configuration services (`metadata-service` and `boot-service`).
* **Composed Node API:** An endpoint (`/composed/nodes/{xname}`) designed to aggregate a node's inventory identity, effective profile, boot parameters, and cloud-init configuration into a single response.

## Usage and Live Testing

To observe the asynchronous reconciliation and write-through materialization in action, follow these steps to run the primary API and the sidecar mock server locally.

1. **Start the servers:**

```bash
pkill -f "go run ./cmd/" || true
go run ./cmd/server/ &
go run ./cmd/composed-api/ &
```

2. **Create a new ProfileBinding:**
This simulates an admin assigning a new profile to a node. Notice the initial status fields are `false`.

```bash
curl -X POST http://localhost:8080/profilebindings \
  -H "Content-Type: application/json" \
  -d '{
    "apiVersion":"v1",
    "kind":"ProfileBinding",
    "metadata":{"name":"test-binding-live"},
    "spec":{
      "targetKind":"Node",
      "targetName":"x1000c0s1b0n0",
      "profile":"compute-new"
    }
  }'
```

*Initial Response:*
```json
{
  "apiVersion":"v1",
  "kind":"ProfileBinding",
  "metadata":{
    "name":"test-binding-live",
    "uid":"profilebinding-67a06bee",
    "createdAt":"2026-03-27T11:56:30.83397-07:00",
    "updatedAt":"2026-03-27T11:56:30.83397-07:00"
  },
  "spec":{
    "targetKind":"Node",
    "targetName":"x1000c0s1b0n0",
    "profile":"compute-new"
  },
  "status":{
    "materializedBoot":false,
    "materializedMetadata":false
  }
}
```

3. **Verify Materialization:**
Wait a few seconds for the background `ProfileBindingReconciler` to process the event, make the HTTP calls to the downstream mock server (port 8081), and update the storage state.

```bash
curl http://localhost:8080/profilebindings
```

*Eventual State Output:*
```json
[
  {
    "apiVersion":"v1",
    "kind":"ProfileBinding",
    "metadata":{
      "name":"test-binding-live",
      "uid":"profilebinding-67a06bee",
      "createdAt":"2026-03-27T11:56:30.83397-07:00",
      "updatedAt":"2026-03-27T11:56:30-07:00"
    },
    "spec":{
      "targetKind":"Node",
      "targetName":"x1000c0s1b0n0",
      "profile":"compute-new"
    },
    "status":{
      "materializedBoot":true,
      "materializedMetadata":true,
      "phase":"Bound"
    }
  }
]
```

## Architecture

This service is built using the OpenCHAMI Fabrica framework, utilizing its in-memory event bus and reconciliation controller patterns to maintain declarative state.

## Next Steps for Agent Hand-off

The foundational wiring and Fabrica generation are complete. The next AI agent taking over this workspace should proceed with the following specific implementation tasks:

1. **Replace Reconciler Mock Clients:** * Update `pkg/reconcilers/nodeset_reconciler.go` to replace `MockSMDClient` with an actual HTTP client capable of querying the OpenCHAMI SMD service.
   * Update `pkg/reconcilers/profilebinding_reconciler.go` to remove the hardcoded `localhost:8081` URLs and implement proper routing and authentication to the real `metadata-service` and `boot-service`.
2. **Implement Dynamic Composed API Handler:**
   * Update `internal/api/composed_node.go` (`GetComposedNodeHandler`). Remove the static mock response. Implement concurrent HTTP GET requests to SMD, `boot-service`, and `metadata-service` for the requested `xname`. Merge the responses into the `ComposedNodeResponse` struct.
3. **Enforce Unknown Profile Policy:**
   * Implement a strict 404 fallback policy. If a requested profile does not exist during binding validation or composed resolution, it must explicitly return an error rather than falling back to a default.
4. **Implement Phase 5 (Campaigns):**
   * Use `fabrica add resource Campaign` to scaffold the Campaign resource.
   * Create a `CampaignReconciler` in `pkg/reconcilers/` to manage canary rollouts, applying `ProfileBindings` to batches of nodes over time based on the Campaign Spec.