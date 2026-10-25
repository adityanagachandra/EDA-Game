# Durable outbox for study-answer post-commit effects

Status: open. This is a separate reliability task from the committed atomic
answer repair in `d0d411d4`.

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
