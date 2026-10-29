# Tautology checking of a cover

A cover is a tautology when its cubes together contain every minterm, so the function it represents is the constant 1. Tautology checking is the workhorse of cube-based logic minimization: whether a cube is redundant, whether a cube can be expanded, and whether two covers are equivalent all reduce to it. It is coNP-complete in general, and the unate recursive paradigm of Espresso is the standard way to make it fast in practice.

## A cover is a tautology when nothing is missed

A cover $F$ over $n$ variables is a tautology exactly when the union of its cubes is all $2^n$ minterms. For small $n$ you can count: union the cubes and count the rows they miss. Zero misses means tautology.

Over $a, b, c$ let F = {1--, -1-, --1}, that is $a + b + c$. The only minterm with $a = b = c = 0$ is missed, so $F$ misses 1 minterm and is not a tautology. Adding the cube 000 would make it one. Counting does not scale, which is why the check is done recursively.

## Cofactoring a cover

The cofactor of a cover with respect to a literal is computed cube by cube. A cube that conflicts with the literal is dropped; a cube that contains the literal or has a dash there survives with that position set to dash. $F$ is a tautology exactly when both $F_x$ and $F_{\overline{x}}$ are, by Shannon expansion.

For F = {1--, -1-, --1} and literal $\overline{a}$: the cube 1-- conflicts and is dropped, and the others survive, so $F_{\overline{a}}$ = {-1-, --1}, 2 cubes. The recursion ends early on special cases: a cover with an all-dash cube is a tautology; a cover that is unate in every variable is a tautology only if it contains the all-dash cube; and a column of all 0s or all 1s means some minterm is missed. Splitting on the most binate variable reaches these cases fastest.

## Cube containment as a tautology question

Cube $c$ is contained in cover $F$ (every minterm of $c$ is covered) exactly when the cofactor $F_c$ is a tautology. This turns the questions a minimizer asks into tautology checks: a cube of a cover is redundant if it is contained in the other cubes, and a cube can be expanded if its larger version is contained in ON plus don't cares.

Is 00- contained in $F$ above? Cofactor with respect to $\overline{a}\,\overline{b}$: 1-- and -1- conflict and drop, leaving {--1}, which is not a tautology. The cube 000 inside 00- is the missed minterm. Each containment check costs one recursive tautology call, and Espresso makes thousands of them per run.

## Going deeper

Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984), chapter 3. Rudell and Sangiovanni-Vincentelli, "Multiple-valued minimization for PLA optimization", IEEE TCAD (1987).
