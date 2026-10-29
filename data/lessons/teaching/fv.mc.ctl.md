# Computation tree logic

Computation tree logic (CTL) reasons about the tree of all futures from a state rather than about one path. Every temporal operator is paired with a path quantifier: E for some path, A for all paths. CTL model checking runs in time linear in the size of the structure and the formula, which made it the logic of the first symbolic model checkers, and it expresses branching properties such as \"a reset state is always reachable\" that LTL cannot.

## Path quantifiers and one step

CTL formulas are true or false at states. EX $\varphi$ holds at $s$ if some successor of $s$ satisfies $\varphi$; AX $\varphi$ holds if every successor does. Because the transition relation is total, every state has at least one successor, so AX $\varphi$ implies EX $\varphi$. Model checkers compute satisfying sets bottom-up: first the set of states satisfying the subformula, then the operator. Take states s0 $\{p\}$, s1 $\{\}$, s2 $\{p\}$, s3 $\{p\}$ with transitions s0 → s1, s2; s1 → s3; s2 → s2; s3 → s0, s1. The $p$-states are s0, s2, s3. EX $p$ holds at s0 (successor s2), s1 (s3), s2 (s2) and s3 (s0): all 4 states. AX $p$ holds at s1 (only s3) and s2 (only s2), but not at s0 (s1 lacks $p$) or s3 (s1 again): 2 states.

## Fixpoint operators EF, AG, EG, AF and EU

EF $\varphi$: some path eventually reaches $\varphi$, the states that can reach a $\varphi$-state (backward reachability). AG $\varphi$: on all paths $\varphi$ always holds, the dual $\neg$EF$\neg\varphi$. EG $\varphi$: some path stays in $\varphi$ forever, the states with an infinite path inside the $\varphi$-states. AF $\varphi$: every path eventually reaches $\varphi$, the dual $\neg$EG$\neg\varphi$. E[$\varphi$ U $\psi$]: some path stays in $\varphi$-states until it reaches a $\psi$-state. In the example above, EF $\neg p$ holds at s0, s1 and s3, which can reach s1, but not at s2, whose only path loops on itself; so AG $p$ holds exactly at s2. EG $p$ holds at s2 (its self-loop), s0 (s0 → s2) and s3 (s3 → s0 → s2): 3 states. With $q$ labelling only s1, E[$p$ U $q$] holds at s1 itself and at s0 and s3, which step to s1 from a $p$-state, while s2 can never leave its loop: 3 states.

## Nesting and the LTL versus CTL divide

Nested formulas are computed inside out. AG EF $p$, \"from every reachable state $p$ can still be reached\", is the reset property: compute EF $p$, then AG of that set. EF AG $p$ says some path reaches a region where $p$ holds forever on all paths. AG($p \rightarrow$ AF $q$) is branching-time response. CTL and LTL are incomparable. AG EF $p$ has no LTL equivalent, because LTL can only speak about the paths that do happen, not about paths that could branch off. F G $p$ has no CTL equivalent: AF AG $p$ is strictly stronger, since it requires every path to reach a point after which all branches keep $p$. CTL model checking is linear in the structure and the formula, while LTL model checking is exponential in the formula, a trade-off that shaped which logic each tool adopted.

## Going deeper

Clarke, Emerson and Sistla, "Automatic Verification of Finite-State Concurrent Systems Using Temporal Logic Specifications" (ACM TOPLAS, 1986); Baier and Katoen, "Principles of Model Checking" (2008), chapter 6.
