# PNG and JPEG Codec Workbench

An offline workbench for writing, inspecting and verifying original PNG and baseline JPEG codecs.

Only a runnable scaffold exists; the codecs and the workbench are planned. The scope, architecture and acceptance criteria are in [the project plan](docs/PROJECT-PLAN.md). The plan includes a reproducible verification strategy for the future build.

- [Product requirements](docs/PROJECT-PLAN.md#product-requirements)
- [Technical design](docs/PROJECT-PLAN.md#technical-design)
- [Implementation backlog and epic](https://github.com/jortega0033/png-jpeg-codec-workbench/issues)

The first plan of phase 1 said 23–36 engineering hours. Its open tickets, with tests and review, add up to 56 to 86 hours; the whole program is larger and is described below. All figures are estimates.

## Program goal

Build a complete family of image codecs from the specifications: PNG, JPEG, WebP, AVIF and JPEG XL, encoders and decoders, with no image library, a workbench that shows every coding decision, and an encoder that competes with the best ones.

The program is planned ([#34](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/34)). Phase 1 is the current epic [#1](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/1) and the only phase with tickets.

**Score.** Share of the public conformance files per format that decode exactly, and compressed size against reference encoders at equal quality.

| Phase | Content | Hours (rough) | Issue |
|---|---|---|---|
| 1 | The current epic and its tickets | 56 to 86 | [#1](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/1) |
| 2 | PNG and JPEG complete | 100 to 170 | [#35](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/35) |
| 3 | Own DEFLATE and better encoders | 110 to 180 | [#36](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/36) |
| 4 | WebP and GIF | 100 to 170 | [#37](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/37) |
| 5 | AVIF and JPEG XL | 250 to 450 | [#38](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/38) |
| 6 | Speed and tools | 110 to 190 | [#39](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/39) |

About 730 to 1250 engineer hours in total: rough estimates, not measurements. The numbers in the score and the finish lines are targets until a run shows them. Only phase 1 has tickets to work; the sub-issues of a later phase are written when it starts. No single session is expected to reach the end.

## License

[MIT](LICENSE) © 2026 Jake Ortega.
