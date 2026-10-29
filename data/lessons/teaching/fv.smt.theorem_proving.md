# Theorem proving for hardware

Model checkers and SAT solvers decide questions about one finite design. Theorem provers prove statements about whole families of designs, such as an adder of every width or a processor against its instruction-set specification, at the price of human guidance. Several floating-point units and processor cores have been verified this way, and theorem provers also host the certified checkers that let other tools be trusted, so a tool designer should know when and how they are used.

## Interactive and automated provers

An interactive theorem prover checks proofs that a human structures with definitions, lemmas and tactics. Widely used systems include ACL2, HOL4, Isabelle/HOL, Coq (now Rocq) and Lean, all of which have been applied to hardware. Most follow the LCF approach: theorems are values of an abstract type that only a small trusted kernel of inference rules can create, so bugs in tactics or automation cannot produce false theorems. Automated provers, first-order resolution provers and SMT solvers, are called from inside interactive provers to discharge routine subgoals, and their results are either reconstructed in the kernel or trusted explicitly. The division of labor is the point: the human supplies the proof structure and the key invariants, and the machine checks every step and closes the easy goals.

## Proofs for every width

A model checker verifies one fixed design; an 8-bit adder and a 64-bit adder are different problems. A theorem prover can prove a parameterized statement for every width $n$ by induction. For a ripple-carry adder, the base case shows the 1-bit full adder computes $a + b + c_{in}$ as a 2-bit value; the inductive step shows that if the $n$-bit adder is correct, appending one full adder gives a correct $(n + 1)$-bit adder, because the carry out of the $n$-bit part has weight $2^n$. Processors are verified by refinement: an abstraction function maps pipelined states to architectural states, and the proof shows that every implementation step corresponds to zero or one specification steps. Proofs about compilers and synthesis passes use structural induction over the syntax tree, following the recursion that defines the pass.

## When to use a theorem prover

Theorem proving costs much more human effort than model checking, so it is used where automatic engines cannot reach: parameterized designs, deep arithmetic such as floating-point division and square root, cache-coherence protocols for any number of nodes, and full processor correctness. Most practical projects combine methods: the prover decomposes the problem and proves the general theorem, and SAT, BDD or model-checking engines discharge finite subproblems. Two cautions apply. A proof shows only that the design meets the stated specification, so a wrong specification gives a correct proof of the wrong thing; specifications are reviewed and tested for that reason. And every assumed axiom or unproved lemma joins the trusted base: a proof that relies on an assumed fact about rounding is only conditional until that fact is discharged.

## Going deeper

Kaufmann, Manolios and Moore, "Computer-Aided Reasoning: An Approach" (2000), on ACL2; Harrison, "Handbook of Practical Logic and Automated Reasoning" (2009); Russinoff, "Formal Verification of Floating-Point Hardware Design" (Springer).
