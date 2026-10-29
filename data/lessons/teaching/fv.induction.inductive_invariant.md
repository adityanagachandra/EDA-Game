# Inductive invariants

BMC can only refute; to prove that a property holds forever, formal engines use induction over time. A property is an inductive invariant if it holds initially and every transition from a state satisfying it lands in a state satisfying it. Many true properties are not inductive, and telling the difference between a real bug, a non-inductive property and a wrong reset is the daily work of anyone running unbounded proofs.

## Initiation and consecution

Let $P$ be a set of states, the property. $P$ is an inductive invariant when two checks pass. Initiation: every initial state is in $P$, $I \rightarrow P$. Consecution: every transition from a $P$-state lands in $P$, $P \wedge T \rightarrow P'$, where primes denote next-state copies. If both hold, mathematical induction over the cycle number shows that $P$ holds in every reachable state. Each check is one SAT query, with only one copy of the transition relation, which is why induction scales to designs that BMC cannot exhaust. Take states s0 to s3, initial state s0, transitions s0 → s1, s1 → s0, s2 → s3, s3 → s3. For $P = \{s0, s1\}$, initiation holds and every transition from s0 or s1 stays inside $P$, so $P$ is an inductive invariant.

## Invariant versus inductive

An invariant is a property that holds in every reachable state. Every inductive invariant is an invariant, but not conversely, because consecution quantifies over all $P$-states, reachable or not. In the example, take $P = \{s0, s1, s2\}$. The reachable states s0 and s1 are in $P$, so $P$ is an invariant, but the unreachable $P$-state s2 steps to s3, outside $P$, so consecution fails. The failing pair (s2, s3) is a counterexample to induction (CTI), not a bug. Four outcomes are possible: $P$ is an inductive invariant; $P$ fails initiation (an initial state violates it); $P$ is an invariant but not inductive (only unreachable $P$-states leave $P$); or $P$ is not an invariant (some reachable state violates it, in which case consecution necessarily fails too).

## Reading the four outcomes

Each outcome calls for a different action. Failed initiation means the reset state already violates the property: usually a wrong reset value or a property that ignores the reset period. A reachable violation is a real bug, and BMC will produce the trace. An inductive invariant is a finished proof. The invariant-but-not-inductive case is the one engineers spend time on: the step check fails from a state the design never reaches, so the property must be strengthened with a lemma that excludes that state, or the engine must strengthen it automatically, as k-induction and IC3 do. Distinguishing the last two cases is the key skill: check whether the CTI's start state is reachable, for example with BMC aimed at it. An unreachable start state means no RTL change is needed.

## Going deeper

Bradley and Manna, "The Calculus of Computation" (2007), the chapters on inductive assertions; Sheeran, Singh and Staalmarck, "Checking Safety Properties Using Induction and a SAT-Solver" (FMCAD 2000).
