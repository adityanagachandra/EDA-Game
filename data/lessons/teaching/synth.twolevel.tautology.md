# Tautology checking by cofactoring

Most operations inside Espresso reduce to one question: is this cover a tautology, equal to 1 everywhere? Containment of a cube in a cover, redundancy of a cube, and equivalence of two covers are all tautology checks on cofactors. Espresso answers it with the unate recursive paradigm: split by Shannon cofactoring until each piece is trivially decided, using unateness to stop early.

## Cofactors of a cover

The cofactor of a cover $F$ with respect to $x = 1$ keeps every cube that does not contain $\bar{x}$ and raises $x$ in it to a dash; cubes containing $\bar{x}$ vanish. For the cover `10-`, `00-`, `-1-` over $a, b, c$, the cofactor with respect to $b = 0$ keeps `10-` and `00-` as `1--` and `0--` and drops `-1-`, leaving 2 cubes. More generally the cofactor with respect to a cube $c$ restricts the cover to the subspace of $c$. By Shannon expansion $F = x F_x + \bar{x} F_{\bar{x}}$, so $F$ is a tautology exactly when both cofactors are.

## Recursive tautology checking

Recursion splits on a variable, usually the most binate one (appearing in both polarities in the most cubes), and checks both cofactors. It stops at easy cases: a cover containing the all-dash cube is a tautology; an empty cover is not; a cover in which some column is all one polarity can be decided by unate rules; and a cover whose total minterm count is below $2^n$ cannot be a tautology. When a branch fails, the path of assignments to the failing leaf is a minterm the cover misses. Counting all uncovered minterms is a stronger question than tautology, but on small covers it can be read off the same recursion.

## Unate covers end the recursion

A cover is **unate** in a variable if that variable appears in only one polarity (or not at all). If a cover is unate in every variable, it is a tautology exactly when it contains the all-dash cube. The reason: assign each variable the value opposite to the polarity in which it appears. Every cube other than the all-dash cube has at least one literal, each such literal is now false, so the cover evaluates to 0 at that point. This is why Espresso's recursion picks binate variables to split on: each split pushes the cofactors toward unateness, where a single containment check decides them. The same paradigm computes complements and checks containment, which is what makes Espresso work on cube lists with dozens of inputs where truth tables would be impossible.

## Going deeper

Brayton et al., "Logic Minimization Algorithms for VLSI Synthesis" (1984), on the unate recursive paradigm.
