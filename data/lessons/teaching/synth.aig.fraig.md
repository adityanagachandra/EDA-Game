# FRAIGs and SAT sweeping

Structural hashing merges identical structures, but different structures can compute the same function. A functionally reduced AIG (FRAIG) has no two nodes with equal or complementary functions. It is built by SAT sweeping: random simulation proposes candidate equivalences, a SAT solver proves or refutes them, and counterexamples refine the candidates. The same machinery is the core of combinational equivalence checking.

## Candidates from simulation

Simulate the AIG on a set of random input patterns, packed as machine words, so each node gets a signature: its output bit on each pattern. Nodes with different signatures are certainly different functions. Nodes with equal signatures, or bitwise complementary ones, are **candidate equivalences**. To group complements together, normalize each signature so that a chosen bit is 0, complementing it if needed. With signatures such as 11111010, 00111100, 01100010, 11000011 and repeats of each, nodes fall into a few candidate classes. Simulation is cheap: 64 patterns per word operation per node.

## Proving with SAT

For each candidate class, choose a representative (usually the node with the lowest level) and prove each other member equivalent to it. A proof is a SAT query: is there an input on which the two nodes differ? Encode both cones in CNF (three clauses per AND node, by the Tseitin encoding) plus an XOR of the two outputs, and ask the solver. UNSAT means equivalent: merge the node into the representative, complementing the edge if needed. A class of $s$ members needs $s - 1$ calls, so the number of calls is the total size of the classes minus their number. Merged nodes simplify later queries, which is why sweeping goes in topological order.

## Counterexamples refine the classes

A SAT result means the two nodes differ on the returned input. That input is a valuable pattern: simulate it, append the new bit to every node's signature, and split all classes accordingly. One counterexample often separates many false candidates at once, much more effectively than random patterns, because it exercises a corner that random simulation missed. The loop of simulation, SAT, and refinement continues until every class is proved or split. ABC's fraig and &fraig commands implement it, and the result is used for verification (equivalence checking) and synthesis (removing redundant logic before rewriting).

## Going deeper

Mishchenko, Chatterjee, Jiang and Brayton, "FRAIGs: a unifying representation for logic synthesis and verification", UC Berkeley ERL technical report (2005); Kuehlmann, Paruthi, Krohm and Ganai, IEEE TCAD (2002).
