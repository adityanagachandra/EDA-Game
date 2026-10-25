# Durable outbox for study-answer post-commit effects

Status: implemented in the working tree, not committed. Correctness checks pass
and the c80 performance gate is within 5.4% of the pre-outbox baseline. Commit is
recommended. This work is separate from the committed atomic answer repair in
`d0d411d4`.

## Why this exists

The committed repair makes the answer verdict, attempt, progress, mastery, XP
and replay body atomic. It does not make the effects that run after that commit
durable:

- weakness propagation
- implicit credit allocation
- provisional prerequisite demotion
- remedial queue insertion
- process-local miss and rush counters

A process failure between the answer commit and those effects can leave the
learner with an accepted answer but stale recommendations, missing credit or a
missing remedial target. The current code logs these failures and continues.

## Implementation

`study_answer_effects` is written in the same transaction as the answer. Each
row carries the submission id, effect kind and JSON payload. `ProcessStudyEffects`
claims a bounded batch, applies each effect and marks it complete. A failed
effect stays pending with exponential backoff. The server starts a worker at
startup and drains every five seconds.

Effects are idempotent: weakness is recomputed, credit uses the attempt key,
remedial insertion dedupes, provisional demotion is repeatable and removal is
repeatable.

## Performance gate

The first version drained effects synchronously inside the answer request. That
cost 33.5% of study-only throughput. Moving the drain to the background worker
recovered almost all of it.

c80 study-only on btrfs disk, three trials each:

| Variant | Requests/s | p95 | Answer 200s | CPU | GC | Max RSS |
|---|---:|---:|---:|---:|---:|---:|
| Pre-outbox base | 533.3 | 249 ms | 10,603 | 239.8% | 512 | 174,708 KB |
| Synchronous drain | 354.7 | 405 ms | 7,025 | 199.9% | 411 | 154,356 KB |
| Durable queue, no sync drain | 504.6 | 249 ms | 10,030 | 231.5% | 492 | 169,748 KB |

The remaining 5.4% is the cost of the two to five effect inserts in the answer
transaction. Batching those inserts is a separate optimization.

The background worker drains immediately at startup and every five seconds. It
does not rely on an in-memory signal for durability, so a missed wakeup delays
an effect but cannot lose it.

## Required behavior

The answer transaction must also write a durable outbox record for every
post-commit effect that must eventually happen.

- the outbox record commits with the answer
- a worker processes pending records after commit
- processing is idempotent by submission id and effect kind
- a failed effect stays pending and is retried
- startup recovers pending records before serving new work
- a replay returns the stored answer body and does not duplicate an effect

## Acceptance tests

- commit succeeds and the worker fails, then a restart processes the record
- two workers claim one record without applying it twice
- implicit credit is applied once for one submission
- provisional demotion and remedial insertion survive a restart
- weak propagation is recomputed without overwriting newer progress
- the replay body is unchanged by outbox processing
- SQLite and PostgreSQL both pass

## Non-goals

This task does not solve the answer-handler throughput collapse. That is a
separate performance investigation.

This task is not required for the correctness of the answer record itself. It
is required before claiming crash-safe post-commit effects.
