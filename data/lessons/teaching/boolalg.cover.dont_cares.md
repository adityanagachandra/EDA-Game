# Incompletely specified functions

Real specifications leave some outputs free. Inputs that never occur (an unused opcode, an unreachable state code) and outputs nobody observes are don't cares. An incompletely specified function has an ON-set, an OFF-set and a don't-care set, and any completely specified function that agrees on the care rows implements it. Don't cares are the main lever of logic optimization: the more a tool can prove irrelevant, the more freedom its minimizer has. Getting them wrong, by assuming a don't care that is reachable, is a silent functional bug.

## The ON, OFF and don't-care sets

Write the function as a triple $(F_{ON}, F_{DC}, F_{OFF})$ that partitions the $2^n$ rows. An implementation $g$ must satisfy $F_{ON} \le g \le F_{ON} + F_{DC}$, so each don't-care row may be set either way, and exactly $2^{|DC|}$ completely specified functions are compatible.

For $f(a, b, c)$ with ON-set $m(1, 3, 5)$ and don't-care set $d(7)$, there are $2^1 = 2$ implementations: $m(1, 3, 5)$ and $m(1, 3, 5, 7)$. With three don't cares there would be 8. In practice the don't-care set is not listed but derived: satisfiability don't cares from impossible input combinations of a node, observability don't cares from the Boolean difference of the outputs with respect to it.

## Minimizing with don't cares

Minimization uses don't cares on the way up and ignores them on the way down. Primes are grown over $F_{ON} + F_{DC}$, so cubes may swallow don't-care rows, but the covering step only requires the ON rows to be covered.

In the example, without the don't care the best cover is $\overline{a}c + \overline{b}c$, 2 cubes. Allowing row 7 lets both merge into the single prime $c$, which covers rows 1, 3, 5 and 7, so the minimum drops to 1 cube. Espresso implements this by passing the don't-care cover to EXPAND (cubes may grow into it) and dropping don't-care minterms from the covering requirement.

## Primes of an incompletely specified function

The primes of an incompletely specified function are the primes of $F_{ON} + F_{DC}$. Some of them lie entirely in the don't-care set and cover no ON row; they are useless for covering and are discarded.

For ON-set $m(1, 3, 5)$ and don't-care set $d(7)$, the only prime of $m(1, 3, 5, 7)$ is $c$, and it covers ON rows, so one useful prime. With ON-set $m(0)$ and don't-care set $d(6, 7)$, the primes of $m(0, 6, 7)$ are $\overline{a}\,\overline{b}\,\overline{c}$ and $ab$, and only the first touches the ON-set. Discarding the others before covering keeps the table small. The flexibility is real, but the minimizer must still prove that the don't cares are don't cares, which in multi-level synthesis is itself a SAT or BDD computation.

## Going deeper

Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984). Savoj and Brayton, "The use of observability and external don't cares for the simplification of multi-level networks", DAC (1990).
