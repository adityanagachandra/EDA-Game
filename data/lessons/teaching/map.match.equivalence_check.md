# Deciding NPN equivalence of two functions

Given a cut function and a candidate cell, the matcher must decide whether some NPN transformation turns one into the other, and find it. The decision combines cheap invariants that refute most pairs with a constructive search that confirms the rest. Doing it by hand for 3-input functions builds the intuition a matcher designer needs about which invariants are strong.

## Necessary conditions first

Two NPN-equivalent functions have the same number of ones or complementary numbers (j and $2^n - j$), the same multiset of cofactor ones counts up to polarity, and the same symmetry structure. Check these before searching. f with ones at {1, 3, 4, 5, 7} has 5 ones, so any equivalent function has 5 or 3 ones. Among candidates with 3 ones, compare cofactor counts: f's complement has 3 ones, and its counts per variable must match the candidate's up to swapping the two cofactors of a variable. All of these tests cost a few bit operations per variable.

## Constructing the transformation

When invariants agree, search for the transformation. Fix the output polarity from the ones count: if the counts differ, the output must be negated. Then pair variables whose cofactor counts agree and read each input's polarity from which cofactor matches. Example: f has ones {1, 5}, which is x1 AND NOT x2 with x3 unused, and g has ones {2, 3}. Both have 2 ones, so the output is not negated. In f, x1 has cofactor counts (positive 2, negative 0), x2 has (0, 2) and x3 has (1, 1); in g, x2 has (2, 0), x3 has (0, 2) and x1 has (1, 1). Pairing them gives g(x1, x2, x3) = f(x2, x3, x1), and a truth-table check confirms it. Equivalently, compare canonical forms: equal forms mean equivalent, and the two canonicalizing transformations compose into the pin assignment.

## Why same counts are not enough

Invariants are necessary, not sufficient. For 3 inputs, 6 of the 14 NPN classes contain functions with exactly 4 ones, among them XOR of two inputs, a single literal, the multiplexer and x1 XOR (x2 AND x3)-like functions, so the ones count alone cannot separate them. Cofactor counts separate many but not all pairs; for larger n, distinct classes can share every simple signature. A matcher must therefore end with an exact check, either by canonical form or by a transformation search followed by a truth-table comparison. Skipping it produces wrong matches, and the mapped netlist then fails equivalence checking.

## Going deeper

Mishchenko, Chatterjee and Brayton, "Technology mapping with Boolean matching, supergates and choices", UC Berkeley technical report, 2005; Benini and De Micheli, "A survey of Boolean matching techniques for library binding", ACM TODAES 1997.
