# IC3 / property-directed reachability

IC3, also called property-directed reachability (PDR), proves safety properties without unrolling the transition relation. It maintains a sequence of frames, over-approximations of the states reachable in at most $i$ steps, and refines them clause by clause, guided by states that could reach a violation. It is the strongest general-purpose unbounded engine for hardware, implemented as pdr in ABC and used through SymbiYosys's abc pdr engine.

## Frames and their invariants

IC3 keeps frames $F_0, F_1, \dots, F_k$, each a set of clauses over the state variables. $F_0$ is the initial states. The frames satisfy three invariants: $F_i$ over-approximates the states reachable within $i$ steps; the clauses of $F_{i+1}$ are a subset of those of $F_i$, so $F_i \subseteq F_{i+1}$ as sets of states; and $F_i \wedge T \rightarrow F'_{i+1}$. For $i < k$, every frame also implies $P$. IC3 never unrolls $T$: every SAT query involves one frame, one copy of the transition relation and one cube, so queries stay small however deep the analysis goes. The frontier frame $F_k$ is the one under construction, and IC3 asks whether any state in it can step to $\neg P$.

## Blocking proof obligations

If a state $s$ in $F_k$ has a transition to $\neg P$, IC3 creates a proof obligation: show that $s$ is unreachable within $k$ steps by blocking it at frame $k$. To block $s$ at frame $i$, check whether $\neg s$ is inductive relative to $F_{i-1}$. If it is, generalize $\neg s$ to a shorter clause and add it to $F_1$ through $F_i$. If not, the query returns a predecessor $t$ of $s$ inside $F_{i-1}$, and $t$ becomes a new proof obligation at frame $i - 1$, which must be blocked first. Obligations thus form a chain back toward the initial states. If a chain reaches a state in $F_0$, it is a concrete path from an initial state to $\neg P$: a real counterexample, returned as a trace. Otherwise every obligation is eventually blocked, and $F_k$ no longer intersects the predecessors of $\neg P$.

## Propagation and convergence

After the frontier is clean, IC3 opens frame $F_{k+1}$ and propagates: for each clause $c$ of $F_i$, if $F_i \wedge T \rightarrow c'$, push $c$ into $F_{i+1}$. If after propagation some $F_i$ and $F_{i+1}$ have the same clauses, then $F_i$ contains the initial states, implies $P$ and satisfies $F_i \wedge T \rightarrow F'_i$: it is an inductive invariant, and $P$ is proved. On a finite system this must happen eventually, since frames can only gain clauses and there are finitely many. The converged frame is a certificate: anyone can re-check initiation, consecution and $F_i \rightarrow P$ with three SAT queries, without trusting IC3. Because the proof is an invariant rather than an exhaustive unrolling, IC3 can prove properties of designs whose diameter makes BMC hopeless.

## Going deeper

Bradley, "SAT-Based Model Checking without Unrolling" (VMCAI 2011); Een, Mishchenko and Brayton, "Efficient Implementation of Property Directed Reachability" (FMCAD 2011); the ABC documentation for pdr.
