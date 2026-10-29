# Golden tests on coverage and pattern counts

When an agent changes an ATPG or DFT skill, the evidence that it helped and did not hurt is a golden suite: reference designs with known baseline coverage and pattern counts, and tolerances the new run must meet. Writing the checks precisely, per design and per metric, is what turns an agent's claim into something a reviewer can approve.

## Per-design limits

A golden test fixes, for each reference design, the baseline test coverage and pattern count, and a tolerance for each: for example, coverage may not drop by more than 0.05 percentage points and pattern count may not grow by more than 2 percent. A design fails if either limit is violated, and the change fails if any design fails. A design whose coverage moved from 97.81 to 97.68 percent fails on coverage even if its pattern count fell; one whose pattern count grew from 9073 to 9798 (8 percent) fails even though coverage barely moved. Gains on some designs never excuse a regression on another.

## Computing the thresholds

Ceilings are computed exactly so that the check is deterministic. With a baseline of 19,657 patterns and a 2 percent allowance, the largest passing count is $\lfloor 19{,}657 \times 1.02 \rfloor = 20{,}050$. Coverage is recomputed from raw counts rather than taken from a rounded report: with $T$ faults, $D$ detected and $U$ proven untestable, $TC = D/(T - U)$. Comparing to three decimals avoids flapping on rounding. Runs must be seeded and deterministic, or the suite will fail and pass at random.

## What the metric can hide

Golden metrics can be gamed. Test coverage rises when more faults are proven untestable even if fewer are detected: baseline 55,714 detected of 58,707 with 1761 untestable gives 97.84 percent; a run with 55,599 detected and 1941 untestable gives 97.94 percent, an apparent gain from detecting 115 fewer faults. A robust suite therefore checks detected counts and fault-list identity alongside coverage, verifies new untestable proofs, and adds a differential check that patterns for unchanged fault classes still detect the same faults.

## Going deeper

The agent-first EDA notes in this curriculum (golden and differential testing, tool-grounded rewards).
