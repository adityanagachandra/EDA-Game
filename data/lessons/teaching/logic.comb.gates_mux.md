# Gates and multiplexers

A combinational circuit computes a Boolean function of its present inputs, with no memory. Its building blocks are gates (AND, OR, NOT, NAND, NOR, XOR, XNOR) and the multiplexer, which selects one of several inputs. A gate-level netlist is the data structure that synthesis produces, that placement and timing consume, and that simulators and equivalence checkers evaluate billions of times, so evaluating one by hand is the first skill.

## Gate functions

Each gate is a small truth table. AND is 1 only when all inputs are 1; OR is 0 only when all inputs are 0; NAND and NOR are their complements; XOR of several inputs is 1 when an odd number of inputs are 1, and XNOR is its complement.

Evaluate $y = \text{NOR}(\text{XOR}(a, b), \text{AND}(b, c))$ at $a = 1, b = 1, c = 0$: $\text{XOR}(1,1) = 0$, $\text{AND}(1,0) = 0$, $\text{NOR}(0,0) = 1$. So $y = 1$. A netlist is evaluated in topological order: a gate's output is computed once all of its inputs are known, exactly the order a levelized compiled-code simulator uses.

## The multiplexer

A $2^k$-to-1 multiplexer has $2^k$ data inputs $D_0, \dots, D_{2^k-1}$, $k$ select lines and one output. The select lines, read as a binary number $s$, choose $D_s$. A 2-to-1 mux computes $y = \bar{s} D_0 + s D_1$.

An 8-to-1 mux with data inputs $D_7 \dots D_0 = 10110010$ (so $D_7 = 1$ and $D_0 = 0$) and select $s_2 s_1 s_0 = 101$ outputs $D_5 = 1$. When the data inputs are themselves signals, a mux implements a function: with select lines $a, b$ and data inputs $(D_0, D_1, D_2, D_3) = (0, c, \bar{c}, 1)$, the output is $\bar{a}\bar{b} \cdot 0 + \bar{a} b c + a \bar{b} \bar{c} + a b$, which is 1 on exactly four of the eight input combinations.

## Counting the ones of a small network

The number of input combinations on which a circuit outputs 1 (its on-set size, or number of minterms) is a basic fingerprint. Two circuits with different on-set sizes cannot be equivalent, so a quick count is a cheap first filter, and random simulation estimates it statistically for large circuits.

For a three-input network, enumerate all eight rows. For $y = \text{OR}(\text{AND}(a, b), \text{XOR}(b, c))$: $\text{XOR}(b,c)$ is 1 on four rows ($b \ne c$), and $\text{AND}(a,b)$ adds the row $a = b = c = 1$ (where $b = c$, so not already counted). The on-set has 5 minterms. Exhaustive enumeration costs $2^n$ evaluations, which is why it stops being an option beyond a few dozen inputs and tools turn to BDDs and SAT.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), chapter 2; Mano and Ciletti, "Digital Design" (5th ed., 2013), chapters 2 and 4.
