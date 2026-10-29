# Miters for equivalence checking

Equivalence checking asks whether two circuits compute the same outputs for every input, for example RTL against its synthesized netlist. The miter turns that universal question into a single satisfiability question: build a circuit that outputs 1 exactly when the two designs disagree, and ask a SAT solver whether it can ever output 1. Every combinational equivalence checker, from ABC's cec to commercial signoff tools, is organized around this construction.

## Building the miter

Feed the same primary inputs to both circuits $A$ and $B$, XOR each pair of corresponding outputs, and OR the XORs into a single miter output. The miter is 1 on an input vector exactly when some output differs there, so the circuits are equivalent if and only if the miter output is unsatisfiable. A satisfying assignment is a distinguishing vector, the counterexample. For small circuits the miter can be evaluated by truth table. Let $A = x_1 \wedge x_2$ and $B = x_1 \wedge (x_2 \vee x_3)$. They agree whenever $x_1 = 0$, and when $x_1 = 1$ they differ only if $x_2 = 0$ and $x_3 = 1$. So exactly 1 of the 8 input vectors, 101, drives the miter to 1, and that vector is what a checker would report. Counting distinguishing vectors is the same as counting the 1s of $A \oplus B$.

## Sizing the miter glue

With the Tseitin convention from the SAT subdomain, the glue on top of the two circuits costs, for $n \ge 2$ output pairs: $n$ XOR gates at 4 clauses each, one $n$-input OR at $n + 1$ clauses, and one unit clause asserting that the miter output is 1. That is $4n + (n + 1) + 1 = 5n + 2$ clauses and $n + 1$ new variables (one per XOR, one for the OR). For 4 output pairs the glue adds 22 clauses and 5 variables. The glue is tiny next to the circuits themselves, which is why the miter costs essentially nothing to build. What makes a miter hard is not its size but the lack of shared structure between $A$ and $B$, which forces the solver to discover their equivalence from scratch.

## Sharing inputs and structure

Because both circuits read the same inputs, the primary inputs appear once in the miter CNF. The total variable count is therefore the shared inputs, plus the gates of $A$, plus the gates of $B$, plus the glue: for 10 inputs, 50 gates in $A$, 60 in $B$ and 2 outputs, $10 + 50 + 60 + (2 + 1) = 123$ variables. In practice checkers do not hand this raw CNF to the solver. They first merge the two circuits into one and-inverter graph with structural hashing, so identical logic in $A$ and $B$ becomes one node, and then prove internal equivalences bottom-up. When most outputs land on the same node after hashing, the remaining SAT problems are small. A miter per output cone, rather than one global miter, also keeps each query focused and reports exactly which output differs.

## Going deeper

Brand, "Verification of Large Synthesized Designs" (ICCAD 1993); Mishchenko, Chatterjee, Brayton and Een, "Improvements to Combinational Equivalence Checking" (ICCAD 2006); the ABC documentation for cec.
