# Tseitin encoding of a circuit

The Tseitin encoding turns a circuit into CNF of linear size by giving every gate output its own variable and writing a few clauses that force that variable to equal the gate's function of its inputs. The result is equisatisfiable with the circuit: it has a model exactly when the circuit can produce the asserted output. Every SAT-based equivalence checker, bounded model checker and ATPG engine starts this way, and clause counts per gate are the first numbers a tool designer should know.

## Clauses for one gate

For $y = a \wedge b$, the constraint $y \leftrightarrow (a \wedge b)$ is $(\neg y \lor a)(\neg y \lor b)(y \lor \neg a \lor \neg b)$: three clauses. A $k$-input AND gives $k$ two-literal clauses ($y$ implies each input) and one $(k+1)$-literal clause (all inputs imply $y$), so $k + 1$ clauses. OR is the dual: $(y \lor \neg x_i)$ for each input and $(\neg y \lor x_1 \lor \dots \lor x_k)$. NAND and NOR are the same with $y$ negated, also $k + 1$.

A two-input XOR $y = a \oplus b$ needs four three-literal clauses, one per falsifying row of the three-variable relation: $(\neg y \lor a \lor b)$, $(\neg y \lor \neg a \lor \neg b)$, $(y \lor \neg a \lor b)$ and $(y \lor a \lor \neg b)$. An inverter $y = \neg a$ needs two: $(y \lor a)(\neg y \lor \neg a)$, though most encoders absorb inverters into literals instead.

## Clauses for a netlist

Encode each gate separately and conjoin. To ask whether the output can be 1, add the unit clause $(y_{out})$. For $y_1 = a \wedge b$, $y_2 = y_1 \lor c$, $y_3 = y_2 \oplus d$ with output $y_3$: $3 + 3 + 4 + 1 = 11$ clauses.

The Plaisted-Greenbaum refinement keeps only the half of each gate's clauses needed for the polarity in which the gate is used. That roughly halves the clause count and preserves satisfiability, but it no longer forces internal variables to their true values, so it is unsafe when a tool reads those values back, for example to extract a counterexample for internal signals.

## Variables and the size of the encoding

The encoding has one variable per primary input plus one per gate: the example has 4 inputs and 3 gates, 7 variables. Both counts are linear in the circuit, which is the content of the reduction from circuit-SAT to CNF-SAT and the reason the two are equally hard.

Linear does not mean free. Modern equivalence checkers encode lazily, only the cone of the signals currently being compared, and simplify the AIG first (structural hashing, constant propagation), because the solver's runtime depends heavily on how the structure is presented. The canonical CNF of the same circuit, by contrast, can be exponential.

## Going deeper

Tseitin, "On the complexity of derivation in propositional calculus" (1968). Plaisted and Greenbaum, "A structure-preserving clause form translation", Journal of Symbolic Computation (1986). Biere et al., "Handbook of Satisfiability" (2nd ed.), chapter on CNF encodings.
