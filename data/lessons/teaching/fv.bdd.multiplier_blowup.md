# Why multipliers blow up BDDs

BDDs made combinational equivalence checking routine for control logic and adders, and then hit a wall on multipliers. The wall is not a tuning problem: Bryant proved that some output bits of integer multiplication have exponential BDDs under every variable order. Understanding that result tells a tool designer when to stop trying to reorder and switch to SAT sweeping, word-level methods or algebraic reasoning instead.

## Bryant's lower bound for multiplication

For an $n$-bit by $n$-bit multiplier, Bryant showed in 1991 that the ROBDD of the middle output bit, bit $n - 1$, has size exponential in $n$ under every variable order; the bound he proved grows like $2^{n/8}$. The argument is a communication argument: whatever the order, there is a level at which the variables above and below must exchange many bits of information to determine the middle bit, and each distinct piece of information is a distinct subfunction at that level, so a distinct node. Low output bits are easy (bit 0 is a single AND) and the highest bits are easier than the middle ones; the middle bits collect the most partial products and carries. The practical meaning is blunt: building BDDs for the outputs of a 32-bit or 64-bit multiplier is hopeless, and dynamic reordering cannot help, because no good order exists.

## Adders versus multipliers

Adders show what good structure looks like. Under the interleaved order $a_0 < b_0 < a_1 < b_1 < \dots$, each sum bit's BDD only has to carry the current carry from one bit position to the next, so the size is linear in $n$. Under the separated order with all $a$ bits first, the diagram must remember all of $a$ before it sees $b$, and the size becomes exponential. So adders are easy only with the right order, while multipliers are hard with every order. The difference is locality: in addition, information flows from bit $i$ to bit $i + 1$ through one carry; in multiplication every output bit mixes products $a_i b_j$ from across both operands, so no order keeps the interacting bits close. Comparators and multiplexers behave like adders; dividers and the middle bits of multipliers do not.

## What verification tools do instead

Equivalence of two similar multiplier netlists, such as before and after local synthesis, is proved by SAT sweeping: the netlists share many internal equivalent nodes, which are proved bottom-up so that each SAT call stays small. For architecturally different multipliers (Booth versus array, for example), few internal equivalences exist, and the effective methods are word-level or algebraic: computer-algebra approaches model each gate as a polynomial and rewrite the specification polynomial, such as $\sum 2^i z_i - (\sum 2^i a_i)(\sum 2^j b_j)$, backwards through the gates until it reduces to zero. Earlier, Bryant and Chen's multiplicative binary moment diagrams (*BMDs) represented word-level multiplication compactly by changing the decomposition, not the order. In model checking, datapath multipliers are usually abstracted as uninterpreted functions so that BDDs or SAT only reason about the control logic around them.

## Going deeper

Bryant, "On the Complexity of VLSI Implementations and Graph Representations of Boolean Functions with Application to Integer Multiplication" (IEEE Transactions on Computers, 1991); Bryant and Chen, "Verification of Arithmetic Circuits with Binary Moment Diagrams" (DAC 1995).
