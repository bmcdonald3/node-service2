# node-service (OpenCHAMI Node API Shim)

This service implements the profile-based node configuration and node API shim described in RFD #121. It provides a node-centric API surface that composes inventory, boot parameters, and cloud-init metadata composition, while decoupling configuration intent from inventory group membership.

## Capabilities

The current implementation provides the following operational features:

* **Resource Scaffolding:** API endpoints for `Node`, `NodeSet`, and `ProfileBinding` resources, backed by file storage and JSON validation.
* **NodeSet Resolution:** The `NodeSetReconciler` asynchronously resolves logical group intent (label selectors and explicit xnames) into static lists of node identifiers.
* **Profile Binding Materialization:** The `ProfileBindingReconciler` listens for binding events and executes a write-through strategy, automatically pushing profile assignments to the downstream configuration services (`metadata-service` and `boot-service`).
* **Composed Node API:** An endpoint (`/composed/nodes/{xname}`) designed to aggregate a node's inventory identity, effective profile, boot parameters, and cloud-init configuration into a single response.

## Testing Status

The following components have been tested and verified:

* **API Validation:** Integration tests verify that the Fabrica middleware correctly rejects malformed JSON and enforces required fields (e.g., rejecting a `ProfileBinding` missing a `profile`).
* **Reconciliation Logic:** Unit tests for both the `NodeSet` and `ProfileBinding` controllers confirm that state transformations and error queueing behave as expected when interfacing with mock clients.
* **Composed View:** Unit testing confirms the custom HTTP handler successfully routes and responds to specific `xname` requests outside of the generated REST routes.
* **End-to-End Loop:** A live test utilizing a secondary sidecar server verified that the asynchronous event bus correctly triggers the `ProfileBinding` worker queue, resulting in successful network requests to downstream services and subsequent database state updates (`MaterializedBoot: true`, `Phase: Bound`).

## Architecture

This service is built using the OpenCHAMI Fabrica framework, utilizing its in-memory event bus and reconciliation controller patterns to maintain declarative state.
