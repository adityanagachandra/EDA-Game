# Subject graph decomposition into NAND2 and INV

Before any matching, the mapper decomposes the optimized network into a subject graph of two-input NANDs and inverters. The rules are mechanical, but the shape they produce decides which library patterns can later be found. A tool designer needs to be able to run the decomposition by hand, count what it produces, and see where double inversions appear and disappear.

## The decomposition rules

Every multi-input operator is first split into two-input operators with a chosen association, for example a AND b AND c as (a AND b) AND c. Then each operator is rewritten: AND(x, y) becomes INV(NAND(x, y)); OR(x, y) becomes NAND(INV x, INV y) by De Morgan; NOT x becomes INV x. Each binary operator yields exactly one NAND2, so a function with $m$ binary operators has $m$ NAND2 nodes. Before cancellation the inverter count is (number of AND) + 2 (number of OR) + (number of NOT). For f = (a AND b) OR c that is 1 + 2 + 0 = 3 raw inverters around 2 NAND2 nodes.

## Cancelling inverter pairs

Rewriting creates back-to-back inverters wherever an inversion meets an inversion: an OR fed by an AND produces INV(INV(NAND(a, b))) on that input, and NOT applied to an AND produces INV(INV(NAND(...))) at the output. Each such pair is deleted, since INV(INV(z)) = z. Continue the example: OR(AND(a, b), c) becomes NAND(INV(INV(NAND(a, b))), INV c), and the pair cancels to NAND(NAND(a, b), INV c). One pair cancels, so the final tree has 2 NAND2 and 3 - 2 = 1 inverter: 3 nodes. In general final inverters = raw inverters minus twice the number of cancelled pairs.

## Why the shape matters

The same function has many decompositions: a 4-input AND can be a balanced tree or a chain, and both have 3 NAND2 and 3 INV nodes after cancellation, yet a structural matcher sees different shapes. A balanced shape lets a balanced NAND4 pattern match; a chain only matches a chain-shaped pattern. Decomposition also decides depth: a balanced tree of $n$ inputs has about $\log_2 n$ NAND levels, a chain has $n - 1$. Mappers therefore balance the AIG before mapping (synth.aig.balancing), store several pattern shapes per cell, or carry several decompositions at once as structural choices. The decomposition is not neutral preprocessing; it is the first optimization decision.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits", section 10.3; Mishchenko, Chatterjee and Brayton, "DAG-aware AIG rewriting", DAC 2006, for the AIG view of the same subject graph.
