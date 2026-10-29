# Satisfiability and controllability don't cares

The nodes feeding a gate are not independent: they are functions of the same primary inputs, so some combinations of their values never appear. Those impossible local patterns, derived from the satisfiability don't cares of the fanin nodes, are controllability don't cares for the gate. Computing them by simulation and using them in local minimization is one of the most effective ways to simplify a node.

## Satisfiability don't cares

Each internal node $y_i = f_i(x)$ contributes the satisfiability don't care $y_i \oplus f_i(x)$: in the joint space of inputs and node values, every point where $y_i$ disagrees with its function is impossible. The full SDC of a network is the union over all nodes. It is large but cheap to describe, since it is just the node equations. For $y_1 = \overline{ab}$ and $y_2 = \overline{b + c}$ over inputs $a, b, c$, the SDC contains every point with $y_1 \ne \overline{ab}$ or $y_2 \ne \overline{b + c}$. Projecting this set onto the local inputs of a downstream node gives that node's controllability don't cares.

## Impossible local patterns

For a node $g$ whose local inputs are $y_1, \dots, y_k$, a local pattern is **controllability don't care** if no primary input vector produces it. Compute the reachable patterns by simulating every input vector (or by image computation with BDDs or SAT for large windows) and take the complement. With $y_1 = \overline{ab}$ and $y_2 = \overline{b + c}$: $y_2 = 1$ needs $b = 0$, which forces $y_1 = 1$, so the pattern $(y_1, y_2) = (0, 1)$ never occurs. Three of the four patterns appear, so one is a controllability don't care for $g$. With three correlated nodes, three of the eight patterns may be unreachable.

## Minimizing a node with them

Feed the unreachable patterns to a two-level minimizer as don't cares when simplifying $g$'s local function. Suppose $g$ must be 1 on the reachable patterns 2, 6 and 7 of $(y_1 y_2 y_3)$ and 0 on the other reachable ones, while 0, 1 and 3 are unreachable. Without don't cares $g$ needs $y_2\bar{y_3} + y_1 y_2$; with them the minimum is just $g = y_2$, one literal. This is node simplification in the SIS sense, and ABC's mfs does the same on mapped networks with SAT-based don't-care computation. The resulting node is not equivalent to the old one in isolation, only in context.

## Going deeper

Hachtel and Somenzi (1996), on don't-care based multi-level minimization; Mishchenko and Brayton, "SAT-based complete don't-care computation for network optimization", DATE (2005).
