# Strengthening an invariant

When a true property is not inductive, the proof needs a stronger statement: the property together with lemmas that rule out the unreachable states from which the step case fails. Finding those lemmas is the creative core of unbounded verification. Whether a human writes them as helper assertions or an engine such as IC3 discovers them as clauses, the method is the same: read the counterexample to induction and exclude its start state.

## Counterexamples to induction

A failing step check returns a counterexample to induction (CTI): a state $s$ satisfying $P$ (or, for k-induction, a path of $P$-states) with a successor that violates $P$. If $s$ is reachable, the CTI extends to a real bug. If it is not, the CTI only shows that $P$ is too weak to be inductive. The first move is to decide which: run BMC aimed at $s$, or inspect it for values the design cannot produce, such as an illegal state encoding or two pointers that disagree with the occupancy counter. CTIs are informative because they are specific: the solver has found an exact state, and the internal register values in it usually point straight at the missing relationship.

## Adding a lemma

A lemma $L$ is a property of the reachable states that excludes the CTI's start state. If $P \wedge L$ is inductive and holds initially, the conjunction is an inductive invariant, so both $P$ and $L$ are proved invariants. Example: counters $x$ and $y$ reset to 0 and both increment every cycle, wrapping at 16, and $P$ is $x = 15 \rightarrow y = 15$. The step check fails from $x = 14, y = 3$, which steps to $x = 15, y = 4$. The lemma $x = y$ holds at reset, is preserved by every step and implies $P$, so $P \wedge (x = y)$ is inductive. Typical lemmas relate redundant state: a state register is one-hot, an occupancy counter equals write pointer minus read pointer, a valid bit implies its data is in range. A lemma must itself hold initially and be proved, not assumed.

## Where lemmas come from

Lemmas come from three places. Designers and verification engineers write helper assertions from their knowledge of the microarchitecture; these are proved as assertions, never added as assumptions, because an unproved assumption would make the main proof unsound. Engines strengthen automatically: k-induction implicitly excludes CTIs whose start states have no long $P$-path leading into them, and IC3 discovers lemmas as clauses, blocking one CTI at a time with a generalized clause that is inductive relative to its frame. Invariant-inference methods propose candidates, from templates, simulation data or a language model, and keep the subset that the prover confirms is inductive together. All three share the same check: whatever the source of a lemma, the prover decides whether it holds.

## Going deeper

Bradley, "SAT-Based Model Checking without Unrolling" (VMCAI 2011); Bradley and Manna, "The Calculus of Computation" (2007).
