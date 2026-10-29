# Implication graphs and decision levels

When a CDCL solver hits a conflict it needs to know why. The implication graph records the answer: one node per assigned literal, edges from the literals that forced an implication, and decision levels that say when each literal was set. Conflict analysis, clause learning and backjumping are all walks over this graph, so reading it fluently is the prerequisite for everything that makes modern solvers fast.

## Decision levels

Each decision opens a new decision level: the first decision is at level 1, the next at level 2, and so on; literals forced before any decision are at level 0. An implied literal sits at the highest level among the literals that forced it, which is the level of the most recent decision when propagation set it. Suppose $x_1 = 1$ is decided at level 1 and $x_4 = 0$ at level 2. If $x_7 = 1$ is implied by $(\neg x_1 \vee x_4 \vee x_7)$, its antecedents are $x_1$ (level 1) and $x_4$ (level 2), so $x_7$ is at level 2. If $x_9 = 0$ is implied by $(\neg x_1 \vee \neg x_9)$, its only antecedent is at level 1, so $x_9$ is at level 1. When implications are listed out of order, compute levels recursively: first the literals whose antecedents are all decisions, then those that depend on them.

## Nodes, edges and reasons

Every assigned literal is a node. A decision has no incoming edges. An implied literal $\ell$ has a reason clause, the clause that became unit and forced it, of the form $(\ell \vee \neg a_1 \vee \dots \vee \neg a_m)$, where each $a_i$ is a literal that was already true. The graph has an edge from each $a_i$ to $\ell$, so the in-degree of $\ell$ equals the size of its reason clause minus one. If $x_5 = 1$ is implied by $(x_5 \vee \neg x_2 \vee x_3 \vee \neg x_8)$, the antecedents are $x_2 = 1$, $x_3 = 0$ and $x_8 = 1$, and the node $x_5$ has three incoming edges. Because propagation only uses literals that are already set, the graph is acyclic, and the trail order of assignment is a topological order of it. Solvers store the graph implicitly: each variable keeps a pointer to its reason clause and its level.

## The conflict and its cone

A conflict adds a special node $\kappa$ with an edge from the negation of each literal of the falsified clause. Walking backwards from $\kappa$ through reason clauses gives the conflict's backward cone. The decisions inside it are a cause of the conflict: setting them the same way always reproduces it. Their negations form a valid learned clause, the all-decision cut. Suppose decisions $x_1 = 1$ (level 1), $x_4 = 0$ (level 2) and $x_6 = 1$ (level 3) are made, and walking back from the conflict reaches $x_1$ and $x_6$ but not $x_4$. The decision clause is $(\neg x_1 \vee \neg x_6)$, two literals, and it already says level 2 was irrelevant. Cutting the graph closer to the conflict, at a unique implication point, usually gives a shorter and more useful clause, which is the next concept.

## Going deeper

Marques-Silva and Sakallah, "GRASP: A Search Algorithm for Propositional Satisfiability" (IEEE Transactions on Computers, 1999); Biere et al. (eds.), "Handbook of Satisfiability" (2nd ed., 2021), the CDCL chapter.
