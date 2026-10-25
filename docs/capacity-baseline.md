# Capacity baseline: Phase 0

Date: 2026-10-25.

This report records measurements against the deployed single-process
architecture. It does not describe a proposed deployment.

## Verdict

Phase 0 is a **conditional fail** as originally measured.

The harness is reproducible and the measured latency and error data are usable.
The retry-safety gate failed: an exact duplicate study answer recorded a second
attempt and awarded XP twice. The durability gate also failed: progress, attempt
and XP were written separately, with no shared transaction.

The Phase 1 section below records the repair. Retry safety and the core answer
transaction now pass. Best-effort effects outside that transaction are listed
as remaining risks.

The measurements do not justify Redis, multiple API replicas or a database
migration for throughput. They justify fixing idempotency and write atomicity
before any architecture change.

## Environment

- Go server built from `./cmd/mathua`.
- SQLite scratch database at `/tmp/opencode/phase0.db`.
- Postgres 16 container on `127.0.0.1:55432` for the comparison run.
- `MATHUA_LOADTEST=1` on the scratch server only. This relaxes rate limiters and
  must never be used in production.
- Load driver: `barrage v0.5.5`.
- Numeric workload: guest creation, one practice fetch, one study answer and a
  scores read. Read workload: graph and leaderboard.
- Symbolic workload: `alg.poly.division` with an equivalent answer, so SymPy is
  invoked instead of the exact-match fast path.

## Reproduce the numeric baseline

```sh
go build -o /tmp/opencode/mathua-loadtest ./cmd/mathua
rm -f /tmp/opencode/phase0.db
DATABASE_URL=/tmp/opencode/phase0.db MATHUA_LOADTEST=1 \
  /tmp/opencode/mathua-loadtest -serve -port 18082 &
mise exec -- barrage run -c loadtest/barrage-study-loop.yaml \
  --duration 45s --concurrency 40 --no-progress \
  --json loadtest/phase0-c40.json --no-report
```

`MATHUA_LOADTEST=1` relaxes the rate limiters. Never set it on a production
server. The SQLite file must be a scratch file, never the production volume.

The symbolic run needs a guest with every prerequisite already mastered. It is
not one command because the current API has no test-only seeding endpoint.

## Results

### SQLite, numeric study loop

| Concurrency | Scenario requests/s | Approx. answers/s | p50 | p95 | p99 | Success |
|---:|---:|---:|---:|---:|---:|---:|
| 20 | 518.5 | 129.6 | 27.1 ms | 41.8 ms | 59.9 ms | 99.9% |
| 40 | 314.4 | 78.6 | 100.2 ms | 145.8 ms | 170.6 ms | 99.774% |
| 80 | 99.5 | 24.9 | 591.1 ms | 1,119.0 ms | 1,307.7 ms | 98.66% |

The read-browse workload stayed between p95 63.6 ms and 78.9 ms. The SQLite
database runner stayed between p95 9.35 ms and 13.15 ms.

The answer path collapses between concurrency 40 and 80. Throughput falls while
latency rises. Read latency stays low, so a saturated read database is not the
cause.

### PostgreSQL, numeric study loop

| Concurrency | Scenario requests/s | Approx. answers/s | p50 | p95 | p99 | Success |
|---:|---:|---:|---:|---:|---:|---:|
| 40 | 261.1 | 65.3 | 112.7 ms | 184.8 ms | 246.7 ms | 99.736% |
| 80 | 83.4 | 20.9 | 649.0 ms | 852.8 ms | 1,271.8 ms | 98.561% |

Postgres is slower than SQLite in this single-process test. That points to
application serialization or CPU work, not the database engine alone.
`submitAnswerWithTask` takes `e.mu` for study session lookup and creation. That
lock is a candidate and has not been profiled.

### Symbolic grading

| Concurrency | Scenario requests/s | Approx. answers/s | p50 | p95 | p99 | Success |
|---:|---:|---:|---:|---:|---:|---:|
| 10 | 203.1 | 101.6 | 41.2 ms | 86.2 ms | 166.2 ms | 99.8% |
| 20 | 118.0 | 59.0 | 155.7 ms | 271.5 ms | 551.9 ms | 99.2% |

The symbolic path is materially slower than the numeric path. At 10% symbolic
traffic, 2,500 answers per second would require about 250 symbolic grades per
second. The measured rate is 59 to 102 per second on one process.

## Error interpretation

Every `POST /api/study/answer` 503 occurred in the final second of a run. Their
response times were 5 to 164 ms, not grading timeouts. They are Barrage
end-of-run cancellations handled by `http.TimeoutHandler`, not application
failures.

The SymPy `ModuleNotFoundError` lines appeared only at server startup, before
the listener started. They are not request-time failures in these runs.

Outside the final second, the application returned no 5xx responses in the
measured runs.

## Retry and concurrency probe

Procedure:

1. Create a guest.
2. Fetch one `arith.add.single` practice question.
3. Submit the exact same body twice.

Result:

- First response: XP 1, status LEARNING.
- Second response: XP 1, status PRACTICING.
- Scores: `xp_total` 2.
- Database: two attempts and two XP events for the same submission.

At Phase 0, the Learn path was not retry-safe. `AwardXP` had the right unique
constraint, but the caller generates a new attempt ID on every request, so the
constraint never sees the same key twice.

## Phase 1 write path

The study answer now records one durable submission row, one attempt, the
progress row, the topic speed row and the XP event in a single transaction.
The exact same submission id replays the stored response. The storage
transaction itself benchmarks at about 93 microseconds per commit on SQLite in
`BenchmarkStudySubmissionCommit`.

Final-code measurements on a scratch database:

| Store | Concurrency | Scenario requests/s | p50 | p95 | p99 | Success |
|---|---:|---:|---:|---:|---:|---:|
| SQLite | 40 | 311.7 | 94 ms | 156 ms | 193 ms | 99.67% |
| SQLite | 80 | 175.7 | 317 ms | 642 ms | 788 ms | 98.86% |
| PostgreSQL | 40 | 266.9 | 99 ms | 141 ms | 175 ms | 99.65% |

Phase 0 comparison:

- SQLite c40: 314.4 scenario requests/s and p95 145.8 ms before, 311.7 and
  156 ms after.
- SQLite c80: 99.5 scenario requests/s and p95 1,119 ms before, 175.7 and
  642 ms after.
- PostgreSQL c40: 261.1 scenario requests/s and p95 184.8 ms before, 266.9 and
  141 ms after.

The c80 PostgreSQL run was discarded: the scratch server was shut down during
the measurement and Barrage recorded connection failures. It is not evidence.

Phase 1 gates:

| Gate | Result |
|---|---|
| Same submission id writes one attempt and one XP event | Pass |
| Concurrent same id writes one logical submission | Pass |
| Same id with a different payload is refused | Pass |
| Failure after attempt, progress or XP insert rolls back | Pass on SQLite and PostgreSQL |
| Replay returns the original body including halted | Pass on SQLite and PostgreSQL |
| Missing submission id is refused | Pass |
| In-memory KP state advances only after commit | Pass |
| Concurrent duplicate is one submission or an explicit 409 | Pass on SQLite and PostgreSQL |
| Core answer writes are atomic | Pass |
| Post-commit weakness, credit and provisional fallback | Not atomic, best effort |
| Idempotency retention policy | Defined, compaction not implemented |

Post-commit effects were classified by learner-visible impact:

| Effect | Visible impact | Status |
|---|---|---|
| Weakness propagation | Recommendation order | Best effort, recomputed from progress |
| Implicit credit | Credit toward encompassing concepts | Idempotent by attempt key, best effort |
| Provisional prerequisite demotion | Prerequisite status and study order | Best effort, not transactional |
| Remedial queue | Next study target | Best effort after commit |
| In-memory miss and rush counters | Halt and rush penalty on later answers | Process local, lost on restart |

Only the answer verdict, progress, mastery, attempt, XP and replay body are in
the transaction. The remaining effects are not described as atomic.

The collapse between c40 and c80 remains. It is smaller on SQLite and it is no
longer a write-atomicity failure, but the answer path still serializes on the
engine and on the SQLite writer.

### Implemented answer-path fixes

Two narrow fixes were implemented after this baseline:

- `GzipMiddleware` now buffers the first 1 KiB and compresses only larger,
  compressible responses. It pools gzip writers and negotiates `Accept-Encoding`
  including `q=0`.
- The leaderboard predicate now excludes rows with no password, email
  or OAuth identity, instead of relying on a `guest_` id prefix. A UUID guest is
  excluded, and an adopted guest with credentials stays ranked.

A btrfs-backed c80 rerun used the same mixed scenario and a study-only scenario.
Five mixed trials and three study-only trials were collected.

| Workload | Before | After | Change |
|---|---:|---:|---:|
| Study-only requests/s, median | 458.6 | 519.5 | +13.3% |
| Study-only p95 | 278 ms | 255 ms | -8.3% |
| Study-only process CPU | 361.7% | 241.6% | -33.2% |
| Study-only GC lines | 2,256 | 513 | -77.3% |
| Mixed requests/s, median | 141.4 | 207.9 | +47.0% |
| Mixed p95, median | 737 ms | 494 ms | -33.0% |
| Mixed process CPU, median | 872.3% | 779.8% | -10.6% |
| Mixed GC lines, median | 2,627 | 819 | -68.8% |

Mixed throughput after the corrected leaderboard predicate ranged from 170.1 to
245.1 requests/s. The median improved by 47% and the range is above every
before trial. RSS rose from about 139-143 MiB to about 155-169 MiB.

The full eligibility optimization was intentionally omitted. Replacing
`EligibilityFor` would duplicate the scheduler's prerequisite and placement
rules for a measured gain of about 4% throughput. The safe path is a later
targeted lookup with an equivalence test against the scheduler.


The storage transaction is not the c80 limit. `BenchmarkStudySubmissionCommit`
measured about 87 microseconds per commit with one CPU and about 90
microseconds with eight CPUs, so the SQLite writer serializes near 11,000
commits per second. The observed c80 rate is about 152 scenario requests per
second, far below that ceiling. The remaining collapse is above storage, in the
answer handler and shared engine state.

## Bottleneck conclusion

The measured bottleneck is the application answer path, not read traffic and not
Postgres. The likely contributors are:

- the engine mutex around study session lookup and creation
- synchronous per-answer writes
- repeated progress and mastery calculations
- symbolic grading in the same process

Profiling is required before choosing a fix. A larger machine would raise the
ceiling, but it would not make retries safe.

## Phase 0 gates, as measured before the repair

| Gate | Result |
|---|---|
| Reproducible current-API scenario | Pass |
| p50, p95 and p99 captured | Pass |
| Error rates captured and explained | Pass |
| Database read load separated | Pass |
| SymPy path exercised | Pass |
| Exact retry cannot duplicate XP | Fail |
| No lost attempt between progress, attempt and XP | Fail |
| Postgres transaction and pool audit complete | Partial, measured but not fully audited |

## Artifacts

- `loadtest/barrage-study-loop.yaml`
- `loadtest/phase0-c20.json`
- `loadtest/phase0-c40.json`
- `loadtest/phase0-c80.json`
- `loadtest/phase0-pg-c40.json`
- `loadtest/phase0-pg-c80.json`
- `loadtest/phase0-symbolic-c10.json`
- `loadtest/phase0-symbolic-c20.json`
- `loadtest/phase1-c40.json`
- `loadtest/phase1-c80.json`

The scratch databases and containers are disposable and must not be pointed at
production.

## Remaining before Phase 2

1. Profile the answer path to identify the lock or CPU work behind the c40 to c80 collapse.
2. Move weakness propagation, implicit credit and provisional fallback to a
   durable outbox or accept their best-effort status.
3. Re-run the baseline against a production-sized database.
4. Implement the `answer_submissions` compaction policy described in `capacity-plan.md`.
5. Run the additive migration on a copy of production data before any deploy.
