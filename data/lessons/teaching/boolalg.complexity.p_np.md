# P, NP and coNP

Complexity classes tell a tool designer what kind of algorithm to look for. A problem in P has an algorithm whose time is polynomial in the input size. A problem in NP has short, quickly checkable certificates for its yes answers, and coNP for its no answers. Satisfiability, the question "can this circuit output 1?", is in NP; tautology and combinational equivalence, "do these circuits agree on every input?", are in coNP. Most of verification lives on the coNP side, which is why proofs, not just answers, matter there.

## Decision problems and polynomial time

Complexity theory classifies decision problems, yes or no questions, by running time as a function of input length. Is there an input that sets this output to 1? Does this graph have a vertex cover of size at most $k$? Optimization and counting versions are related but distinct.

The input length matters more than it seems. An algorithm that takes $O(2^n)$ steps on a function given as an explicit truth table of $2^n$ bits is linear in its input, so many truth-table problems are in P. The same question asked of a formula or circuit, whose size can be polynomial in $n$, may be NP-hard. 2-SAT, Horn-SAT and checking equality of two ROBDDs in one shared manager (a pointer comparison) are in P; 3-SAT is not known to be.

## Certificates and NP

A problem is in NP when every yes-instance has a certificate of polynomial length that a polynomial-time verifier accepts, and no no-instance has one. For SAT the certificate is a satisfying assignment, checked by evaluating every clause. For vertex cover it is the set of at most $k$ vertices, checked edge by edge.

Membership in NP is an upper bound: trying every certificate solves the problem in exponential time, and nothing better is implied. P is contained in NP, and whether they are equal is open. SAT is also self-reducible: a decision procedure finds a model with about $n$ further calls, fixing one variable at a time and keeping the value under which the formula stays satisfiable.

## coNP and the shape of verification problems

coNP is the class of problems whose no-instances have short certificates. Tautology is the standard example: a single input vector on which the circuit outputs 0 refutes it. Combinational equivalence is coNP-complete for the same reason, since a distinguishing vector is a short certificate of non-equivalence.

This asymmetry shapes tools. When a SAT solver finds the miter of two circuits satisfiable, the model is a counterexample anyone can check in linear time. When it reports UNSAT, the circuits are equivalent, but no short certificate of that is known in general; if every unsatisfiable formula had a polynomial-size proof checkable in polynomial time, NP would equal coNP. Solvers therefore emit proofs such as DRAT, which can be long but are checked by small independent programs.

## Going deeper

Arora and Barak, "Computational Complexity: A Modern Approach" (2009), chapters 1 and 2. Sipser, "Introduction to the Theory of Computation" (3rd ed., 2012), chapter 7.
