# Test compaction

ATPG generates one test cube per targeted fault, with many inputs left as don't-cares. Tester time and memory are proportional to the number of patterns, so the cubes are merged and redundant tests are removed. Compaction is a set-cover problem in disguise, solved with greedy heuristics; it routinely shrinks pattern counts several-fold and interacts directly with scan compression.

## Merging test cubes

A **test cube** specifies only the inputs needed for its fault, with X elsewhere. Two cubes are **compatible** if no position has 0 in one and 1 in the other; compatible cubes merge into one pattern by taking the specified bit wherever either specifies one. **Static compaction** merges after generation: process cubes in order and merge each into the first compatible merged pattern, else start a new one. For X0X1X0, 0X1XXX, X0XXXX, XXX0XX and X0XXXX, the first three merge into 0011X0, the fourth conflicts with it at position 4 and starts a new pattern, and the fifth merges into the first: 2 patterns. The result depends on order, and finding the minimum is NP-hard. **Dynamic compaction** does better by asking ATPG to target extra faults within the don't-cares of the current cube before releasing it.

## Minimum covering sets

Given a fault-detection table from fault simulation, the smallest subset of tests that detects every detected fault is a **set cover** problem. Exact solutions use ILP for small instances; production tools use greedy selection (repeatedly take the test detecting the most uncovered faults) or iterative methods. Essential tests, the only ones detecting some fault, must always be kept.

## Reverse-order fault simulation

A cheap, effective trick: fault-simulate the tests in **reverse** generation order with fault dropping and keep a test only if it detects a fault not yet detected. Late tests target hard faults and incidentally detect many easy ones, so early tests often become redundant. For tests whose detection sets are T6 = {f1, f4, f5, f7}, T5 = {f1, f2, f4, f5, f7}, T2 = {f6} and others subsumed, the reverse pass keeps T6, T5 and T2 and drops the rest. Compaction trades against compression: fully specified merged patterns have more care bits, which a linear decompressor may fail to encode.

## Going deeper

Bushnell and Agrawal (2000), chapter 7; Pomeranz, Reddy and others on static compaction; Wang, Wu and Wen (2006), chapter 4.
