# FSRS replay harness

This directory holds an offline, synthetic replay harness. It compares three
scheduling models over the same review sequence:

- `sm2`: the scheduler in `internal/scheduler/sm2.go`.
- `fsrs_capped`: the FSRS long-term model with the pre-Phase 3 stability cap at
  36500 days.
- `fsrs_corrected`: the FSRS long-term model without the stability cap. The
  interval is still capped at 36500 days.

The capped model is implemented inside this harness. It is not production code
and nothing imports it. The corrected model calls `internal/fsrs` directly.

## Input

`fixtures.json` contains synthetic review histories only. Each review has:

- `gap_days`: days since the previous review.
- `rating`: FSRS rating 1 to 4.
- `quality`: SM-2 quality 0 to 5.

A fixture may provide `initial_stability` and `initial_difficulty` to start
above the reference cap. This is synthetic state, not a database row.

## Run

```sh
go run ./tools/fsrs-replay \
  -fixtures tools/fsrs-replay/fixtures.json \
  -out tools/fsrs-replay/report.json
```

Omit `-out` to write the report to standard output.

The harness has no network calls, no database access and no clock. The output
is deterministic.

## Output

`report.json` records the final due day, interval, stability and difficulty for
each strategy. It also records the difference between capped and corrected
FSRS, and between SM-2 and corrected FSRS.

The report is evidence for Phase 3. It is not a production scheduler result.

## Test

```sh
go test ./tools/fsrs-replay
```

The tests cover deterministic output, the stability-cap fixture, the
above-cap input fixture and invalid fixture rejection.
