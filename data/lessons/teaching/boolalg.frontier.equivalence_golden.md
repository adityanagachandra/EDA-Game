# Equivalence checking as the golden test

When an agent proposes a logic transform, a rewrite, a resynthesized cone, a remapped block, the golden test is a combinational equivalence check between the original and the proposal. It is the Boolean domain's version of a golden regression: deterministic, complete, and cheap enough to run on every proposed change. A transform that has not passed it is not a transform, only a suggestion. Building the check correctly, including what it assumes about don't cares, is what makes an agentic optimization loop safe.

## The miter

The miter joins the two circuits on shared primary inputs, XORs their outputs, and asks a SAT solver whether the XOR can be 1. Any model is a counterexample input; UNSAT means the circuits agree on every input. With several outputs, XOR them pairwise and OR the results.

Suppose $f = d + \overline{a}\,\overline{d} + b\overline{c}d + \overline{b}c$ is rewritten as $g = d + \overline{a}\,\overline{d} + bd + \overline{b}c$. Both reduce to $d + \overline{a} + \overline{b}c$, since $b\overline{c}d$ and $bd$ are both absorbed by $d$, so the miter has no model and the rewrite is accepted. If instead $g$ dropped the term $\overline{a}\,\overline{d}$, the miter would be satisfied by the inputs with $a = 0$, $d = 0$ and not ($b = 0$, $c = 1$), and the solver would return one of them.

## The size of the check

With the Tseitin encoding the miter's CNF is linear in the two circuits: each $k$-input AND, OR, NAND or NOR adds $k + 1$ clauses, each two-input XOR or XNOR adds 4, each inverter 2, then 4 for the miter XOR and 1 for the unit clause asserting it. An original block with one two-input XOR, one three-input AND, one two-input OR and an inverter (13 clauses), compared with a rewrite of 14 clauses, gives a miter of $13 + 14 + 4 + 1 = 32$ clauses.

Linear size is what makes the golden test affordable on every transform in a loop. Runtime is a different matter, since UNSAT miters of restructured logic can be hard; practical checkers first hash both circuits together so identical structure merges, then sweep internal equivalences so the final miter is small.

## Equivalence up to don't cares

Transforms often exploit don't cares: inputs that cannot occur, or outputs nobody observes. Then the right check is equivalence on the care set: the miter is conjoined with a constraint excluding the don't-care inputs, and only mismatches on care inputs count. If $f$ and $g$ differ on five inputs and three of them are certified unreachable, two real counterexamples remain and the transform is rejected.

The don't-care set is part of the claim and must come from a trusted analysis, such as reachability or a verified environment constraint, never from the agent that proposed the transform. A check that accepts any don't-care set the agent supplies can be passed by declaring every mismatch a don't care.

## Going deeper

Kuehlmann and Krohm, "Equivalence checking using cuts and heaps", DAC (1997). Mishchenko, Chatterjee, Brayton and Een, "Improvements to combinational equivalence checking", ICCAD (2006).
