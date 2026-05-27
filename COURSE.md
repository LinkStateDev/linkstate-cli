# LinkState — Course & Curriculum Plan

## TL;DR

**LinkState** — hands-on платформа для network engineers, infrastructure engineers и SRE, которые хотят перейти от ad-hoc scripts к production-grade NetOps tooling на Go.

Главная идея курса:

> **Build real NetOps tools one tested checkpoint at a time — in your editor, against real network labs.**

Педагогика — Boot.dev-style progression: маленькие проверяемые шаги, быстрый feedback, постоянный progress. Но реализация не browser IDE. Код пишется локально в editor студента, а `lst` управляет workspace, checkpoints, tests, submit, progress и later labs.

Структура продукта:

```text
LinkState
├── Track 1: Foundation — Go for NetOps
├── Track 2: Build Your NetOps Framework      ← flagship / main paid track
├── Track 3: Network Observability in Go
└── Track 4: Protocol Internals
```

Главный paid artifact — не config push, а **state testing / network validation framework**. Это безопаснее, быстрее даёт value, проще тестируется, и потом становится core dependency для deploy, observability и device lifecycle automation.

Ключевая curriculum spine:

```text
model network data
  → fetch network state
  → test expected vs actual state
  → render intended config
  → deploy safely
  → provision/deprovision devices
  → observe continuously
  → understand protocols underneath
```

---

# 1. Core Positioning

## One-liner

> **Build production-grade NetOps tools in Go — with tests, real labs, and production patterns.**

## Longer positioning

LinkState teaches network engineers, infrastructure engineers, and SREs to build real network automation tooling in Go. Students work in their own editor, run guided checks through `lst`, and eventually validate their code against Containerlab/netlab-based labs.

The platform is not a generic Go course and not cert-prep. It is an artifact-driven path from small local tools to a reusable NetOps framework.

## Internal shorthand

Internally, LinkState can be thought of as:

> **Boot.dev cadence for NetOps tooling, but local-first and lab-backed.**

Publicly, avoid saying “Boot.dev for network engineers” as the main positioning. It makes the product derivative and may attract generic beginner-programming expectations.

Better public framing:

```text
Build real NetOps tools in Go.
In your editor. With tests. Against real labs.
```

## What LinkState is not

LinkState is not:

- a browser coding game;
- a CCNP/CCIE lab prep product;
- a Cisco/Juniper/Arista vendor-cert course;
- a generic Go course with network-themed examples;
- a collection of blog posts;
- a video course;
- a cloud lab rental platform;
- an AI-for-networking hype product.

---

# 2. Curriculum Model

## Hierarchy

```text
Track → Module → Lesson → Checkpoint
```

## Definitions

### Track

A large learning journey with a coherent outcome.

Examples:

```text
Foundation — Go for NetOps
Build Your NetOps Framework
Network Observability in Go
Protocol Internals
```

### Module

A substantial artifact or capability inside a track.

Examples:

```text
IP Toolkit
State Tester
Config Renderer
Device Lifecycle
Prometheus Exporter
DNS Parser
```

### Lesson

A connected explanatory unit, usually 20–45 minutes of student time in flagship tracks and 10–25 minutes in Foundation.

A lesson teaches one concept or implementation slice.

### Checkpoint

A small, testable increment inside a lesson. Usually 5–12 minutes of work.

A checkpoint has:

- a clear task;
- target file(s);
- expected behavior;
- private checker coverage;
- optional visible example tests;
- progress state;
- `lst test` / `lst submit` integration.

## Important principle

> **Lessons teach concepts. Checkpoints verify progress. Modules produce artifacts.**

Or:

> **A checkpoint is a small, testable increment in a larger artifact.**

---

# 3. Learning Experience

## Desired student loop

The ideal loop is:

```bash
lst start 00_go_for_netops_primer
lst test
lst submit --next
lst test
lst submit --next
...
```

Not:

```bash
lst fetch lesson-01
cd lesson-01
code .
open browser
...
lst fetch lesson-02
cd ../lesson-02
code .
...
```

The product should minimize navigation friction. Small checkpoints only work if moving between them is nearly frictionless.

## Stable workspace principle

Students should not open a new directory for every lesson.

Foundation:

```text
~/linkstate/foundation/
├── go.mod
├── lessons/
├── internal/
├── CURRENT.md
└── .linkstate/
```

Flagship module:

```text
~/linkstate/netops-framework/state-tester/
├── go.mod
├── cmd/
├── internal/
├── testdata/
├── CURRENT.md
└── .linkstate/
```

One workspace per Foundation track or per flagship artifact/module.

## Core `lst` commands

Minimum user-facing workflow:

```bash
lst init
lst start <track-or-module>
lst status
lst open
lst next
lst test
lst submit
lst resume
lst doctor
```

Later:

```bash
lst test --watch
lst lab up
lst lab down
lst lab reset
lst solution
lst hint
```

## `lst start`

Starts a track or module.

Example:

```bash
lst start 00_go_for_netops_primer
lst start netops-framework/state-tester
```

Responsibilities:

- create workspace;
- download starter skeleton;
- write `.linkstate/manifest.yml`;
- write initial `CURRENT.md`;
- open lesson page or print URL;
- optionally open editor;
- set current checkpoint.

## `lst next`

Moves to the next checkpoint or lesson.

Responsibilities:

- inspect progress;
- determine next checkpoint;
- check access/entitlement;
- update `CURRENT.md`;
- print exact file and line;
- open browser if configured;
- optionally focus editor;
- never overwrite user code silently.

Example output:

```text
Next checkpoint:
  State Tester / Parse FRR BGP Summary / Checkpoint 2

Task:
  Parse non-Established BGP peers.

Read:
  https://linkstate.dev/tracks/netops-framework/state-tester/parse-frr#checkpoint-2

Edit:
  internal/parser/frr.go:42

Run:
  lst test
```

## `lst test`

Runs the current checkpoint checker.

During development:

```bash
lst test
```

Expected behavior:

- run visible/sanity tests if present;
- run private compiled checker;
- run current checkpoint checks;
- optionally run regression checks for previous checkpoints;
- for lab checkpoints, run lab assertions;
- produce friendly, domain-specific feedback.

## `lst submit`

Submission should run checks again.

```bash
lst submit
lst submit --next
```

Responsibilities:

- run local checker;
- validate current checkpoint;
- sync progress if authenticated;
- unlock next checkpoint;
- optionally advance immediately with `--next`.

## `lst resume`

Restores learning context after a break.

```bash
lst resume
```

Responsibilities:

- show current track/module/lesson/checkpoint;
- open lesson URL or print it;
- show workspace;
- show target file and line;
- optionally open editor.

## Editor integration

Editor integration should be a progressive enhancement, not a requirement.

Core behavior must work in plain terminal:

```text
Edit:
  ~/linkstate/foundation/lessons/04-peer-state/parser.go:14
```

Supported levels:

```text
Level 0: print URL/file/line
Level 1: open browser
Level 2: VS Code-style `code -r -g` support
Level 3: custom editor command templates
Level 4: optional extensions/plugins later
```

Do not block learning if editor auto-open fails.

## `CURRENT.md`

Each workspace should contain a generated `CURRENT.md` file.

Example:

```markdown
# Current LinkState checkpoint

Track:
Build Your NetOps Framework

Module:
State Tester

Lesson:
Parse FRR BGP summary

Checkpoint:
Parse non-Established peers

Read:
https://linkstate.dev/tracks/netops-framework/state-tester/parse-frr#checkpoint-2

Edit:
internal/parser/frr.go:42

Goal:
Return State="Active" and Prefixes=0 when State/PfxRcd contains a state name.

Run:
lst test

Submit:
lst submit --next
```

This provides a universal fallback even without IDE integration.

---

# 4. Public Reading vs Gated Action

## Public reading is marketing

Public pages should include:

- conceptual explanation;
- diagrams;
- command examples;
- partial code snippets;
- expected output;
- task description;
- CTA to actually build.

Public pages prove expertise and attract search/community traffic.

## Registered/paid action is the product

Gated parts should include:

- `lst start` / `lst fetch` starter repos;
- private compiled checkers;
- lab topology bundles;
- submit/progress sync;
- reference solutions;
- premium modules;
- support bundle tooling;
- certificates later.

Core distinction:

> **The lesson explains the concept. The platform provides the environment, verification, progression, and artifact.**

Or shorter:

> **The lesson is public. The build system is gated.**

## Why this matters

If everything is behind a paywall, cold technical users have no reason to trust the product.

If everything is public, the product looks like a blog with a CLI wrapper.

The correct boundary is:

```text
Understanding → public
Execution / verification / progress → gated
```

---

# 5. Private Checkers and Anti-Piracy Model

## Current direction

Use compiled checker binaries, not public full `_test.go` files, as the primary grading mechanism.

This is compatible with checkpoints.

Checkpoint is a learning/progress concept. It does not require shipping source tests to students.

## Principle

> **Checkpoints are public learning milestones. Checkers are private executable graders.**

And:

> **Visible tests document the contract. Private checkers grade the solution.**

## Starter repos may include visible tests

Visible tests should be small examples, not full grading logic.

They help students understand the API and run basic `go test` locally.

Private compiled checkers cover:

- edge cases;
- hidden inputs;
- malformed data;
- lab state;
- regression from previous checkpoints;
- JSON schema / exit codes;
- production gotchas.

## Checker delivery

Prefer:

```text
one checker binary per module per version per OS/arch
```

Example:

```text
state-tester-checker_2026.05.1_darwin_arm64
state-tester-checker_2026.05.1_linux_amd64
state-tester-checker_2026.05.1_windows_amd64.exe
```

The checker accepts:

```bash
--workspace <path>
--checkpoint <id>
--format json
```

Example:

```bash
state-tester-checker \
  --workspace ~/linkstate/netops-framework/state-tester \
  --checkpoint parser.non-established
```

## White-box and black-box checks

Private checker can use both.

### White-box style

For early checkpoints, checker can inject temporary tests into a temp copy of workspace:

```text
copy workspace to temp dir
write private temporary _test.go
run go test
delete temp dir
```

Do not write private tests into student workspace.

### Black-box style

For later checkpoints, checker can run the student CLI:

```bash
go run ./cmd/netops check --inventory testdata/inventory.yml --checks testdata/checks.yml --json
```

Then verify:

- stdout/stderr;
- exit code;
- JSON schema;
- behavior under failure cases.

## `lst submit` trust model

For MVP:

```text
lst submit → run local private checker → sync progress
```

This is enough for early product validation.

For certificates/team reporting later:

```text
lst submit → upload solution snapshot → server-side validation
```

Do not treat local checker output as high-stakes proof forever.

## Anti-piracy goal

Do not try to build perfect DRM.

Reasonable v1 goals:

```text
1. Anonymous users cannot download starter repos.
2. Free users cannot download premium module assets.
3. Private tests are not shipped as normal files.
4. Reference solutions are gated.
5. Lab bundles are gated.
6. Official product is easier and more current than a pirated copy.
```

---

# 6. Track Overview

```text
Track 1: Foundation — Go for NetOps
  Goal: teach Go primitives through small network tools.
  Lab: no Docker, no real devices.
  Access: free registered.

Track 2: Build Your NetOps Framework
  Goal: build a reusable state/config/deploy/lifecycle framework.
  Lab: FRR/containerlab/netlab.
  Access: flagship paid track, with free preview.

Track 3: Network Observability in Go
  Goal: turn checks, events, and telemetry into metrics and alerts.
  Lab: FRR + optional telemetry components.
  Access: paid elective.

Track 4: Protocol Internals
  Goal: understand protocols by building parsers, encoders, and state machines.
  Lab: mostly fixtures/pcaps/binary data; optional lab later.
  Access: paid elective.
```

---

# 7. Track 1 — Foundation: Go for NetOps

## Goal

Foundation teaches enough Go and software structure for students to start the flagship track.

It should feel like:

> “I am already building useful network tools.”

Not:

> “I am doing generic beginner Go exercises before the real course starts.”

## Constraints

Foundation should be low-friction:

```text
No Docker.
No Containerlab.
No SSH to real devices.
No NetBox.
No vendor images.
```

Required:

```text
Go + lst + local editor
```

## Foundation artifacts

By the end, the student has built:

```text
netipcalc       — small IP/CIDR toolkit
config-fetcher  — concurrent mock config fetcher
mini-netcheck   — offline expected-vs-actual state checker
```

## Suggested structure

```text
Foundation — Go for NetOps
├── 00. Introduction to LinkState Workflow
├── 01. IP Toolkit
├── 02. Inventory & Device Model
├── 03. Config Fetcher with MockDriver
├── 04. Offline State Assertions
└── 05. Foundation Capstone: Mini NetCheck
```

## Module 00 — Introduction to LinkState Workflow

### Goal

Get the student through the first `lst` loop quickly.

### Concepts

- workspace layout;
- lessons vs checkpoints;
- `lst test`;
- `lst submit --next`;
- `CURRENT.md`;
- how public lesson + local code work together.

### Example flow

```bash
lst init
lst start 00_go_for_netops_primer
lst test
lst submit --next
```

### Checkpoints

```text
1. Install and initialize lst.
2. Start Foundation workspace.
3. Run first checker.
4. Submit first checkpoint.
```

No significant coding yet. The goal is activation.

---

## Module 01 — IP Toolkit

### Goal

Teach parsing, validation, errors, bitwise operations, CLI output, and small Go functions through IP/CIDR tasks.

### Artifact

```bash
netipcalc
```

Example commands:

```bash
netipcalc info 10.0.1.15/24
netipcalc contains 10.0.0.0/16 10.0.1.15
netipcalc overlap 10.0.0.0/24 10.0.0.128/25
netipcalc json 10.0.1.15/24
```

### Lessons and checkpoints

#### Lesson 01: Parse IP and CIDR input

Checkpoints:

```text
1. Parse IPv4 octets.
2. Validate octet range.
3. Parse CIDR prefix.
4. Return structured errors for bad input.
```

#### Lesson 02: Network math

Checkpoints:

```text
1. Convert IPv4 to uint32.
2. Calculate network address.
3. Calculate broadcast address.
4. Calculate usable host count.
```

#### Lesson 03: Prefix relationships

Checkpoints:

```text
1. Implement contains(ip).
2. Implement contains(prefix).
3. Implement overlap(prefixA, prefixB).
4. Add table output.
```

#### Lesson 04: CLI and JSON output

Checkpoints:

```text
1. Add subcommands.
2. Add stable human output.
3. Add JSON output.
4. Add exit codes for invalid input.
```

### Skills taught

- functions;
- structs;
- errors;
- parsing;
- bitwise logic;
- CLI shape;
- JSON output;
- small artifact completion.

---

## Module 02 — Inventory & Device Model

### Goal

Introduce the idea that network automation starts with a model of devices, roles, platforms and management addresses.

This becomes the foundation for every later tool.

### Artifact

```bash
invcheck
```

Example:

```bash
invcheck --inventory inventory.yml
invcheck --inventory inventory.yml --json
```

### Example inventory

```yaml
devices:
  - name: leaf1
    role: leaf
    site: lab1
    platform: frr
    mgmt_ip: 172.20.20.11

  - name: spine1
    role: spine
    site: lab1
    platform: frr
    mgmt_ip: 172.20.20.21
```

### Lessons and checkpoints

#### Lesson 01: Model a device

Checkpoints:

```text
1. Define Device struct.
2. Add Role, Site, Platform, MgmtIP fields.
3. Add validation for required fields.
4. Add MgmtURL or Address helper method.
```

#### Lesson 02: Load inventory from YAML/JSON

Checkpoints:

```text
1. Decode inventory file.
2. Validate duplicate device names.
3. Validate duplicate management IPs.
4. Return structured validation errors.
```

#### Lesson 03: Filter and group devices

Checkpoints:

```text
1. Filter devices by role.
2. Filter devices by platform.
3. Group devices by site.
4. Produce stable sorted output.
```

### Skills taught

- structs;
- methods;
- slices/maps;
- YAML/JSON decoding;
- validation;
- deterministic output;
- source-of-truth thinking.

---

## Module 03 — Config Fetcher with MockDriver

### Goal

Introduce driver interfaces and concurrent collection without requiring real devices.

### Artifact

```bash
config-fetcher
```

Example:

```bash
config-fetcher --inventory inventory.yml --out configs/
```

### Core interface

```go
type Driver interface {
    GetConfig(ctx context.Context, device Device) (string, error)
}
```

### Lessons and checkpoints

#### Lesson 01: Define a driver interface

Checkpoints:

```text
1. Define Driver interface.
2. Implement MockDriver.
3. Fetch config for one device.
4. Return wrapped errors.
```

#### Lesson 02: Fetch configs for inventory

Checkpoints:

```text
1. Iterate over all devices.
2. Save configs to disk.
3. Create stable output paths.
4. Handle partial failures.
```

#### Lesson 03: Concurrent fetch

Checkpoints:

```text
1. Start goroutine per device.
2. Collect results through channel.
3. Add sync.WaitGroup.
4. Preserve deterministic report order.
```

#### Lesson 04: Context and timeout

Checkpoints:

```text
1. Pass context into Driver.
2. Add per-device timeout.
3. Classify timeout errors.
4. Return nonzero exit code on failed fetches.
```

### Skills taught

- interfaces;
- mocks;
- goroutines;
- channels;
- WaitGroup;
- context;
- partial failure;
- file output.

---

## Module 04 — Offline State Assertions

### Goal

Introduce the central NetOps idea:

> expected state vs actual state

This is the conceptual bridge to the flagship State Tester.

### Artifact

```bash
mini-netcheck
```

Example:

```bash
mini-netcheck --expected expected.yml --actual fixtures/
mini-netcheck --expected expected.yml --actual fixtures/ --json
```

### Example checks

```yaml
checks:
  - name: leaf1 uplink should be up
    type: interface_state
    device: leaf1
    interface: eth1
    expected: up

  - name: leaf1 should have service route
    type: route_exists
    device: leaf1
    prefix: 10.10.0.0/16
```

### Lessons and checkpoints

#### Lesson 01: Define check model

Checkpoints:

```text
1. Define Check struct.
2. Define CheckResult struct.
3. Load checks from YAML.
4. Validate unknown check types.
```

#### Lesson 02: Evaluate interface state from fixtures

Checkpoints:

```text
1. Parse interface fixture.
2. Find interface by name.
3. Compare expected vs actual state.
4. Return clear failure reason.
```

#### Lesson 03: Evaluate route existence from fixtures

Checkpoints:

```text
1. Parse route fixture.
2. Implement route_exists.
3. Implement route_absent.
4. Add structured discrepancy output.
```

#### Lesson 04: Report results

Checkpoints:

```text
1. Human table output.
2. JSON output.
3. Stable exit codes.
4. Summary counts.
```

### Skills taught

- domain modeling;
- expected/actual comparison;
- structured results;
- exit codes;
- foundation of validation tooling.

---

## Module 05 — Foundation Capstone: Mini NetCheck

### Goal

Combine Foundation concepts into a single offline tool.

### Artifact

```bash
mini-netcheck
```

Final behavior:

```bash
mini-netcheck --inventory inventory.yml --checks checks.yml --fixtures fixtures/ --json
```

### Capstone requirements

The tool should:

- load inventory;
- load checks;
- collect actual state from fixtures via a `Collector`/`Driver` interface;
- run checks;
- produce human output;
- produce JSON output;
- return stable exit codes;
- handle partial failures.

### Checkpoints

```text
1. Wire inventory + checks together.
2. Implement fixture collector.
3. Run all checks.
4. Add JSON output.
5. Add exit codes.
6. Produce final report.
```

### Bridge to flagship

The last lesson should explicitly say:

> In Foundation, your collector reads fixtures. In the NetOps Framework track, the same architecture talks to real FRR devices in a local lab.

Optional teaser:

```bash
lst lab up
netops check --inventory lab.yml --checks checks.yml
```

But Docker/lab should not be required for Foundation completion.

---

# 8. Track 2 — Build Your NetOps Framework

## Goal

This is the flagship track.

Students build a reusable NetOps framework with:

```text
inventory
  → drivers / collectors
  → state checks
  → reports / CI
  → config rendering
  → safe deploy
  → source-of-truth integration
  → device lifecycle automation
```

The flagship should feel like:

> “I am building the automation framework my network team should have.”

## Access

Recommended:

```text
Free registered:
  Introduction + first lab / first state checks preview

Paid:
  full state tester, reports, renderer, deployer, source-of-truth, lifecycle, capstone
```

## Lab

Primary lab:

```text
FRR-based spine/leaf topology via netlab/containerlab
```

Initial topology:

```text
spine1      spine2
  | \      / |
  |  \    /  |
leaf1    leaf2
```

Later lifecycle modules can add or simulate `leaf3`.

Use netlab as the underlying lab foundation. `lst lab` wraps it for student UX.

## Flagship artifacts

By the end, students build something like:

```bash
netops check
netops render
netops plan
netops apply
netops device provision
netops device deprovision
```

## Suggested structure

```text
Build Your NetOps Framework
├── 00. Framework Introduction
├── 01. Lab & Real Device Access
├── 02. Inventory, Drivers & Collectors
├── 03. State Tester
├── 04. Reports & CI
├── 05. Config Renderer
├── 06. Config Deployer
├── 07. Source of Truth Integration
├── 08. Device Lifecycle: ZTP, Provisioning & Deprovisioning
└── 09. Capstone: NetOps Controller
```

---

## Module 00 — Framework Introduction

### Goal

Show the destination before implementation begins.

### Concepts

- why scripts become brittle;
- framework vs one-off tool;
- expected vs actual state;
- read-only automation before mutation;
- driver abstraction;
- lab-backed learning;
- how Foundation concepts transfer.

### Target architecture

```text
Source of Truth / Inventory
          ↓
Drivers / Collectors
          ↓
Parsers / Normalizers
          ↓
Checks / Validators
          ↓
Reports / CI
          ↓
Render / Plan / Deploy
          ↓
Post-check / Rollback / Audit
```

### Checkpoints

```text
1. Start NetOps Framework workspace.
2. Inspect lab topology.
3. Run first framework command against fixtures.
4. Understand module artifact layout.
```

---

## Module 01 — Lab & Real Device Access

### Goal

Introduce the real lab without overwhelming the student.

### Artifact

```bash
netops show
```

Example:

```bash
lst lab up
netops show --device leaf1 --command "show ip bgp summary"
netops show --all --command "show ip route"
```

### Lessons and checkpoints

#### Lesson 01: Bring up the lab

Checkpoints:

```text
1. Run lst lab up.
2. Wait for lab health checks.
3. Inspect nodes.
4. Run lst doctor for lab status.
```

#### Lesson 02: Execute show command on one device

Checkpoints:

```text
1. Define Device connection info.
2. Connect to one FRR device.
3. Run show command.
4. Return stdout/stderr/errors clearly.
```

#### Lesson 03: Parallel show command

Checkpoints:

```text
1. Run command across all devices.
2. Add timeout per device.
3. Handle partial failures.
4. Print deterministic report.
```

### Why this module matters

This is the first “same code, real network” moment. It should be free preview or at least partially free.

---

## Module 02 — Inventory, Drivers & Collectors

### Goal

Upgrade Foundation inventory and mock driver into real lab-aware drivers.

### Artifacts

```bash
netops inventory validate
netops collect
```

### Core interfaces

```go
type Driver interface {
    Show(ctx context.Context, device Device, command string) (string, error)
    GetConfig(ctx context.Context, device Device) (string, error)
}

type Collector interface {
    Collect(ctx context.Context, device Device) (DeviceState, error)
}
```

### Lessons and checkpoints

#### Lesson 01: Framework inventory

Checkpoints:

```text
1. Load lab inventory.
2. Validate roles/platforms.
3. Attach connection parameters.
4. Filter target devices.
```

#### Lesson 02: FRR driver

Checkpoints:

```text
1. Implement Show().
2. Implement GetConfig().
3. Add context timeout.
4. Wrap connection errors.
```

#### Lesson 03: Collector abstraction

Checkpoints:

```text
1. Define DeviceState.
2. Collect interface output.
3. Collect route output.
4. Collect BGP summary output.
```

#### Lesson 04: Normalize raw outputs

Checkpoints:

```text
1. Store raw command output.
2. Normalize into typed state.
3. Preserve raw output for debugging.
4. Add partial-state reporting.
```

---

## Module 03 — State Tester

## Goal

Build the first main paid artifact:

```bash
netops check
```

This is the center of the whole curriculum.

State Tester validates expected network state against actual network state from the lab.

## Why first

State testing should come before config deployment because it is:

- read-only;
- safer;
- easier to support;
- easy to demonstrate;
- directly useful in CI/post-change checks;
- reused by deployer, lifecycle and observability modules.

## Command

```bash
netops check --inventory lab.yml --checks checks.yml
netops check --inventory lab.yml --checks checks.yml --json
```

## Initial check types

Start with three high-value checks:

```text
bgp_peer_state
route_exists
interface_state
```

Then add absence checks:

```text
route_absent
bgp_peer_absent
interface_unused
config_line_absent
```

Later optional checks:

```text
prefix_count
bgp_prefix_received
lldp_neighbor
interface_mtu
bfd_session_state
no_routes_from_peer
```

## Example checks file

```yaml
checks:
  - name: leaf1 peers with spine1
    type: bgp_peer_state
    device: leaf1
    peer: 10.0.0.1
    expected: established

  - name: leaf1 has service route
    type: route_exists
    device: leaf1
    prefix: 10.10.0.0/16

  - name: leaf1 uplink is up
    type: interface_state
    device: leaf1
    interface: eth1
    expected: up
```

## Lessons and checkpoints

### Lesson 01: Check model

Checkpoints:

```text
1. Define CheckSpec.
2. Define CheckResult.
3. Load checks from YAML.
4. Validate unsupported check types.
```

### Lesson 02: Interface state check

Checkpoints:

```text
1. Parse interface state from FRR output.
2. Implement interface_state.
3. Return pass/fail result.
4. Include actual vs expected values.
```

### Lesson 03: Route exists check

Checkpoints:

```text
1. Parse route table output.
2. Implement route_exists.
3. Implement route_absent.
4. Add prefix normalization.
```

### Lesson 04: BGP peer state check

Checkpoints:

```text
1. Parse show ip bgp summary row.
2. Handle Established sessions.
3. Handle Active/Idle/Connect sessions.
4. Implement bgp_peer_state.
5. Return useful discrepancy reasons.
```

### Lesson 05: Multi-device execution

Checkpoints:

```text
1. Group checks by device.
2. Collect state per device.
3. Run checks concurrently.
4. Preserve deterministic result order.
```

### Lesson 06: Partial failures

Checkpoints:

```text
1. Continue if one device fails.
2. Mark checks as unknown when state is missing.
3. Distinguish check failure from tool failure.
4. Add timeout classification.
```

## Artifact outcome

By the end:

```bash
netops check --inventory lab.yml --checks checks.yml
```

Outputs:

```text
✓ leaf1 bgp peer spine1 established
✓ leaf1 route 10.10.0.0/16 exists
✗ leaf2 eth1 expected up, got down

Result: 1 failed, 2 passed
```

---

## Module 04 — Reports & CI

### Goal

Turn State Tester from a script into a production-grade tool.

### Why separate module

Report quality, JSON output and exit codes are not polish. They are what makes a tool usable in CI and automation pipelines.

### Artifacts

```bash
netops check --json
netops check --junit
netops check --ci
```

### Exit codes

```text
0 = all checks passed
1 = checks failed
2 = tool/runtime error
```

### Lessons and checkpoints

#### Lesson 01: Human report

Checkpoints:

```text
1. Add stable table output.
2. Add summary counts.
3. Sort results deterministically.
4. Include failure reasons.
```

#### Lesson 02: JSON output

Checkpoints:

```text
1. Define JSON schema.
2. Marshal check results.
3. Include metadata and duration.
4. Make output stable for CI.
```

#### Lesson 03: Exit codes

Checkpoints:

```text
1. Exit 0 on pass.
2. Exit 1 on check failure.
3. Exit 2 on runtime/tool failure.
4. Test mixed partial failure behavior.
```

#### Lesson 04: CI integration

Checkpoints:

```text
1. Add make target.
2. Add GitHub Actions or generic CI example.
3. Store JSON artifact.
4. Fail pipeline on network check failure.
```

#### Lesson 05: Audit report

Checkpoints:

```text
1. Write report file.
2. Include inventory/checks version hash.
3. Include lab/tool version.
4. Produce reproducible audit output.
```

---

## Module 05 — Config Renderer

### Goal

Generate intended device configs from typed data.

This should come after State Tester so the student already understands inventory, device models and expected state.

### Artifact

```bash
netops render
```

Example:

```bash
netops render --inventory lab.yml --vars fabric.yml --out rendered/
```

### Concepts

- intended config;
- typed data before templates;
- deterministic rendering;
- per-platform templates;
- validation before render;
- bootstrap vs full config;
- removal config later for lifecycle.

### Lessons and checkpoints

#### Lesson 01: Render from device model

Checkpoints:

```text
1. Define render input model.
2. Render hostname and interfaces.
3. Render loopbacks.
4. Write per-device config files.
```

#### Lesson 02: BGP config templates

Checkpoints:

```text
1. Add ASN and neighbors.
2. Render BGP router stanza.
3. Render per-neighbor config.
4. Validate missing peer data.
```

#### Lesson 03: Deterministic output

Checkpoints:

```text
1. Sort interfaces.
2. Sort neighbors.
3. Normalize whitespace.
4. Compare rendered output with golden file.
```

#### Lesson 04: Multi-platform template boundary

Checkpoints:

```text
1. Add platform field.
2. Choose template by platform.
3. Return error on unsupported platform.
4. Keep common data model platform-neutral.
```

#### Lesson 05: Bootstrap vs intended config

Checkpoints:

```text
1. Render bootstrap config.
2. Render full intended config.
3. Explain difference between bootstrap and full config.
4. Add command selector: bootstrap/full.
```

### Output example

```text
rendered/
├── leaf1.conf
├── leaf2.conf
├── spine1.conf
└── spine2.conf
```

---

## Module 06 — Config Deployer

### Goal

Apply rendered configs safely.

This module should be mutation-aware and conservative. It must build on State Tester and Config Renderer.

### Artifact

```bash
netops plan
netops apply
netops rollback
```

### Concepts

- dry-run by default;
- diff;
- plan file;
- snapshot before change;
- apply;
- post-check;
- rollback;
- canary;
- partial failure.

### Lessons and checkpoints

#### Lesson 01: Build a plan

Checkpoints:

```text
1. Compare rendered config to live config.
2. Produce per-device diff.
3. Build plan object.
4. Save plan.json.
```

Example:

```bash
netops plan --inventory lab.yml --rendered rendered/ --out plan.json
```

#### Lesson 02: Dry-run UX

Checkpoints:

```text
1. Human diff output.
2. Summary: add/change/remove counts.
3. Refuse apply without explicit flag.
4. Add --device filter.
```

#### Lesson 03: Snapshot before apply

Checkpoints:

```text
1. Fetch current configs.
2. Store snapshot with timestamp.
3. Hash snapshot.
4. Attach snapshot metadata to plan.
```

#### Lesson 04: Apply to one device

Checkpoints:

```text
1. Apply config to one lab device.
2. Handle apply failure.
3. Run post-check.
4. Mark device result.
```

#### Lesson 05: Canary deploy

Checkpoints:

```text
1. Apply to canary device.
2. Run State Tester after canary.
3. Continue only if checks pass.
4. Stop and report if canary fails.
```

#### Lesson 06: Rollback

Checkpoints:

```text
1. Restore snapshot.
2. Run post-rollback checks.
3. Report rollback result.
4. Preserve audit trail.
```

### Important rule

Do not position this as “easy safe network deploys”. Teach that safe deployment requires preflight, plan, apply, validate and rollback.

---

## Module 07 — Source of Truth Integration

### Goal

Move from static files to a source-of-truth workflow.

Start with YAML as the simple source of truth. Later add NetBox read-only integration.

### Artifact

```bash
netops sot sync
netops inventory from-netbox
```

### Concepts

- source-of-truth vs live state;
- intended state;
- schema validation;
- mapping external data to internal model;
- read-only integration first;
- NetBox later, not as first dependency.

### Lessons and checkpoints

#### Lesson 01: Static source-of-truth schema

Checkpoints:

```text
1. Define fabric model.
2. Validate roles, ASNs and links.
3. Validate IP uniqueness.
4. Generate inventory from source-of-truth.
```

#### Lesson 02: Intended checks from source-of-truth

Checkpoints:

```text
1. Generate bgp_peer_state checks.
2. Generate interface_state checks.
3. Generate route_exists checks.
4. Run generated checks with State Tester.
```

#### Lesson 03: Intended config from source-of-truth

Checkpoints:

```text
1. Feed source-of-truth into renderer.
2. Render configs for all devices.
3. Detect invalid source data before render.
4. Produce full plan.
```

#### Lesson 04: NetBox read-only integration

Checkpoints:

```text
1. Query NetBox API.
2. Map devices/interfaces/IPs to internal model.
3. Handle missing fields.
4. Generate checks/config from NetBox data.
```

### NetBox note

NetBox should not appear too early. It adds operational complexity. Introduce it after the static source-of-truth path is already clear.

---

## Module 08 — Device Lifecycle: ZTP, Provisioning & Deprovisioning

### Goal

Teach full device lifecycle automation:

```text
planned → ztp_pending → bootstrapped → configured → validated → active → draining → deprovisioned → retired
```

This is the advanced integration module where inventory, renderer, deployer, state tester and source-of-truth come together.

### Naming

Avoid `kickstarter` publicly. Better names:

```text
Device Lifecycle
ZTP & Deprovisioning
Switch Provisioning Pipeline
Fabric Lifecycle Automation
```

Recommended title:

```text
Device Lifecycle: ZTP, Provisioning & Deprovisioning
```

### Commands

```bash
netops device status leaf03
netops ztp render leaf03
netops device provision leaf03 --dry-run
netops device provision leaf03 --apply
netops device drain leaf03 --dry-run
netops device drain leaf03 --apply
netops device deprovision leaf03 --dry-run
netops device deprovision leaf03 --apply --confirm leaf03
```

### Core concepts

- lifecycle state machine;
- bootstrap config vs full intended config;
- preflight checks;
- ZTP simulation in lab;
- post-provision validation;
- drain before removal;
- dependency discovery;
- absence checks;
- source-of-truth state update;
- audit trail;
- rollback/snapshot;
- explicit confirmation for destructive actions.

### Lifecycle states

```go
type DeviceLifecycleState string

const (
    StatePlanned       DeviceLifecycleState = "planned"
    StateZTPPending    DeviceLifecycleState = "ztp_pending"
    StateBootstrapped  DeviceLifecycleState = "bootstrapped"
    StateConfigured    DeviceLifecycleState = "configured"
    StateValidated     DeviceLifecycleState = "validated"
    StateActive        DeviceLifecycleState = "active"
    StateDraining      DeviceLifecycleState = "draining"
    StateDeprovisioned DeviceLifecycleState = "deprovisioned"
    StateRetired       DeviceLifecycleState = "retired"
)
```

### Lab approach

Use ZTP simulation for v1.

Do not start with real vendor-specific DHCP/TFTP/ONIE/ZTP mechanics. They are valuable later but too heavy for first lifecycle module.

In lab:

```text
1. planned leaf03 exists in source-of-truth.
2. netops ztp render leaf03 generates bootstrap config.
3. lst lab adds or enables leaf03 with bootstrap config.
4. netops device provision leaf03 applies intended config.
5. netops check validates BGP/interfaces/routes.
6. lifecycle_state becomes active.
```

### Lessons and checkpoints

#### Lesson 01: Lifecycle state model

Checkpoints:

```text
1. Add lifecycle_state to Device.
2. Implement valid state transitions.
3. Add netops device status.
4. Refuse invalid transitions.
```

#### Lesson 02: Bootstrap config rendering

Checkpoints:

```text
1. Render minimal bootstrap config.
2. Include hostname and management access.
3. Include basic routing or lab reachability.
4. Keep bootstrap separate from full config.
```

Command:

```bash
netops ztp render leaf03
```

#### Lesson 03: Provisioning preflight

Checkpoints:

```text
1. Verify device exists in source-of-truth.
2. Verify hostname/mgmt IP/ASN uniqueness.
3. Verify interface peer references.
4. Refuse provision if active device already exists.
```

Command:

```bash
netops device provision leaf03 --dry-run
```

#### Lesson 04: Provision and validate

Checkpoints:

```text
1. Apply full intended config.
2. Run interface_state checks.
3. Run bgp_peer_state checks.
4. Transition configured → validated → active.
```

Command:

```bash
netops device provision leaf03 --apply
```

#### Lesson 05: Deprovision plan

Checkpoints:

```text
1. Discover BGP peers.
2. Discover routes originated by device.
3. Discover interface references from neighbors.
4. Produce dry-run deprovision plan.
```

Command:

```bash
netops device deprovision leaf03 --dry-run
```

#### Lesson 06: Drain before removal

Checkpoints:

```text
1. Generate drain config.
2. Withdraw originated routes.
3. Disable or shut selected sessions/interfaces.
4. Validate no required routes depend on the device.
```

Command:

```bash
netops device drain leaf03 --apply
```

#### Lesson 07: Remove config and validate absence

Checkpoints:

```text
1. Remove neighbor references from spines.
2. Remove stale interface descriptions/references.
3. Implement bgp_peer_absent.
4. Implement route_absent.
5. Transition deprovisioned → retired.
```

Command:

```bash
netops device deprovision leaf03 --apply --confirm leaf03
```

#### Lesson 08: Lifecycle audit report

Checkpoints:

```text
1. Record pre-change snapshot.
2. Record plan.
3. Record validation results.
4. Write audit report.
```

### Important safety principles

Provisioning and deprovisioning must teach production safety:

```text
1. Dry-run by default.
2. Explicit confirmation for destructive operations.
3. Snapshot before change.
4. Two-phase deprovision: drain, then remove.
5. Validate absence after removal.
6. Update source-of-truth last, after validation.
7. Preserve audit trail.
```

---

## Module 09 — Capstone: NetOps Controller

### Goal

Combine the framework into one cohesive tool.

### Final artifact

```bash
netops
```

Capabilities:

```bash
netops check
netops render
netops plan
netops apply
netops rollback
netops device provision
netops device deprovision
netops audit
```

### Capstone scenario

A lab starts with a two-leaf fabric. Student must:

```text
1. Validate current state.
2. Add leaf03 to source-of-truth.
3. Render bootstrap and intended config.
4. Provision leaf03.
5. Validate leaf03 is active.
6. Introduce a controlled fault.
7. Detect it with State Tester.
8. Fix through deployer.
9. Drain and deprovision leaf03.
10. Produce audit report.
```

### Capstone checkpoints

```text
1. Current-state validation.
2. Add planned device.
3. Render provision plan.
4. Apply and post-check.
5. Detect injected drift.
6. Remediate drift.
7. Drain device.
8. Remove device references.
9. Validate absence.
10. Generate audit report.
```

### Outcome

The student graduates with a serious portfolio artifact:

> A Go-based NetOps framework that validates state, renders config, deploys safely, provisions devices, deprovisions devices, and produces audit reports against a real lab.

---

# 9. Track 3 — Network Observability in Go

## Goal

Turn network state, events and telemetry into metrics and alerts.

This track should reuse concepts from the NetOps Framework:

```text
inventory
state checks
collectors
reports
structured events
```

## Suggested structure

```text
Network Observability in Go
├── 01. Prometheus Exporter
├── 02. Syslog/Event Receiver
├── 03. SNMP Trap Receiver
├── 04. gNMI Subscriber
├── 05. BMP Monitoring
└── 06. Alert Enrichment & Webhooks
```

---

## Module 01 — Prometheus Exporter

### Goal

Expose network state as Prometheus metrics.

### Artifact

```bash
netexporter
```

Example:

```bash
netexporter --inventory lab.yml --listen :9100
```

Metrics:

```text
net_bgp_session_up{device="leaf1",peer="spine1"} 1
net_interface_up{device="leaf1",interface="eth1"} 1
net_route_present{device="leaf1",prefix="10.10.0.0/16"} 1
```

### Lessons and checkpoints

```text
1. Build basic HTTP metrics endpoint.
2. Reuse State Tester collectors.
3. Export interface metrics.
4. Export BGP peer metrics.
5. Add scrape timeout handling.
6. Add labels safely and consistently.
```

---

## Module 02 — Syslog/Event Receiver

### Goal

Receive network events and normalize them into structured data.

### Artifact

```bash
netsyslog
```

Example:

```bash
netsyslog --listen :5140 --inventory inventory.yml
```

### Concepts

- UDP server;
- message parsing;
- structured event model;
- inventory enrichment;
- buffering;
- forwarding.

### Lessons and checkpoints

```text
1. Listen on UDP.
2. Parse basic syslog message.
3. Enrich event with inventory device metadata.
4. Classify event severity.
5. Forward JSON events to file/webhook.
```

---

## Module 03 — SNMP Trap Receiver

### Goal

Decode and enrich SNMP traps.

### Artifact

```bash
trap-receiver
```

### Concepts

- trap receiver;
- OID mapping;
- vendor-specific payloads;
- enrichment;
- alert generation.

### Lessons and checkpoints

```text
1. Receive trap payload.
2. Decode simple trap.
3. Map OIDs to names.
4. Enrich with inventory.
5. Emit structured alert.
```

SNMP is more legacy than gNMI, but still relevant for many network environments.

---

## Module 04 — gNMI Subscriber

### Goal

Introduce modern streaming telemetry.

### Artifact

```bash
gnmi-exporter
```

### Concepts

- gNMI target;
- Subscribe RPC;
- telemetry path;
- updates;
- streaming vs polling;
- mapping telemetry to metrics.

### Lessons and checkpoints

```text
1. Connect to gNMI target.
2. Subscribe to one path.
3. Decode updates.
4. Convert updates to internal event model.
5. Export selected values as Prometheus metrics.
```

This module strengthens modern NetOps positioning.

---

## Module 05 — BMP Monitoring

### Goal

Monitor BGP state and route churn using BGP Monitoring Protocol.

### Artifact

```bash
bmp-monitor
```

### Concepts

- BMP session;
- peer up/down messages;
- route monitoring;
- route churn;
- metrics/events.

### Lessons and checkpoints

```text
1. Accept BMP connection.
2. Parse peer up/down messages.
3. Parse route monitoring message.
4. Track per-peer route counters.
5. Emit churn metrics/events.
```

This should be advanced. Do not put it before Prometheus/syslog/gNMI.

---

## Module 06 — Alert Enrichment & Webhooks

### Goal

Turn raw events into useful alerts.

### Artifact

```bash
netalert
```

### Concepts

- event normalization;
- inventory enrichment;
- deduplication;
- routing;
- webhooks;
- auditability.

### Lessons and checkpoints

```text
1. Define alert event model.
2. Enrich event with device/site/role.
3. Deduplicate repeated alerts.
4. Send webhook payload.
5. Add retry/backoff.
```

---

# 10. Track 4 — Protocol Internals

## Goal

Understand networking protocols by implementing parsers, encoders and state machines in Go.

This is not a “write your own full network stack” track. Avoid root/raw-socket-heavy requirements early.

The focus:

```text
wire formats
binary parsing
encoders
state machines
checksums
fixtures
pcaps
fuzz/property tests
```

## Suggested structure

```text
Protocol Internals
├── 00. Binary Parsing Foundations
├── 01. UDP Datagram Parser/Encoder
├── 02. DNS Parser/Resolver
├── 03. BGP Message Parser
├── 04. BGP FSM as Pure Function
├── 05. OSPF LSDB Parser
├── 06. NetFlow/IPFIX Parser
└── 07. STP BPDU Parser
```

## Why this track exists

This track serves:

- SREs who need deeper protocol intuition;
- network engineers preparing for senior/staff interviews;
- engineers who want to understand what their automation tools observe;
- advanced students who enjoy protocol mechanics.

---

## Module 00 — Binary Parsing Foundations

### Goal

Teach the low-level parsing primitives required by protocol modules.

### Concepts

- byte slices;
- `encoding/binary`;
- Big Endian;
- bit flags;
- length-prefixed fields;
- TLVs;
- checksums;
- parser errors;
- fuzz/property tests.

### Checkpoints

```text
1. Read fixed-width fields.
2. Parse bit flags.
3. Parse length-prefixed payload.
4. Parse TLV sequence.
5. Return offset-aware errors.
```

---

## Module 01 — UDP Datagram Parser/Encoder

### Goal

Parse and encode UDP datagrams from bytes.

### Artifact

```bash
udpdecode
```

Example:

```bash
udpdecode packet.bin
```

### Scope

- parse UDP header;
- source/destination ports;
- length;
- checksum field;
- payload;
- encode datagram;
- optional checksum lesson.

### Checkpoints

```text
1. Parse UDP header.
2. Validate length.
3. Extract payload.
4. Encode UDP header.
5. Optional: checksum.
```

Avoid raw sockets in the first version.

---

## Module 02 — DNS Parser/Resolver

### Goal

Build a small DNS parser and simple resolver.

### Artifact

```bash
dnsmini
```

Example:

```bash
dnsmini query example.com A
```

### Scope

- DNS header;
- question section;
- A/AAAA records;
- name compression;
- UDP query;
- timeout/retry.

### Checkpoints

```text
1. Parse DNS header.
2. Parse question.
3. Parse A record.
4. Implement name compression decoding.
5. Encode simple query.
6. Send UDP query with timeout.
```

DNS compression is an excellent production gotcha.

---

## Module 03 — BGP Message Parser

### Goal

Parse BGP messages without building a full BGP speaker.

### Artifact

```bash
bgpdecode
```

Example:

```bash
bgpdecode update.bin
bgpdecode --hex "..."
```

### Scope

- marker;
- length;
- type;
- OPEN;
- KEEPALIVE;
- UPDATE;
- NOTIFICATION;
- path attributes;
- NLRI basics.

### Checkpoints

```text
1. Parse BGP common header.
2. Parse OPEN message.
3. Parse KEEPALIVE.
4. Parse NOTIFICATION.
5. Parse UPDATE path attributes.
6. Parse simple NLRI.
```

Do not implement a full BGP speaker in this track.

---

## Module 04 — BGP FSM as Pure Function

### Goal

Model BGP state transitions as testable pure functions.

### Artifact

```text
BGP FSM library + small simulator CLI
```

Example:

```bash
bgpfsm --state Idle --event ManualStart
```

### Core idea

```go
func NextState(current State, event Event) State
```

### Checkpoints

```text
1. Define states.
2. Define events.
3. Implement Idle transitions.
4. Implement Connect/Active transitions.
5. Implement Established transitions.
6. Add transition table tests.
```

This is excellent for interviews and protocol reasoning.

---

## Module 05 — OSPF LSDB Parser

### Goal

Parse OSPF LSAs and build a graph from LSDB data.

### Artifact

```bash
ospfgraph
```

### Scope

- Router LSA;
- Network LSA;
- LSDB as graph;
- shortest path calculation.

### Checkpoints

```text
1. Parse LSA header.
2. Parse Router LSA.
3. Parse Network LSA.
4. Build graph.
5. Run shortest path.
6. Output graph summary.
```

This is advanced and should come after DNS/BGP.

---

## Module 06 — NetFlow/IPFIX Parser

### Goal

Parse flow telemetry formats.

### Artifact

```bash
flowdecode
```

### Why add this

NetFlow/IPFIX connects Protocol Internals to Observability and is more relevant to telemetry than STP for many modern environments.

### Scope

- template records;
- data records;
- field mapping;
- exporter/session state;
- flow summaries.

### Checkpoints

```text
1. Parse NetFlow/IPFIX header.
2. Parse template record.
3. Store template state.
4. Parse data record using template.
5. Emit JSON flow event.
```

---

## Module 07 — STP BPDU Parser

### Goal

Understand STP by parsing BPDUs and modeling root bridge selection.

### Artifact

```bash
bpduparse
```

### Scope

- BPDU fields;
- bridge ID;
- root ID;
- path cost;
- root bridge election;
- simple port role decision.

### Checkpoints

```text
1. Parse BPDU header.
2. Parse bridge ID.
3. Compare bridge priorities.
4. Determine root bridge from BPDUs.
5. Simulate simple port role choice.
```

This is useful but lower priority than DNS/BGP/NetFlow.

---

# 11. Recommended Launch Slice

Do not launch all tracks at once.

The full curriculum is the vision. The first product should prove the core learning loop.

## MVP / Alpha slice

```text
BGP / State Tester Alpha
├── Foundation mini-ramp
│   ├── Introduction
│   ├── Inventory & Device Model
│   └── Offline State Assertions
└── NetOps Framework preview
    ├── Lab & Real Device Access
    ├── State Tester: interface_state
    ├── State Tester: route_exists
    ├── State Tester: bgp_peer_state
    └── Reports & CI: JSON + exit codes
```

## First paid artifact

```bash
netops check
```

It should validate:

```text
- BGP peer state
- route existence
- interface state
```

Against:

```text
- fixtures first
- FRR lab soon after
```

## Why this first

This slice proves:

```text
1. web + local lst workflow
2. checkpoint model
3. private checker model
4. domain-specific Go learning
5. state testing value
6. local lab feasibility
```

Without building:

```text
- full web app
- dashboard
- NetBox
- deployer
- lifecycle automation
- observability track
- protocol internals track
```

## Suggested alpha shape

```text
4-6 lessons
10-18 checkpoints
1 starter workspace
1 compiled checker
0-1 lab scenario
5-10 alpha users
```

Goal:

> A stranger can complete the BGP/state tester alpha without your help and say, “I would pay for the full version.”

---

# 12. Build Order

## Phase 0 — Validation

```text
1. Simple landing page.
2. Artifact post: “I built a BGP/state validator in Go.”
3. Waitlist around early access to the guided lab.
4. 5-10 conversations with target engineers.
```

Do not validate with a generic “course coming soon” landing page. Validate with a concrete artifact.

## Phase 1 — Local alpha, no backend

```text
1. Static lesson pages.
2. `lst start/test/next/submit` local only.
3. One starter workspace.
4. One compiled checker.
5. Local progress file.
6. Manual feedback collection.
```

No auth, billing, dashboard or certificates yet.

## Phase 2 — Lab alpha

```text
1. Add `lst lab up`.
2. Add FRR lab.
3. Add lab checks.
4. Add `lst doctor` basic diagnostics.
5. Test on macOS, Linux, WSL2, Apple Silicon.
```

## Phase 3 — Paid v1

```text
1. Add account/auth.
2. Add remote progress sync.
3. Gate premium starter repos/checkers.
4. Add billing.
5. Launch State Tester + Reports & CI as first paid module.
```

## Phase 4 — Expand flagship

```text
1. Config Renderer.
2. Config Deployer.
3. Source-of-Truth Integration.
4. Device Lifecycle.
5. Capstone.
```

## Phase 5 — Electives

```text
1. Network Observability.
2. Protocol Internals.
3. Multi-vendor drivers.
4. EVPN/gNMI advanced tracks.
```

---

# 13. Track Dependencies

## Concept dependency graph

```text
Foundation / Inventory
  → NetOps Framework / Inventory & Drivers

Foundation / Config Fetcher
  → NetOps Framework / Lab Driver

Foundation / Offline State Assertions
  → NetOps Framework / State Tester

NetOps Framework / State Tester
  → Reports & CI
  → Config Deployer post-checks
  → Device Lifecycle validation
  → Observability metrics

Config Renderer
  → Config Deployer
  → Device Lifecycle provisioning/deprovisioning

Source of Truth
  → State Tester generated checks
  → Renderer intended config
  → Device Lifecycle state transitions

Protocol Internals
  → Better parsers and deeper debugging intuition
```

## Most important dependency

State Tester should be built before Config Deployer.

Reason:

> You should know how to verify the network before teaching students how to mutate it.

---

# 14. Example Checkpoint Sequence: State Tester

This example shows how checkpoints can build one growing application without exposing private tests.

## Workspace

```text
~/linkstate/netops-framework/state-tester/
├── go.mod
├── cmd/
│   └── netops/
│       └── main.go
├── internal/
│   ├── model/
│   │   └── state.go
│   ├── parser/
│   │   ├── bgp.go
│   │   ├── routes.go
│   │   └── interfaces.go
│   ├── checks/
│   │   ├── bgp_peer_state.go
│   │   ├── route_exists.go
│   │   └── interface_state.go
│   └── report/
│       └── json.go
├── testdata/
├── CURRENT.md
└── .linkstate/
    ├── manifest.yml
    └── progress.json
```

## Lesson: Model check results

### Checkpoint 1: Define `CheckResult`

Task:

```text
Create CheckResult with Name, Type, Device, Status, Expected, Actual, Reason.
```

### Checkpoint 2: Add status helpers

Task:

```text
Add Passed(), Failed(), Unknown() helpers.
```

### Checkpoint 3: Stable result ordering

Task:

```text
Sort by device, type, name.
```

## Lesson: Parse FRR BGP summary

### Checkpoint 1: Parse Established peer

Input:

```text
10.0.0.1 4 65001 1234 1201 42
```

Expected:

```text
State = Established
Prefixes = 42
```

### Checkpoint 2: Parse Active/Idle peers

Input:

```text
10.0.0.2 4 65002 900 887 Active
```

Expected:

```text
State = Active
Prefixes = 0
```

### Checkpoint 3: Parse full output

Task:

```text
Skip headers, parse rows, return []BGPSession.
```

## Lesson: Implement `bgp_peer_state`

### Checkpoint 1: Find peer by address

### Checkpoint 2: Compare expected state

### Checkpoint 3: Return discrepancy reason

### Checkpoint 4: Add missing peer behavior

## Lesson: Reports & CI

### Checkpoint 1: Human report

### Checkpoint 2: JSON output

### Checkpoint 3: Exit code 1 for failed checks

### Checkpoint 4: Exit code 2 for tool/runtime errors

---

# 15. Naming Recommendations

## Track names

Recommended:

```text
Foundation: Go for NetOps
Build Your NetOps Framework
Network Observability in Go
Protocol Internals
```

## Tool names

Use short, practical names:

```text
netipcalc
config-fetcher
mini-netcheck
netops
netexporter
netsyslog
trap-receiver
gnmi-exporter
bmp-monitor
udpdecode
dnsmini
bgpdecode
ospfgraph
flowdecode
bpduparse
```

## Avoid

Avoid public course names that sound like toys:

```text
kickstarter
coding puzzles
toy lab
```

For the provisioning/deprovisioning module, use:

```text
Device Lifecycle
ZTP & Deprovisioning
Switch Provisioning Pipeline
```

---

# 16. What to Defer

Do not build these before the first successful alpha:

```text
- browser IDE
- VS Code extension
- JetBrains plugin
- certificates
- team dashboard
- server-side grading
- cloud labs
- NetBox-heavy flow
- full config deployer
- device lifecycle module
- observability track
- protocol internals track
- multi-vendor support
- EVPN automation
```

These are valid future features, but they should be earned by usage.

---

# 17. Summary Decisions

| Area | Decision |
|---|---|
| Pedagogy | Small checkpoints inside larger artifact lessons |
| Execution | Local editor + `lst`, no browser IDE v1 |
| Public content | Reading pages public as marketing |
| Product value | Starter repos, checkers, labs, submit, progress, solutions gated |
| Tests | Compiled private checkers, not full public unit tests |
| First paid artifact | State Tester / `netops check` |
| First mutation module | Config Deployer after State Tester + Renderer |
| Lifecycle module | Device Lifecycle: ZTP, Provisioning & Deprovisioning |
| Foundation | No Docker, no lab, local fixtures/mocks only |
| Flagship lab | FRR via netlab/containerlab |
| Observability priority | Prometheus exporter → syslog → traps → gNMI → BMP |
| Protocol priority | Binary parsing → UDP → DNS → BGP → OSPF/NetFlow/STP |
| Launch strategy | One guided artifact first, not full platform |

---

# 18. Final Curriculum Snapshot

```text
LinkState

Track 1: Foundation — Go for NetOps
  00. Introduction to LinkState Workflow
  01. IP Toolkit
  02. Inventory & Device Model
  03. Config Fetcher with MockDriver
  04. Offline State Assertions
  05. Capstone: Mini NetCheck

Track 2: Build Your NetOps Framework
  00. Framework Introduction
  01. Lab & Real Device Access
  02. Inventory, Drivers & Collectors
  03. State Tester
       - interface_state
       - route_exists / route_absent
       - bgp_peer_state / bgp_peer_absent
       - later: prefix_count, lldp_neighbor, bfd_session_state
  04. Reports & CI
       - human output
       - JSON
       - exit codes
       - audit report
  05. Config Renderer
       - intended config
       - platform templates
       - bootstrap vs full config
  06. Config Deployer
       - plan
       - dry-run
       - snapshot
       - apply
       - canary
       - rollback
  07. Source of Truth Integration
       - static YAML first
       - generated checks
       - generated configs
       - NetBox read-only later
  08. Device Lifecycle: ZTP, Provisioning & Deprovisioning
       - lifecycle state machine
       - bootstrap config
       - ZTP simulation
       - provision
       - validate
       - drain
       - deprovision
       - validate absence
       - retire
       - audit
  09. Capstone: NetOps Controller

Track 3: Network Observability in Go
  01. Prometheus Exporter
  02. Syslog/Event Receiver
  03. SNMP Trap Receiver
  04. gNMI Subscriber
  05. BMP Monitoring
  06. Alert Enrichment & Webhooks

Track 4: Protocol Internals
  00. Binary Parsing Foundations
  01. UDP Datagram Parser/Encoder
  02. DNS Parser/Resolver
  03. BGP Message Parser
  04. BGP FSM as Pure Function
  05. OSPF LSDB Parser
  06. NetFlow/IPFIX Parser
  07. STP BPDU Parser
```

---

# 19. The Most Important Practical Advice

Do not start by building all of this.

Start with one module that proves the product:

```text
State Tester Alpha
  - one workspace
  - public lessons
  - local `lst`
  - compiled checker
  - 10-18 checkpoints
  - BGP/interface/route checks
  - JSON output + exit codes
  - optional FRR lab at the end
```

The long-term product is a platform.

The first release should be:

> **one guided NetOps artifact that feels magical.**

If that works, the full curriculum can grow naturally.
