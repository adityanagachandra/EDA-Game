# Agent-proposed invariants filtered by the prover

Unbounded proofs often stall for lack of a strengthening invariant, and proposing candidate invariants is exactly the kind of guess a language model or agent can make cheaply. The formal engine then filters the guesses. The Houdini algorithm does this filtering with nothing but induction checks: start from all candidates and repeatedly drop the ones that are not inductive relative to the rest. Whatever survives is sound by construction, so the proposer never needs to be trusted.

## Candidates in, survivors out

An agent proposes candidate invariants $c_1, \dots, c_m$, each a set of states (in practice, an assertion over the design's signals). Houdini looks for the largest subset whose conjunction is inductive. One round works as follows: let $\text{Inv}$ be the conjunction (intersection) of all surviving candidates; remove every candidate that fails initiation (excludes an initial state) or that some state of $\text{Inv}$ can leave in one transition, $\text{Inv} \wedge T \not\rightarrow c'$. Example: states s0 to s3, initial s0, transitions s0 → s1, s1 → s0, s1 → s2, s2 → s0; candidates $c_1 = \{s0, s1, s2, s3\}$, $c_2 = \{s0, s1\}$, $c_3 = \{s0\}$. In round 1, $\text{Inv} = \{s0\}$, whose only successor s1 is outside $c_3$, so $c_3$ is removed, while $c_1$ and $c_2$ survive: 2 candidates.

## Iterating to the fixpoint

Removing a candidate weakens the conjunction, so a state that was excluded before may now be inside $\text{Inv}$, and its successors can break another candidate. So Houdini repeats rounds until one removes nothing. In the example, round 2 has $\text{Inv} = c_1 \cap c_2 = \{s0, s1\}$; s1 now steps to s2, which $c_2$ excludes, so $c_2$ is removed. Round 3 has $\text{Inv} = c_1$, all states, every successor stays inside $c_1$, and nothing is removed. One candidate survives. At the fixpoint, initiation holds and $\text{Inv} \wedge T \rightarrow c'$ for every survivor, so $\text{Inv} \wedge T \rightarrow \text{Inv}'$: the conjunction of the survivors is an inductive invariant, and it is the largest inductive subset of the candidates. If it implies the target property, the property is proved; otherwise the agent proposes more candidates and the loop continues.

## Rounds, cost and trust

Each round costs one SAT query per surviving candidate (or one query with activation literals for all of them), and each round except the last removes at least one candidate, so there are at most $m + 1$ rounds for $m$ candidates; the example used 3, counting the final round that removes nothing. Two properties make this a good fit for agents. Soundness does not depend on the proposer: a wrong candidate is simply removed, so an agent can propose hundreds of guesses, including hallucinated or irrelevant ones, and the prover keeps only what is true together. And failures are informative: a removed candidate comes with a counterexample to induction that the agent can read to propose a better one. The cost to watch is the number of candidates, since every bad guess costs solver time, and a seed set of mutually supporting good candidates matters more than volume.

## Going deeper

Flanagan and Leino, "Houdini, an Annotation Assistant for ESC/Java" (FME 2001); Bradley, "SAT-Based Model Checking without Unrolling" (VMCAI 2011), on inductive strengthening.
