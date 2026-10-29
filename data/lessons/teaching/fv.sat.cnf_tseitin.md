# The Tseitin encoding

SAT solvers read conjunctive normal form, but designs arrive as circuits. The Tseitin encoding bridges the two: it gives every gate output a fresh variable and writes a few clauses that force that variable to equal the gate's function. The result is linear in the circuit size and equisatisfiable with it. Every SAT-based engine in this course, from equivalence checking to bounded model checking, starts by building this CNF, so its size is the first number a tool designer estimates.

## Clauses for one gate

For $y = a \wedge b$ the encoding writes $y \leftrightarrow (a \wedge b)$ as three clauses: $(\neg y \vee a)$, $(\neg y \vee b)$ and $(y \vee \neg a \vee \neg b)$. The two binary clauses say that $y$ true forces each input true; the long clause says that both inputs true force $y$ true. A $k$-input AND therefore costs $k$ binary clauses plus one clause of $k + 1$ literals, so $k + 1$ clauses. OR is the dual: $(y \vee \neg a_i)$ for each input and $(\neg y \vee a_1 \vee \dots \vee a_k)$, again $k + 1$ clauses. NAND and NOR only flip the polarity of $y$. A 2-input XOR $y = a \oplus b$ needs one 3-literal clause per row of its truth table: $(\neg y \vee a \vee b)$, $(\neg y \vee \neg a \vee \neg b)$, $(y \vee \neg a \vee b)$, $(y \vee a \vee \neg b)$, so 4 clauses. A NOT gate gets no variable at all: its readers use the negated literal of its input.

## Counting a whole circuit

Under this convention a circuit's CNF has one variable per primary input and one per non-inverter gate, and its clauses are the sum of the per-gate costs plus one unit clause asserting the output (or the miter, or the bad signal). Take 4 inputs, two 2-input AND gates, one 3-input OR gate, one XOR gate and one NOT gate. Variables: $4 + 4 = 8$, since the NOT adds none. Clauses: $3 + 3 + 4 + 4 + 1 = 15$. The count is linear in the number of gates, which is the whole point: a netlist with a million gates becomes a CNF with a few million clauses, well within reach of a modern solver. The encoding is equisatisfiable, not equivalent: every satisfying assignment of the circuit extends uniquely to the gate variables, and every model of the CNF restricted to the inputs is a satisfying input vector.

## Why new variables pay off

Without new variables, a CNF over the inputs alone can be exponentially large. The parity $x_1 \oplus \dots \oplus x_n = 1$ is false on the $2^{n-1}$ even-parity rows, and any clause over these variables that is implied by parity must mention all $n$ of them, so it excludes exactly one row. The smallest direct CNF therefore has $2^{n-1}$ clauses. A Tseitin chain of $n - 1$ two-input XORs costs $4(n - 1)$ clauses plus one unit clause, $4n - 3$ in total. For $n = 8$ that is 29 clauses against 128, and the gap doubles with every input. Production encoders go further: the Plaisted-Greenbaum variant keeps only the clauses needed for the polarity in which each gate is used, and AIG-based tools such as ABC encode only AND nodes, with inverters on edges.

## Going deeper

Biere, Heule, van Maaren and Walsh (eds.), "Handbook of Satisfiability" (2nd ed., 2021), the chapter on CNF encodings; Kroening and Strichman, "Decision Procedures" (2nd ed., 2016), chapter 1.
