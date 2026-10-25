# Capacity plan

Status: Phase 0, measurement only. This document does not authorize a storage
migration, Redis or multiple API replicas.

## 1. Target and assumptions

The planning target is 50,000 concurrent learners. Concurrency alone does not
give an answer rate. The answer rate is:

```
answers/second = concurrent learners / seconds between answers
```

At one answer every 20 seconds, 50,000 concurrent learners produce 2,500
answers per second. At one answer every 30 seconds they produce about 1,667
answers per second. A learner who is reading a lesson, reviewing a profile or
waiting between questions answers more slowly.

The Phase 0 measurements below are expressed in scenario requests per second and
in answer submissions per second. They must not be read as proof that 50,000
concurrent learners are supported.

The earlier estimate of 5,000 to 12,000 writes per second is not a measured
requirement and is not a Phase 0 gate. Phase 0 counts actual writes per answer
and reports them instead of treating the estimate as fact.

## 2. Phase 0 scope

Phase 0 answers four questions.

1. Can the current APIs be driven by a reproducible load test?
2. Where do latency and errors actually appear?
3. Is the current answer write path safe under retries and concurrent
   submissions?
4. Does Postgres remove a measured bottleneck or only move it?

Phase 0 does not change the production architecture. It does not migrate
SQLite, add Redis or add replicas. Those changes need evidence from this phase.

## 3. Audit findings

### 3.1 Study answer path

`POST /api/study/answer` calls `Engine.SubmitStudyAnswer`, which performs three
separate operations:

1. `applyGradedAttempt` updates progress, mastery evidence, SM-2, weakness and
   other learner state.
2. `RecordAttempt` inserts an attempts row.
3. `AwardXP` inserts an XP event keyed by `attempt:<attemptID>`.

The attempt ID is built from student, concept and the current nanosecond. Every
request creates a new ID. The client sends `retryAttempts: 1` and its comment
claims the engine replays an existing attempt ID. It does not. An exact retry is
therefore graded again and writes a second attempt and a second XP event.

`AwardXP` is idempotent only when the caller supplies the same event key. The
database has the right unique constraint. The caller does not supply a stable
key.

No transaction spans progress, attempt and XP. A crash between calls can leave
progress updated without an attempt, or an attempt without its XP event.

### 3.2 Durability rule

An accepted answer must be durably recorded before the API returns XP and a
verdict. The durable record is the attempts row plus the learner state changed
by that answer.

XP, progress and the attempt row must commit together. Alternatively the
attempt and an outbox record must commit together and a worker must finish the
remaining effects. A cache may hold the live session, but it must never be the
only record of an accepted answer, an XP award or a mastery change.

The synchronous set is the attempt row, the XP event, the progress and mastery
change and the idempotency key. These affect the next recommendation and the
learner's visible XP, so they cannot be eventually consistent. Activity feed
entries, leaderboard aggregates, analytics and notification email may be
asynchronous.

Every answer request must carry a stable idempotency key. A retry with the same
key returns the original result and writes nothing new. The current study path
does not satisfy this rule.

### 3.3 Idempotency retention

An idempotency key may be compacted only after its response is no longer needed.
It must never be deleted while the attempt and XP rows it guards still exist,
because deleting the key reopens the duplicate-award window.

The safe policy is:

- keep the full response for at least the client retry window, currently seven days
- after that, keep the key, request hash and attempt id as a tombstone
- do not delete tombstones while the guarded attempt exists
- if a tombstone is later pruned, prune it only together with the attempt and XP
  rows it guards, never on its own

A request that arrives after compaction receives an explicit expired-submission
response instead of being graded again.

### 3.4 Database behaviour

SQLite runs with WAL, foreign keys, a five second busy timeout and eight open
connections. Postgres uses twenty-five open connections and five idle
connections. Postgres supports the same repository interface.

The answer path is not one transaction. Several storage methods commit
separately. Postgres removes SQLite's single-writer limitation but does not by
itself make the application write path atomic.

### 3.5 Shared state

Live quiz sessions, diagnostic sessions and rate limiters are held in process
memory. A second API replica cannot resume a session created by another replica.
A restart loses live sessions. This is a scaling blocker, but Phase 0 measured
the current single-process path first.

### 3.6 Grading

Most grading is Go. Symbolic grading can call SymPy after an exact-match fast
path. The local probe needed a deliberately equivalent expression to exercise
SymPy. Submitting the exact expected string does not.

## 4. Load-test scenario

`loadtest/barrage-study-loop.yaml` now drives the current APIs:

- create a guest
- fetch `GET /api/lessons/arith.add.single/practice`
- submit one `POST /api/study/answer`
- read `GET /api/scores/{student}`
- read the graph and leaderboard in a separate browse scenario

The guest is created inside each journey because Barrage has no setup phase.
The result is a new-learner journey, not a returning learner answering twenty
questions. The SQLite database runner adds the read shapes used by profile,
history and recommendation pages.

The symbolic probe uses `alg.poly.division`, submits an equivalent expression
such as `(answer)` and therefore exercises SymPy instead of the exact-match fast
path.

## 5. Phase 0 gates

Phase 0 passes only when all of these are true.

| Gate | Result |
|---|---|
| Reproducible scenario against current APIs | Pass |
| Measured p50, p95 and p99 with error rates | Pass |
| Database read work separated from answer writes | Pass |
| Grading path measured with SymPy actually invoked | Pass |
| Retry or concurrent duplicate cannot duplicate XP | Fail |
| No lost attempt between progress, attempt and XP | Fail, no shared transaction |
| Postgres transaction and pool audit complete | Partial |
| Baseline report states pass or fail explicitly | Pass, this document and `capacity-baseline.md` |

The retry and durability failures are not acceptable for 50,000 concurrent
learners. They are also present at the current traffic level.

## 6. What the measurements justify

The measurements justify:

- profiling the answer path before any infrastructure change
- adding a stable attempt or idempotency key
- making progress, attempt and XP one atomic operation
- auditing the Postgres write path with concurrent submissions
- keeping the current single-process architecture for now

The measurements do not yet justify:

- Redis
- multiple API replicas
- a SQLite to Postgres migration for throughput alone
- a larger machine as the primary fix

Postgres may still be required for availability and backup recovery. That is a
reliability decision, not a throughput conclusion from this phase.

## 7. Next phase conditions

Phase 1 starts only after Phase 0 fixes or explicitly accepts the retry and
durability gaps.

The first engineering work is:

1. Add a stable idempotency key to the study answer request.
2. Make the attempt write and XP award idempotent under the same key.
3. Audit the Postgres transaction boundaries for progress, attempt and XP.
4. Profile the in-process answer path and confirm which lock or CPU work causes
   the collapse between concurrency 40 and 80.
5. Repeat the baseline against a production-sized database.

Only then should the plan decide between a larger machine, Postgres, Redis or
replicas.
