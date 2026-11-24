# PNG and JPEG Codec Workbench

> This document specifies phase 1 of the program in [#34](https://github.com/jortega0033/png-jpeg-codec-workbench/issues/34). Hour figures in it are early estimates for that phase; its open tickets add up to 56 to 86 hours (2026-11-24). Work that it defers is planned in the later phases only where the program names it.

## Product requirements

Give people learning image formats a local place to see what a codec reads and writes. They can open an original sample or a supported file, inspect the chunk or marker sequence, and preview the pixels reconstructed by the project's decoder. They can also encode a bundled raster, inspect the output and measure its differences from the source.

**First release.** Support noninterlaced 8-bit grayscale, RGB and RGBA PNG plus single-component 8-bit baseline sequential JPEG with one scan and no restart markers. PNG is lossless at the sample level. JPEG exposes lossy differences and uses the project's fixed quantizer, with no quality slider in this release. A PNG preview may show transparent pixels over a checkerboard while the pixel inspector retains their original straight-alpha component values. JPEG displays grayscale replicated into all three color channels. The benchmark uses synthetic images the repository can distribute.

**User journeys.** Learn how PNG filters reconstruct a row and how JPEG quantization changes a block. Open a supported file or select a generated fixture, then compare its dimensions and decoded samples with a reference. When a format or size is outside the contract, get the reason without losing the last valid preview. Export an encoded sample to the local workspace and reopen it with the workbench and an independent reader.

**Success criteria.** A new contributor can build and open the same offline Chrome page on Linux, reproduce independently checked PNG and JPEG sample outputs and inspect the decoder's reported format structure. Reference checks use exact sample comparisons for PNG and only the documented tolerances on specific JPEG fixtures. A mistaken input cannot replace the currently displayed result or allocate outside the documented caps.

**Later milestones.** Add PNG indexed color and Adam7 interlace after the first release. Add color JPEG, restart intervals and configurable quantization in separate slices. Expand the shareable comparison corpus, quality and speed report, and reference coverage after each profile expands. Full PNG or JPEG conformance is a separate project goal with its own matrix.

## Technical design

### Runtime and boundaries

- Use Go with standard-library packages, a local HTTP server bound to `127.0.0.1`, and HTML, CSS and JavaScript shipped in the binary. Record the verified Go and Chrome versions in release evidence. Use the installed toolchain offline with network-backed dependency and toolchain discovery disabled. No GPU or external service is needed.
- Keep PNG and JPEG codec packages behind an explicit raster contract. One raster has positive dimensions, a layout of gray, RGB or straight RGBA, and tightly packed 8-bit samples in row-major order. Format parsing returns normalized samples and an ordered, bounded inspection report. The inspector records offsets and structural fields, with no embedded file texts or image metadata values shown.
- The preview consumes raw RGBA created by expanding this raster. Use Canvas `ImageData`, with a separate raw pixel inspector reading from the expanded bytes. Load neither the original input nor the encoded output as an image element for this preview. Native browser decoding may be a separately labeled optional comparison after the baseline.
- Provide `GET /`, `GET /api/fixtures`, `POST /api/decode`, and `POST /api/encode?format=png|jpeg`. Encoding takes a bounded raw raster envelope from the built-in fixture or current preview. The JPEG endpoint requires a gray raster and gives an explicit unsupported-layout error for RGB/RGBA. The CLI supports encoding, decoding to raw samples and writing an inspection report, and shares the production codec packages. Exports use a deliberate local save action.
- Cap incoming compressed files at 8 MiB, PNG dimensions at 1024×1024, and JPEG dimensions at 512×512. Validate positive width and height, channel count and bounded multiplication before reserving output. Cap report entries at 256, report omitted entries by count, and never retain segment payloads in the report. Limit request bodies, cap concurrent codec operations to two, and set a 12 MiB limit on the encoded file returned to the user. A third work request can return a clear busy response. Use stable error categories for invalid data, unsupported profile, resource limit and busy. Failed operations leave the preceding successful UI result intact.

### PNG profile

- Verify the signature, require one first `IHDR`, a contiguous `IDAT` run and one final `IEND`. Check the declared length and CRC-32 of every chunk before using it. Support 8-bit color types 0, 2 and 6 with compression, filter and interlace methods 0. Unknown critical chunks are unsupported. Reject profile-affecting `tRNS` and `PLTE` in this release. Skip bounded ancillary metadata, and document that the workbench displays raw samples without profile transformations or metadata preservation. [W3C Portable Network Graphics Specification](https://www.w3.org/TR/png/)
- Decode the IDAT payload as one zlib stream using `compress/zlib`. Bound inflated data at `height × (1 + width × componentCount)` plus one detection byte and require the exact expected length, a clean zlib end, and no trailing compressed member. Advance across IDAT chunk boundaries without resetting the decompressor. Reconstruct filter types 0 through 4 with modulo-256 sample arithmetic, zero neighbors outside the row and the image, and bytes per pixel matching the supported profile. Use wider intermediates for Average and Paeth.
- Encode the same three layouts. Frame independent chunks, write correct CRCs, and use deterministic filter selection with a fixed tie rule. For the first pass choose one filter per row by the sum of absolute signed residuals across all five candidates. Keep output chunking stable without requiring byte equality across Go toolchain versions. The project's code owns framing and filtering. Go provides generic zlib compression and checksum primitives. `image/png` belongs only to comparisons and fixture preparation. [Go zlib package](https://pkg.go.dev/compress/zlib) · [Go PNG reference reader and writer](https://pkg.go.dev/image/png)

### JPEG profile

- Accept one positive-sized 8-bit SOF0 frame, one grayscale component sampled 1×1, 8-bit nonzero quantization entries and one full sequential SOS. Permit bounded application/comment segments and skip their payloads. Reject progressive frames, restart intervals, additional components or scans, and other coding processes with an unsupported-profile error. A valid-looking but incomplete structure yields invalid data. Do not treat TIFF, JPEG 2000 or alternate extensions as baseline JPEG.
- Use length-checked SOI, DQT, DHT, SOF0, SOS and EOI parsing. Build canonical Huffman codes from declared counts and symbols, enforce baseline symbol bounds and reject impossible code tables. Treat stuffed FF bytes and true markers distinctly. Decode the expected row-major 8×8 blocks using DC differences and AC runs, enforcing the block extent and end of scan before EOI. The pixel pipeline applies inverse zigzag, dequantization, a separable inverse DCT, level shift, rounding and clamp. Crop right and bottom padded blocks to the exact frame size. [JPEG Committee overview of JPEG 1](https://jpeg.org/jpeg/) · [ITU-T T.81 baseline coding reference](https://www.w3.org/Graphics/JPEG/itu-t81.pdf)
- Encode by edge-replicating into partial blocks, subtracting 128, applying separable forward DCT, and quantizing before zigzag and entropy coding. Specify the 64 positive entries of an original fixed quantizer. Use the supported baseline coefficient categories and deterministic, valid Huffman tables owned by this code. Emit the corresponding DQT and DHT and apply scan byte stuffing and end padding. The workbench reports the quantizer and measured pixel differences. `image/jpeg` is allowed as an independent verifier only. [Go JPEG reference reader and writer](https://pkg.go.dev/image/jpeg)

### Verification fixtures

| Fixture | Independent expected result |
| --- | --- |
| 2×2 RGBA PNG containing opaque and transparent straight colors | The exact two source rows match the decoder and `image/png` after conversion to `image.NRGBA`. At least one fully transparent pixel has nonzero RGB and the raw pixel inspector retains those original components. |
| One 2-pixel RGBA row, `10,20,30,255,11,22,33,250` | With Sub filter and a four-byte stride, the residual bytes are `10,20,30,255,1,2,3,251`. Reconstruct each original byte. |
| Hand-authored grayscale JPEG with one 8×8 block and all-zero quantized coefficients | A one-symbol DC-zero table and one-symbol AC-EOB table each use code `0`. The entropy payload is `0x3f` with the two code bits and six one bits of padding. Every output sample is 128. Build the fixture independently of the production encoder. |
| Odd-sized 13×9 grayscale JPEG fixture with block-aligned constant regions | Preview has 117 grayscale pixels and no padding columns or rows. Compare field dimensions and reconstructed samples with Go's independent decoder, using at most two levels difference at each pixel on this fixture. Use separate AC reference samples with nonzero coefficients to exercise zigzag, signed values and nontrivial entropy paths. |
| PNG with changed CRC, incomplete zlib data, oversized dimensions or unknown critical type, and JPEG with truncated segment or scan | Each input produces a stable error category without a partial preview. A following valid file gives the same output as in a fresh session. |

Record which independent reference generates or reads each fixture. Compare decoded sample arrays, not bytes produced by different encoders. Restrict any padded-image or lossy result claims to the committed fixture set and measured tolerance. For PNG use the W3C filter rules and independent reference output. For JPEG confirm both external decode of workbench output and workbench decode of the hand-authored fixture, so matching defects in its writer and reader cannot cancel silently.

## Epic · Deliver a local PNG and grayscale JPEG codec lab

Ship the workspace where learners can inspect either format and trace pixels from encoded input to independently checked raster output. The planned effort is **23–36 engineering hours** including verification. IC-01 defines the contracts, then PNG and JPEG work can proceed in parallel. IC-03 extends the PNG reader into a lossless export path while IC-05 adds JPEG output after the decoder milestone. IC-06 ties both paths into one repeatable offline run. All verification results are future deliverables.

### IC-01 · Build the raster contract and first fixture preview

**Estimated effort.** 3–4 hours. **Dependencies.** None.

- Start the loopback server and show an original 2×2 raw RGBA fixture when the page opens. Match the four pixel values in a keyboard-accessible sample inspector. Build the reusable raster and error contracts.
- Verify inputs for zero size, bad channel count, mismatched sample length and size-cap violation before allocation. Show structured errors without replacing a valid preview.
- Run offline from a clean checkout with the installed toolchain. A trace shows the preview built from sample bytes delivered by the app. It shows no browser decode of the input file and no external request.

### IC-02 · Decode the noninterlaced PNG profile

**Estimated effort.** 4–6 hours. **Dependencies.** IC-01.

- Own chunk framing, checksum verification, bounded zlib access and all five inverse row filters. Render raw decoded pixels and an ordered inspection report for a supported input.
- Use independently generated PNG fixtures for each supported layout and all five filters. Assert the exact RGBA fixture and the Sub vector. Compare all admitted fixtures with Go's independent decoder.
- Include passing adjacent IDAT chunks and clean failures for ordering, unsupported layout, invalid CRC, truncated stream and excess inflated data. Report which stage rejected each input.

### IC-03 · Export supported rasters as lossless PNG

**Estimated effort.** 3–5 hours. **Dependencies.** IC-01 and IC-02.

- Write IHDR, deterministically filtered and compressed IDAT data and IEND. Export and reopen all three raster layouts. Retain nonzero RGB at zero alpha in the original RGBA test fixture.
- Compare all decoded source and output samples exactly with an independent PNG reader. Verify output with at least one PNG containing multiple IDAT chunks. Normalize test images using a straight-alpha reference model.
- Show source and reopened dimensions, sample layout and byte length. No debug fixture path or source metadata is written to an export.

### IC-04 · Decode one-scan baseline grayscale JPEG

**Estimated effort.** 5–8 hours. **Dependencies.** IC-01.

- Build parser, entropy decoder, inverse transform, bounded block assembly and ordered marker report. Show the original size after cropping partial edge blocks.
- Decode the hand-authored one-block fixture to 64 exact samples of 128. Check DC prediction across blocks, nonzero signed AC values, byte stuffing and the odd-sized reference comparison. Expected JPEG reference differences are restricted to the stated fixtures.
- Return distinct deterministic errors for unsupported frames, incomplete streams, bad table declarations and blocks that do not fit the profile before presenting an image.

### IC-05 · Write baseline grayscale JPEG and quantify loss

**Estimated effort.** 5–8 hours. **Dependencies.** IC-04.

- Add forward transforms, fixed quantization, legal baseline tables, stuffing and marker writes. Export grayscale sample sizes 8×8 and 13×9 and reopen with an independent JPEG decoder.
- Check for exact frame dimensions after reference decode. Confirm a flat 128 source remains flat and include an original patterned source where the preview shows a measured difference after encoding. Document the comparison metric as maximum and mean absolute sample difference.
- Reject unsupported raster layouts without discarding the original preview. Show the re-decoded JPEG and measured difference from the original as separately identified views.

### IC-06 · Verify both codec paths in the offline workbench

**Estimated effort.** 3–5 hours. **Dependencies.** IC-02, IC-03 and IC-05.

- From a clean Linux checkout record the Go and Chrome versions, run the self-contained unit and reference checks, and verify the import, inspect, export and reopen workflows for both formats in Chrome.
- Use the original fixtures and bounded regression cases for size, partial inputs and resource limits. Confirm a failed request leaves the prior valid preview intact. Check readout of a transparent raw PNG pixel and the padded-edge JPEG fixture.
- Attach local output captures and the exact observed result of each check to the project verification report. Record the supported profiles and their limitations and keep codec code, test-reference use and the UI preview path distinguishable.
