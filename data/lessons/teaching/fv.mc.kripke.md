# Kripke structures

Model checking needs a mathematical model of the design's behavior, and the standard one is the Kripke structure: states, initial states, a transition relation and a labelling of states with the atomic propositions true in them. A synchronous circuit maps onto one directly, with latch valuations as states. Every temporal-logic semantics, reachability algorithm and state-explosion argument in this subdomain is stated over Kripke structures.

## States, transitions and labels

A Kripke structure is $K = (S, I, R, L)$: a finite set of states $S$, initial states $I \subseteq S$, a transition relation $R \subseteq S \times S$, and a labelling $L$ that maps each state to the set of atomic propositions true there, such as $\{\text{req}, \text{busy}\}$. A path is an infinite sequence of states $s_0 s_1 s_2 \dots$ with $(s_i, s_{i+1}) \in R$, and the behaviors of the system are the paths that start in $I$. The reachable states are those on some such path; properties only constrain reachable behavior. To count them, search forward from the initial states. With states s0 to s4, $I = \{s0\}$ and transitions s0 → s1, s1 → s2, s2 → s1, s3 → s4, s4 → s3, the reachable states are s0, s1 and s2: 3 states, while s3 and s4 are unreachable and irrelevant to any property.

## From circuits to Kripke structures

A synchronous circuit with latches $x_0, \dots, x_{n-1}$ gives a Kripke structure whose states are the $2^n$ latch valuations. The initial states are the reset values, and $(s, t) \in R$ when some input value makes the next-state functions map $s$ to $t$, so free inputs show up as nondeterminism. Take latches x2, x1, x0, reset to 000, with $x_0' = \neg x_0$, $x_1' = x_1 \oplus x_0$ and $x_2' = x_2$. From 000 the latches step to 001, 010, 011 and back to 000, while $x_2$ stays 0: 4 reachable states out of 8. If $x_2'$ were the free input $i$ instead, each step could set $x_2$ either way and all 8 valuations would be reachable. The labelling comes from the signals: a proposition such as busy is true in the states where the busy output is 1.

## Totality and deadlocks

Temporal logics are defined over infinite paths, so model checkers require the transition relation to be total: every state has at least one successor. A state with no successor is a deadlock; paths that reach it end, and properties such as G p would be judged on finite paths with no defined meaning. Circuits are total automatically, since the next-state functions always produce a next state. Abstract models, protocol models and models with constraints are not: an assumption that rules out every input in some state leaves that state without successors. Tools either reject such models or make them total by adding a self-loop to every deadlock state, and they often warn, because an unintended deadlock usually means the environment model is wrong. Counting the deadlock states is counting the states with an empty successor list.

## Going deeper

Clarke, Grumberg, Kroening, Peled and Veith, "Model Checking" (2nd ed., 2018), chapter 2; Baier and Katoen, "Principles of Model Checking" (2008), chapter 2.
