# Extracting the cheapest term

An e-graph holds many equivalent implementations; extraction picks one. For tree cost, where each occurrence of a subterm is paid for separately, a bottom-up dynamic program over e-classes finds the optimum even when the e-graph has cycles. Working it by hand shows how operator costs steer strength reduction: whether $a \cdot 2$ becomes a shift, an add or stays a multiply.

## Cost of a class from its e-nodes

Give each operator a cost, here multiply, shift and add, with leaves and constants free. The cost of an e-node is its operator cost plus the costs of its child classes, and the cost of a class is the minimum over its e-nodes. For $C_6 = \{C_0 \cdot C_4, C_0 \ll C_5, C_0 + C_0\}$ representing $a \cdot 2$, with multiply 4, shift 1 and add 1, the three options cost 4, 1 and 1, so $C_6$ costs 1. For $C_7 = \{C_0 \cdot C_2, C_6 + C_0\}$ representing $a \cdot 3$, the options cost 4 and $1 + 1 = 2$, so $C_7$ costs 2. Leaves are computed first, then classes whose children are known.

## The root term

Continue up to the root. With multiply 8, shift 1 and add 3 for $a \cdot 3 + b \cdot 4$: $C_6 = a \cdot 2$ costs 1 (shift), $C_7 = a \cdot 3$ costs $\min(8, 3 + 1) = 4$ via $(a \ll 1) + a$, $C_8 = b \cdot 4$ costs 1 via $b \ll 2$, and the root $C_9 = C_7 + C_8$ costs $3 + 4 + 1 = 8$. The extracted term is $((a \ll 1) + a) + (b \ll 2)$. Change the costs and the answer changes: with an expensive add, $a \cdot 3$ stays a multiply. That sensitivity is why the cost model is the natural extension point for a customer-specific objective such as power.

## Cycles and DAG cost

Rules like $x \cdot 1 = x$ put a class inside itself: $C_0 = \{a, C_0 \cdot C_5\}$. Plain recursion would loop. Iterative relaxation handles it: start every class at infinity, repeatedly update each class from its e-nodes whose children are finite, and stop at a fixpoint. A cyclic e-node never wins, because it costs its class's own cost plus something nonnegative. Tree cost overcounts shared subterms; real hardware computes a shared subterm once, which is DAG cost. DAG-cost extraction is NP-hard and is solved by integer linear programming, by greedy heuristics with sharing awareness, or by limiting the search, and that choice matters for datapath area.

## Going deeper

Willsey et al., "egg: fast and extensible equality saturation", POPL (2021), on extraction; Tate et al., POPL (2009).
