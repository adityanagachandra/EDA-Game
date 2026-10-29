# Boolean networks, BDDs and AIGs

Before logic becomes cells it lives in a technology-independent representation of Boolean functions. Three dominate: the Boolean network (a DAG of nodes, each with a local function), the binary decision diagram (BDD), a canonical graph for a single function under a variable order, and the and-inverter graph (AIG), a DAG of two-input AND nodes with optionally inverted edges. Synthesis, equivalence checking and formal verification all choose among them, and the choice decides what is cheap.

## Truth tables and their size

The plainest representation of an $n$-input function is its truth table, with $2^n$ rows. It is canonical (two functions are equal exactly when their tables are equal) but exponential, so it is used only for small cuts, typically up to 6 to 16 inputs, inside larger algorithms. The number of distinct $n$-input Boolean functions is $2^{2^n}$: 16 for $n = 2$ and 65,536 for $n = 4$. Technology mappers and rewriting engines in ABC enumerate small cuts and compare their truth tables, which is fast precisely because the cuts are small.

## And-inverter graphs

An AIG represents any function with two-input AND nodes and inverters on edges. An $n$-input AND needs $n - 1$ nodes. An $n$-input OR also needs $n - 1$, by De Morgan's law: $a + b = \overline{\bar a \cdot \bar b}$. A two-input XOR needs 3 nodes, and so does a two-to-one multiplexer. A network built from a 4-input AND, a 3-input OR and one XOR without shared logic has $3 + 2 + 3 = 8$ AND nodes. AIGs are not canonical, but they are compact, uniform and fast to rewrite, which is why ABC and Yosys use them as the working form for logic optimization and why the AIGER format exists to exchange them.

## BDDs and canonicity

A reduced ordered BDD fixes a variable order and merges identical subgraphs, so each function has exactly one BDD for that order. Equivalence becomes pointer comparison. Size depends strongly on the order. An $n$-input AND has $n$ internal nodes; $n$-input parity has $2n - 1$ (one node at the top level, two at each level below). Some functions, such as multipliers, have exponentially large BDDs under every order. So BDDs serve where canonicity pays (small blocks, some model checking, don't-care computation), while SAT on AIGs carries most large equivalence and formal checks.

## Going deeper

Bryant, "Graph-Based Algorithms for Boolean Function Manipulation", IEEE Transactions on Computers (1986); Mishchenko, Chatterjee and Brayton, "DAG-Aware AIG Rewriting" (DAC 2006); the ABC documentation.
