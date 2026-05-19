## Plan: Local Programmable Proxy in Go

Build a phased local programmable proxy in Go: deliver an MVP first (HTTP forward proxy + CONNECT tunneling + config-driven rules + SQLite logging + request/response interception), then add HTTPS interception (local CA MITM) as a second phase. This approach minimizes risk while preserving your end goal of full programmable control.

**Steps**
1. Phase 1 - Project foundation and runtime skeleton.
2. Create module/layout and startup wiring: go module, binary entrypoint, config loader, graceful shutdown, and central app container. This enables parallel development of proxy, rules, and persistence. *No dependency.*
3. Define core domain contracts up front: Rule, Matcher, Action, RequestContext, ResponseContext, Interceptor, RuleStore, LogStore. Include deterministic rule priority semantics (first match by priority, explicit default allow). *Depends on 1.*
4. Phase 2 - Proxy engine baseline (MVP transport path).
5. Implement forward proxy handler for HTTP requests and CONNECT tunneling for HTTPS passthrough. Add hop-by-hop header sanitization, upstream transport tuning, timeout limits, and per-request correlation IDs for traceability. *Depends on 2-3.*
6. Add pre-forward request interceptor chain with rule evaluation and early actions: allow, block (403/451), redirect, request-header modification, and request logging hooks. *Depends on 5.*
7. Phase 3 - Response interception and broad modification path.
8. Implement response pipeline with safe transform stages: decode (gzip/deflate/br if available), classify content type, transform body, re-encode, and header fixups (Content-Length and Transfer-Encoding). Apply fail-open behavior on uncertainty or transform failure. *Depends on 5; parallel with 10.*
9. Satisfy the any-content-type goal via policy tiers: text/json/xml/js/css/html in-place rewrite, known binary types pass-through with logging, unknown/multipart/streaming pass-through with explicit reason codes. This keeps behavior predictable without corrupting payloads. *Depends on 8.*
10. Add size/time guards for modification (max body bytes, transform timeout, memory budget) and enforce bypass-on-limit to prevent OOM/latency spikes. *Parallel with 8; depends on 5.*
11. Phase 4 - Persistence, configuration, and live rule updates.
12. Implement SQLite schema and repositories for traffic logs, rule evaluation events, and optional counters. Use batched or async write path to reduce request-path latency. *Depends on 2-3; parallel with 8-10.*
13. Implement config-file rule source (JSON/YAML) and hot-reload watcher with atomic swap (RWMutex + versioned snapshot). Validate rules before activation and preserve last-known-good snapshot on parse errors. *Depends on 3; parallel with 12.*
14. Merge config-file rules with runtime cache and expose deterministic startup order: load config -> validate -> activate -> start listener. *Depends on 12-13.*
15. Phase 5 - Focus mode and policy features.
16. Implement time-window policies (focus mode) as rule predicates with local timezone handling and schedule conflict resolution. Add policy tests for overnight windows and DST transitions. *Depends on 6 and 13.*
17. Add tracker/ad blocking presets (domain and URL patterns) as optional starter rule bundles loaded from config. *Depends on 13.*
18. Phase 6 - Hardening and observability.
19. Add structured logging, error taxonomy, and metrics counters (requests, blocks, rewrites, bypass reasons, upstream latency buckets). *Depends on 5-14.*
20. Add robustness protections: panic recovery in transform path, circuit-breaker style bypass on repeated transform failures, graceful shutdown drain, and bounded goroutine usage in tunnel copies. *Depends on 5-14.*
21. Phase 7 - Optional HTTPS interception (Phase 2 product milestone).
22. Implement local CA generation/loading, per-host leaf cert issuance + cache, TLS handshake split, and MITM request/response interception path. Keep this behind explicit config flags and per-domain allow/deny controls. *Depends on 5, 8, and 20.*
23. Add certificate lifecycle docs and trust-store setup scripts (Linux/macOS/Windows), plus guardrails for sensitive-domain bypass (banking, payments, auth providers). *Depends on 22.*

**Relevant files**
- /home/devansharora18/Documents/GitHub/allseer/README.md - Expand with architecture, setup, proxy configuration, and security model.
- /home/devansharora18/Documents/GitHub/allseer/go.mod - Module and dependency declarations.
- /home/devansharora18/Documents/GitHub/allseer/cmd/allseer/main.go - Process boot, config load, app wiring, shutdown signals.
- /home/devansharora18/Documents/GitHub/allseer/internal/config/config.go - Config schema, parsing, validation, hot-reload integration.
- /home/devansharora18/Documents/GitHub/allseer/internal/proxy/server.go - HTTP handler entrypoint and routing for HTTP vs CONNECT.
- /home/devansharora18/Documents/GitHub/allseer/internal/proxy/forward.go - Upstream forwarding path and header normalization.
- /home/devansharora18/Documents/GitHub/allseer/internal/proxy/tunnel.go - CONNECT tunnel lifecycle and io copy orchestration.
- /home/devansharora18/Documents/GitHub/allseer/internal/rules/engine.go - Rule priority, matching, action resolution.
- /home/devansharora18/Documents/GitHub/allseer/internal/rules/matcher.go - Domain/path/header/body/time-window predicates.
- /home/devansharora18/Documents/GitHub/allseer/internal/intercept/request.go - Request pre-processing and early response actions.
- /home/devansharora18/Documents/GitHub/allseer/internal/intercept/response.go - Decode/transform/re-encode pipeline.
- /home/devansharora18/Documents/GitHub/allseer/internal/transform/transformer.go - Content-type-aware body rewriting and safe bypass behavior.
- /home/devansharora18/Documents/GitHub/allseer/internal/storage/sqlite/db.go - DB init, migrations, connection settings.
- /home/devansharora18/Documents/GitHub/allseer/internal/storage/sqlite/rules_repo.go - Rule persistence/cache integration (if needed for snapshots).
- /home/devansharora18/Documents/GitHub/allseer/internal/storage/sqlite/log_repo.go - Traffic/event log writes and query helpers.
- /home/devansharora18/Documents/GitHub/allseer/internal/mitm/ca.go - CA load/generation and persistence.
- /home/devansharora18/Documents/GitHub/allseer/internal/mitm/intercept.go - TLS interception flow and cert cache usage.
- /home/devansharora18/Documents/GitHub/allseer/config/rules.example.yaml - Starter rules and focus-mode examples.
- /home/devansharora18/Documents/GitHub/allseer/test/integration/proxy_flow_test.go - End-to-end HTTP/CONNECT behavior.
- /home/devansharora18/Documents/GitHub/allseer/test/integration/rewrite_safety_test.go - Compression/content-length/transfer-encoding validation.

**Verification**
1. Unit tests: rule matcher/action resolution, config validation, content-type policy decisions, and response transform functions (including compression round-trips).
2. Integration tests: client -> proxy -> local upstream for allow/block/redirect/header rewrite; CONNECT passthrough behavior; hot-reload rule activation without restart.
3. Safety tests: large payload bypass, unknown content-type bypass, broken compressed payload fail-open, content-length correctness after rewrite, chunked response handling.
4. Persistence tests: SQLite migration startup, concurrent log inserts, query performance with indexes, and graceful shutdown flush behavior.
5. Focus mode tests: fixed-time and overnight windows, timezone correctness, deterministic action precedence with overlapping rules.
6. Optional MITM tests (phase 2): local CA trust flow, cert cache reuse, per-domain bypass enforcement, and TLS handshake error handling.
7. Manual validation: set OS proxy to localhost listener, browse mixed sites, verify blocked trackers, rewritten responses, and expected log entries.

**Decisions**
- Scope decision: Phased delivery (MVP first, HTTPS interception second).
- Rule management decision: Config-file based rules with hot reload as initial control plane.
- Modification decision: Any-content-type intent implemented with safe policy tiers and fail-open bypass for unsupported/high-risk payloads.
- Included in MVP: HTTP forward proxy, CONNECT tunnel passthrough, rule engine, request/response interception, SQLite logs, focus mode.
- Excluded from MVP: TLS decryption/MITM, websocket rewriting, HTTP/2-specific rewrites, complex streaming multipart transforms.

**Further Considerations**
1. Rule expression format: choose glob-only for simpler UX or glob+regex for advanced matching (recommended: glob-first with optional regex flag).
2. Log retention strategy: choose fixed-day TTL purge or size-based compaction (recommended: TTL + periodic vacuum).
3. Phase 2 launch gate: require explicit acceptance criteria on security docs and trust-store automation before enabling MITM by default.
