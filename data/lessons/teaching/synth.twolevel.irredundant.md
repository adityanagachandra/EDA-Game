# Irredundant covers

After expansion every cube is prime, but the cover may still contain cubes that the others already cover. IRREDUNDANT removes them. It is not simply deleting every redundant cube: removing one can make another necessary. Espresso classifies the cubes and solves a small covering problem over the ones that are only partially redundant.

## Relatively essential cubes

A cube of a cover is **relatively essential** if it covers some minterm (outside the don't-care set) that no other cube of the cover covers. Such a cube must stay in every subcover. For the cover `-110`, `00-1`, `1-10`, `11-0`, `111-` of a function of $a, b, c, d$, each cube can be checked by asking whether it is contained in the union of the others plus the don't cares. That check is a tautology test on the cofactor of the other cubes with respect to the cube in question. Relatively essential primes are the cover-relative analogue of essential primes, but they are relative to the current cover, not to all primes.

## Redundant cubes and why they interact

A cube contained in the union of the others is **redundant**. Every cube is either relatively essential or redundant. Redundant cubes cannot all be removed at once: if cubes $A$ and $B$ cover the same leftover minterms, each is redundant while the other is present, but removing both uncovers those minterms. Espresso therefore splits the redundant cubes into totally redundant ones (covered by the relatively essential cubes and don't cares alone, safe to drop) and partially redundant ones (whose removal must be coordinated). In a cover of five primes where three are redundant, typically some are totally redundant and the rest compete for the same minterms.

## Minimum irredundant subcover

The partially redundant cubes and the minterms only they cover form a small covering problem: choose the fewest of them so that, with the relatively essential cubes, every ON minterm is covered. The result is an **irredundant cover**: no cube can be removed. For a cover of four primes such as `0--0`, `000-`, `1-11`, `11-1`, every cube turns out relatively essential and all four stay. Espresso solves this covering heuristically, since the subproblem is NP-hard in general but usually tiny. Prime plus irredundant is the local optimality Espresso guarantees: a cover can satisfy both properties and still not be minimum.

## Going deeper

Brayton et al., "Logic Minimization Algorithms for VLSI Synthesis" (1984), on IRREDUNDANT_COVER; Hachtel and Somenzi (1996), on heuristic two-level minimization.
