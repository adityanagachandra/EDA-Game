# Assume, assert and cover

A formal testbench has three kinds of directive. Assertions state what the design must do and are proved. Assumptions state what the environment does and are trusted. Covers ask whether something can happen and return a witness trace. Getting the balance right is what separates a meaningful proof from one that passes for trivial reasons, so this is the most important methodology concept in property checking.

## Three directives, three jobs

assert property (p) asks the tool to prove that $p$ holds on every trace allowed by the assumptions, or to return a counterexample. assume property (q) restricts the traces the tool considers: inputs are chosen only in ways that satisfy $q$, so assumptions encode the interface contract, such as \"the sender never asserts valid during reset\". cover property (r) asks for one trace on which $r$ happens and returns it as a witness. Every counterexample obeys every assumption, so a counterexample is always a legal environment behavior. The same property can play different roles in different proofs: in assume-guarantee reasoning, block A asserts a property about its outputs and block B, which consumes them, assumes it. The composition is sound only when every assumed property is proved somewhere.

## Over-constraining assumptions

Assumptions are not checked, so a wrong one silently removes behaviors from the proof. The extreme case is a contradiction: if the assumptions cannot all hold on any trace after reset, there are no traces at all, and every assertion passes vacuously. A milder case is common: an engineer adds assume property (!req) to silence a confusing counterexample, and the design now never receives a request, so its assertions pass while saying almost nothing. Even a correct-looking assumption can be false in silicon, for example if the real sender can raise valid during reset; then the proofs do not cover the real behavior. The defenses are to keep assumptions few and traceable to the interface specification, to prove them as assertions on the block that drives those signals, and to check them with covers.

## Covers and vacuity

An implication with an antecedent that can never be true passes vacuously: req |-> ##[1:3] ack is proved if req is stuck at 0. The guard is a cover on the antecedent, cover property (req), which must be reachable. More generally, teams write covers for the normal behaviors a block must exhibit, such as a full FIFO, a back-to-back transfer or a retry; a reached cover returns a witness trace, useful for sanity checks and waveform review, and an unreachable cover flags either dead logic or assumptions that exclude legal behavior. Cover mode in tools such as SymbiYosys searches for these witnesses directly. A proof is trustworthy only when its key covers are reached under the same assumptions, so in practice every significant assertion comes with at least one cover.

## Going deeper

Cerny, Dudani, Havlicek and Korchemny, "SVA: The Power of Assertions in SystemVerilog" (2nd ed., 2015); Seligman, Schubert and Kumar, "Formal Verification: An Essential Toolkit for Modern VLSI Design" (2nd ed., 2023).
