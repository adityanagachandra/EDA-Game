# Constraint satisfaction problems in CNF

A constraint satisfaction problem (CSP) has variables with finite domains and constraints over subsets of them. Register assignment, scheduling, track assignment in routing, graph coloring for multiple patterning and FPGA placement legality are all CSPs. Solving them with a SAT solver means encoding each finite-domain variable into Booleans, and the encoding choice changes clause counts by orders of magnitude and changes how well unit propagation reasons about the problem.

## One-hot encoding of a finite domain

The direct, or one-hot, encoding gives a variable with domain $\{1, \dots, d\}$ one Boolean $v_j$ per value. One at-least-one clause $(v_1 \lor \dots \lor v_d)$ says the variable takes some value, and the pairwise at-most-one clauses $(\neg v_i \lor \neg v_j)$ for $i < j$ say it takes only one. That is $1 + \binom{d}{2}$ clauses: for $d = 4$, $1 + 6 = 7$.

Pairwise at-most-one is quadratic, so for large domains encoders use a sequential counter or a commander encoding, which needs $O(d)$ clauses and $O(d)$ auxiliary variables. The one-hot form propagates well: setting one value true immediately falsifies all the others.

## Graph coloring as a CSP

$k$-coloring a graph $G = (V, E)$ is the classic CSP: each vertex is a variable with domain $\{1, \dots, k\}$, and each edge forbids equal colors. In one-hot CNF each vertex costs $1 + \binom{k}{2}$ clauses, and each edge costs one clause $(\neg x_{u,c} \lor \neg x_{v,c})$ per color. The total is $|V|(1 + \binom{k}{2}) + |E| k$.

A triangle with 3 colors needs $3 \cdot (1 + 3) + 3 \cdot 3 = 21$ clauses and 9 variables. Coloring instances are symmetric (any permutation of colors maps solutions to solutions), so practical encodings add symmetry-breaking clauses, for example fixing the colors of a clique, without which the solver wastes time on equivalent branches.

## Binary encoding of a domain

The binary, or log, encoding represents a value by $b = \lceil \log_2 d \rceil$ bits. It uses far fewer variables, but if $d$ is not a power of two, the $2^b - d$ unused codes must be forbidden, one clause each over the $b$ bits. A domain of 5 values needs 3 bits, and 3 exclusion clauses.

The trade-off is propagation strength. In one-hot form, learning that a value is impossible is a single literal; in binary form it is a clause over all bits, and constraints between variables become long clauses that propagate weakly. That is why hardware-oriented tools often prefer one-hot or order encodings for scheduling and placement, and binary encodings when variable count is the bottleneck.

## Going deeper

Walsh, "SAT v CSP", Principles and Practice of Constraint Programming (2000). Sinz, "Towards an optimal CNF encoding of Boolean cardinality constraints", CP (2005). Rossi, van Beek and Walsh (eds.), "Handbook of Constraint Programming" (2006).
