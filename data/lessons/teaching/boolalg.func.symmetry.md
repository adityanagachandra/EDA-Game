# Symmetric functions

A function is symmetric in a pair of variables when swapping them never changes the output, and totally symmetric when any permutation of its inputs leaves it unchanged. Symmetry is structure a tool can exploit: synthesis builds symmetric functions with counter-based circuits, technology mappers treat symmetric pins of a cell as interchangeable, BDD reordering keeps symmetric variables adjacent, and Boolean matching uses symmetry classes to cut down the permutations it must try.

## Totally symmetric functions

A totally symmetric function depends only on the weight of the input, the number of inputs at 1. It is described by the set $W$ of weights that make it 1. AND of $n$ inputs has $W = \{n\}$, OR has $W = \{1, \dots, n\}$, parity has the odd weights, and the majority of three inputs has $W = \{2, 3\}$.

The ON count is the sum of binomial coefficients over $W$, $\sum_{k \in W} \binom{n}{k}$. For majority of three, $\binom{3}{2} + \binom{3}{3} = 3 + 1 = 4$ ON rows. For a five-input function that is 1 at weights 1 and 4, $\binom{5}{1} + \binom{5}{4} = 10$.

## Symmetry in a pair of variables

$f$ is symmetric in $x$ and $y$ exactly when $f_{x\overline{y}} = f_{\overline{x}y}$: the two rows that differ by swapping $x$ and $y$ (01 against 10) always agree, and rows with $x = y$ are unchanged by the swap anyway.

For $f(a, b, c) = ab + c$: $f_{a\overline{b}} = c$ and $f_{\overline{a}b} = c$, so $f$ is symmetric in $(a, b)$. For $(a, c)$: $f_{a\overline{c}} = b$ but $f_{\overline{a}c} = 1$, so it is not. Pairwise symmetry is an equivalence relation on the variables (for the plain, non-skew form), so its classes can be found with at most $\binom{n}{2}$ cofactor comparisons. Variants include skew symmetry, where swapping with complementation leaves $f$ unchanged.

## Counting symmetric functions

Because a totally symmetric function is a choice of output for each weight $0, 1, \dots, n$, there are exactly $2^{n+1}$ of them, against $2^{2^n}$ functions in total. For $n = 3$ that is 16 of 256. Fixing outputs at some weights halves the count per fixed weight: those with $f(0 \dots 0) = 0$ number $2^n$.

This count is why symmetric functions are cheap: $n + 1$ bits describe them, and a circuit can compute the weight with an adder tree and then decode it. It also means a random function is almost never symmetric, so a matcher that detects symmetry has found something worth using.

## Going deeper

Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", chapter on symmetries. Mishchenko, "Fast computation of symmetries in Boolean functions", IEEE TCAD (2003).
