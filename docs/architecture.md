# Implementer map

Status: **only the scaffold exists, no product behaviour** (packages are empty and verified in CI only). This page maps the boundaries that [PROJECT-PLAN.md](PROJECT-PLAN.md) requires onto a module layout so that parallel tickets do not collide. It adds no requirement; the plan wins where they differ. The names are the orchestrator's proposal, fixed by the scaffold; a ticket that changes one corrects this page in the same pull request.

## Planned module map

| Boundary | Module | Owns | Parent ticket |
|---|---|---|---|
| Entry point | `cmd/codeclab` | Subcommands `serve`, `decode`, `encode`, `inspect` | #2 (`serve`), #3 to #6 (codec commands), #7 |
| Raster and report | `internal/raster` | Raster contract, RGBA expansion, limits, error categories, inspection report | #2 |
| PNG | `internal/pngcodec` | Chunks, inflate bound, filters; encoder | #3 (decode), #4 (encode) |
| JPEG | `internal/jpegcodec` | Markers, Huffman, transforms, blocks; encoder and difference | #5 (decode), #6 (encode) |
| Fixtures | `internal/fixtures` | Original rasters and hand-authored streams | #2, then each codec ticket |
| HTTP | `internal/httpapi` | Routes, guards, limits, busy rule, static assets | #2, then each codec ticket |
| Browser assets | `web/ (embedded)` | Preview, sample inspector, structure report, export | #2, then each codec ticket |
| Verification | `internal/verify` | Both paths, regressions, report | #7 |

Dependency direction: `httpapi` and `cmd/codeclab` depend on `pngcodec`, `jpegcodec`, `fixtures` and `raster`; `pngcodec`, `jpegcodec` and `fixtures` depend on `raster` only. `raster` imports nothing from the module.

## Freeze points

Parallel tickets rely on these shapes. A ticket that changes a frozen shape updates this page and every dependent ticket in the same pull request.

- **Raster, error and inspection report contracts**: ticket #2; the PNG and JPEG lines build on them in parallel.
- **HTTP shapes** (`GET /api/fixtures`, `POST /api/decode`, `POST /api/encode?format=png|jpeg`, raster envelope, errors): #2.
- **Fixed JPEG quantiser and Huffman tables**: #6, written in the docs.

## What is verified where

Go is verified in CI only (see [CLAUDE.md](../CLAUDE.md)), on Linux, the supported platform. Manual browser steps are owner checks in Chrome.
