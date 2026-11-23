# CLAUDE.md

Working instructions for this repository. Read this file first, then the documents it links to. It does not repeat them.

## What this is

A local lab for people learning image formats: open a supported PNG or grayscale JPEG, read its chunks or markers, inspect the decoded samples, encode a raster and measure the difference. Go standard library, own framing, filters, transforms and entropy coding, `127.0.0.1` only, original fixtures.

**Nothing is implemented yet.** The repository holds the plan, the delivery tooling and the backlog in epic #1. Commands marked planned below do not work yet. Do not describe a planned thing as working, and do not claim a passing baseline: there is none.

## Sources of truth

| Topic | Document |
|---|---|
| Product scope, design, fixtures with expected values | [docs/PROJECT-PLAN.md](docs/PROJECT-PLAN.md) |
| Planned module map, freeze points, planned versus implemented | [docs/architecture.md](docs/architecture.md) |
| Backlog | Epic issue #1 and its child tickets on GitHub. The epic's table is authoritative for order, dependencies and estimates; tickets state their own owner-decision defaults. |
| Loop setup and its limits | [LOOP.md](LOOP.md), [loop-constraints.md](loop-constraints.md) |
| Project subagents, origin and local changes | [.claude/agency-agents/README.md](.claude/agency-agents/README.md) |
| Ticket workflow and the step at which each role runs | [.claude/skills/implement-ticket/SKILL.md](.claude/skills/implement-ticket/SKILL.md), section "Agent roles" below |
| How a ticket goes from the queue to a ready pull request, who may do what | "Delivery workflow" below |
| UI and copy rules | [docs/ui-quality.md](docs/ui-quality.md) |
| Decisions the owner made, with dates | [docs/owner-decisions.md](docs/owner-decisions.md) |

The plan wins for product design; the epic and tickets win for order and scope. A ticket that changes a contract fixes the plan in the same pull request.

## Stack and constraints

- Go module, **no external dependencies**: standard library only, in the application and in the tests.
- One executable: command line and loopback server share the codec packages; browser assets are embedded.
- `image/png` and `image/jpeg` are verifiers in tests and fixture preparation, never in production code.
- Offline: same-origin assets only. Original synthetic images only.
- PNG: 8-bit gray, RGB, RGBA, not interlaced. JPEG: 8-bit baseline, one grayscale component, one scan, no restart markers.
- Indexed colour, Adam7, colour JPEG, restart intervals, a quality control, colour management and conformance claims are later work.
- A runtime version change, a dependency, an external request or a non-loopback listener needs the owner's approval first, not a pull request.

## Commands

Go is not installed on the owner's PC and nothing is installed globally. **Go checks run in CI only.** Run `go version` first: if it is missing, do not install it, and do not claim you ran any command below.

| Purpose | Command | Status |
|---|---|---|
| Build | `go build -o codeclab ./cmd/codeclab` | planned: real with the scaffold |
| Format check (empty output is good) | `gofmt -l .` | planned: real with the scaffold |
| Static check | `go vet ./...` | planned: real with the scaffold |
| Tests | `go test ./...` | planned: real with the scaffold |
| Run the app | `./codeclab serve` | planned: #2 |
| Inspect a file | `./codeclab inspect <file>` | planned: #3 (PNG), #5 (JPEG) |
| Encode a raster | `./codeclab encode --format png <raster.json>` | planned: #4 (PNG), #6 (JPEG) |
| Release checks | `./codeclab verify` | planned: #7 |

## Architecture boundaries

The planned module map is in [docs/architecture.md](docs/architecture.md). Invariants every change must keep (the plan is the source):

- One raster contract: positive size, layout gray, RGB or straight RGBA, tight 8-bit row-major samples. The PNG and JPEG packages know the raster and the report, not HTTP or the page.
- Own code does framing, filters, transforms and entropy coding. Production may use `compress/zlib` and the checksum packages; `image/png` and `image/jpeg` appear only in tests and fixture preparation.
- Limits before allocation: input 8 MiB, PNG 1024×1024, JPEG 512×512, output 12 MiB, 256 report entries with an omitted count, two codec operations at once. Sizes use checked multiplication.
- Four stable error categories: invalid data, unsupported profile, resource limit, busy. A failed operation returns no partial raster and keeps the last valid preview.
- PNG: signature, first `IHDR`, contiguous `IDAT`, final `IEND`, length and CRC checked before use, types 0, 2 and 6 at 8 bit, no interlace. `PLTE`, `tRNS` and unknown critical chunks are unsupported.
- The IDAT run is one zlib stream, read to exactly the expected length plus one detection byte, with a clean end and no trailing member.
- PNG filters 0 to 4 in modulo-256 arithmetic with wider intermediates; encoding picks the filter by the sum of absolute signed residuals with a fixed tie rule. Samples are lossless, also colour under zero alpha.
- JPEG: one 8-bit SOF0 frame, one 1×1 component, one full scan, no restart markers; nonzero 8-bit quantisers; canonical Huffman tables with checked counts; stuffed bytes and markers kept apart; padded blocks cropped to the frame.
- JPEG output uses the fixed documented quantiser and own valid tables. Loss is reported as maximum and mean absolute sample difference, only for the measured images.
- Every fixture names its independent reference. Decoded samples are compared, never bytes of two encoders. The one-block JPEG fixture is built without the production encoder.
- The preview draws raw RGBA bytes with `ImageData`. Neither the input nor the output is ever loaded as an image element, and embedded texts or metadata values are never shown.
- The server binds `127.0.0.1`, checks host and origin on writes, caps bodies and builds no path from an HTTP value. Exports are a deliberate save without source names or metadata.
- The sample inspector and every action work by keyboard. Values, errors and differences are text; the checkerboard never replaces the numeric alpha.

## Git workflow

- `main` is the default branch. Work on one branch per ticket: `<issue-number>-<short-slug>`. Never commit to `main`.
- Small commits with imperative subjects and a type prefix (`feat:`, `fix:`, `docs:`, `test:`, `chore:`).
- Order for every ticket: local evidence, Code Reviewer, Reality Checker, then a draft pull request that says `Part of epic #1` and `Closes #<ticket>`, then CI on that draft, then `gh pr ready` only when CI is green.
- Do not push or open a pull request unless the owner asked in the current session; starting `/next-ticket` is that request for one ticket. Only the owner merges and releases. Do not push to `main`, force-push, or merge your own pull request.
- Never commit exported images (`out/`), `.env*`, or build output.
- Do not change `.claude/agents/`, `LOOP.md`, `loop-constraints.md` or `loop-budget.md` as part of a feature ticket.
- Open owner decisions on a ticket are settled with the ticket's stated default. Where a ticket states none, take the conservative option and say so in the pull request.

## Evidence before a ticket is called done

Paste real output into the pull request; never summarise a run you did not make.

1. For each acceptance criterion of the issue: the name of the test that proves it, or the exact step and what was observed.
2. Go results come from CI, and green CI on the current diff is a condition for `gh pr ready`: the link to the CI run on the draft pull request with the Go version and the results of `gofmt -l .`, `go vet ./...` and `go test ./...`. If Go is not available locally, say so.
3. For a codec change: the fixture, its independent reference, and the expected and observed samples or the measured difference.
4. For a limit or error change: the input, the category and stage, and that a following valid file decodes as in a fresh session.
5. For a UI change: the browser and version, the steps and what was observed, and the manual steps left as owner checks.
6. `git status` shows no stray files, and no dependency, external URL or non-loopback listener was added.
7. Anything not verified, said plainly.

## Agent roles

Project subagents are in `.claude/agents/`: eight Agency Agents roles and `loop-verifier`, which belongs to the Loop setup. Each Agency Agents role file starts with a section "In this repository", written for this project. That section and this file override the generic text below it. Sources, licence, local changes and the name collision with a user-level install are in [.claude/agency-agents/README.md](.claude/agency-agents/README.md).

When asked to implement, continue or finish a GitHub issue, follow the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Claude Code loads it on its own for ticket work; `/implement-ticket <issue number>` starts it by hand. Delegate to the roles without waiting to be asked: each role's `description` ends with a "Use proactively ..." sentence, which is what makes Claude Code delegate. The playbook gives the exact trigger of each role.

| When | Roles | Access |
|---|---|---|
| Before coding | Software Architect | `Read, Grep, Glob` |
| While coding | Backend Architect, Frontend Developer | can edit |
| After coding, when the trigger matches | Application Security Engineer, Accessibility Auditor | `Read, Grep, Glob` |
| After coding, when the trigger matches | API Tester | can edit |
| Gate on every ticket | Code Reviewer on the full diff, until no material finding is left | `Read, Grep, Glob` |
| Acceptance gate on every ticket, after Code Reviewer | Reality Checker on each acceptance criterion and the recorded evidence, and that CI is configured to run the checks (no CI result is needed yet) | `Read, Grep, Glob` |

- The two gates run on every ticket, whatever its size. A ticket is not called done before both have passed.
- Roles limited to `Read, Grep, Glob` cannot run commands. Give them the diff, the acceptance criteria and whatever command output exists.
- The session or role that wrote a change does not count as its reviewer.
- If a role is missing or cannot be delegated to, the session does the before-coding, while-coding and specialist steps itself and says so. It cannot pass the two gates on its own change: do not push or open the pull request. Stop, report that the gates could not run, and leave the run for the owner. This overrides the playbook's sentence "do that step yourself" for the two gates.
- In the pull request description, list under "Agent roles run" each role that ran, what it was asked and its verdict, and any step done without a role.

## Delivery workflow

How a ticket reaches the owner. Each piece is described in its own file; this part only links them.

- **Start:** the owner runs `/next-ticket` ([`.claude/skills/next-ticket/SKILL.md`](.claude/skills/next-ticket/SKILL.md)); it picks one leaf ticket (never a `tracking` parent or a ticket with open sub-issues; the last sub-issue's pull request also closes the parent), creates the branch and worktree, and runs the playbook [`.claude/skills/implement-ticket/SKILL.md`](.claude/skills/implement-ticket/SKILL.md). Nothing is scheduled.
- **Gate hook:** [`.claude/hooks/role_gate.py`](.claude/hooks/role_gate.py) blocks `gh pr create` and `gh pr ready` until Code Reviewer and Reality Checker ran after the last code edit. `ROLE_GATE=off` in the environment Claude Code is started from switches it off for that session; only the owner does that.
- **Permissions:** [`.claude/settings.json`](.claude/settings.json). Agent configuration edits ask in the main checkout and in ticket worktrees; CI, `CLAUDE.md`, `AGENTS.md` and Loop-file edits ask only in the main checkout; in a ticket worktree they are limited by the playbook's scope rule and reviewed by the owner at merge (shell-made edits are not covered, so a ticket never changes agent configuration); force pushes, pushes to `main`, `gh pr merge`, releases and reads of `.env` are denied.
- **CI:** [`.github/workflows/ci.yml`](.github/workflows/ci.yml). Every job must be green on the current diff before a pull request is marked ready. The agent does not poll: it ends its turn with "Draft pull request #N is waiting for CI" and reads `gh pr checks` once when woken. Reviewers read the diff and evidence from `.claude/state/review/`. An acceptance step that needs a person is listed under "Owner checks" in the pull request.
- **Pull request:** [`.github/pull_request_template.md`](.github/pull_request_template.md). Tickets use [`.github/ISSUE_TEMPLATE/implementation-ticket.md`](.github/ISSUE_TEMPLATE/implementation-ticket.md).
- **Labels:** `status:ready`, `status:blocked`, `status:in-progress`, `needs-decision`, `tracking` (parent container, never worked; its sub-issues are the work items).
- **Changelog:** every ticket pull request adds one line under Unreleased in [CHANGELOG.md](CHANGELOG.md).
- **UI and copy:** [docs/ui-quality.md](docs/ui-quality.md); the `ui-ux-pro-max` and `unslop` skills are in `.claude/skills/`.
- **Decisions:** [docs/owner-decisions.md](docs/owner-decisions.md). Security reports: [SECURITY.md](SECURITY.md).
- **Merge and release:** only the owner merges and releases ([docs/RELEASING.md](docs/RELEASING.md)). An agent stops at a ready pull request.

## Loop

Issue triage runs at L1: report-only, started by hand, no scheduler, no GitHub writes. See [LOOP.md](LOOP.md).
