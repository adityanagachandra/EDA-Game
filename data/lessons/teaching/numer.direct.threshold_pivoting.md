# Threshold partial pivoting

Sparse LU must serve two masters: pivots large enough for stability and pivots placed to minimize fill. Threshold partial pivoting reconciles them by accepting any pivot that is not too small relative to its column and then choosing among those by Markowitz count. Every unsymmetric sparse solver used in circuit simulation has some form of this rule.

## The threshold test

Fix a threshold $0 < u \le 1$. An entry $a_{ij}$ of the active column $j$ is eligible as a pivot if
$$|a_{ij}| \ge u \max_k |a_{kj}|.$$
With $u = 1$ this is partial pivoting: only the largest entry qualifies. With $u = 0.1$ any entry within a factor 10 of the largest qualifies.

Example: column entries $8.0, -3.5, 0.6, -6.2$ with $u = 0.5$. The cutoff is $0.5 \cdot 8.0 = 4.0$, so the entries 8.0 and -6.2 are eligible and -3.5 and 0.6 are not. Lowering $u$ to 0.05 makes all four eligible, giving the ordering heuristic more freedom.

## Choosing among eligible pivots by Markowitz count

Among eligible entries, pick the one with the smallest Markowitz count $(r_i - 1)(c_j - 1)$. Within one column $c_j$ is shared, so the choice reduces to the row with the fewest nonzeros; ties are broken by larger magnitude.

Example: in a column with entries 8.0 (row with 6 nonzeros), -6.2 (row with 3 nonzeros), -3.5 (row with 2 nonzeros), and $u = 0.5$, the eligible entries are 8.0 and -6.2, and -6.2 wins with its shorter row. With $u = 0.25$ the entry -3.5 becomes eligible and wins instead.

Real solvers search across several columns, not one, and many keep the previous pivot sequence during refactorization while it still passes the threshold test, which is very effective in Newton loops where values change slowly.

## Tuning the threshold

Every multiplier is bounded by $1/u$, so one elimination step can grow an entry by at most a factor $1 + 1/u$, and $k$ steps by $(1 + 1/u)^k$ in the worst case. With $u = 0.1$ that is 11 per step and 1331 after three steps; with $u = 0.01$ it is 101 per step. The bound is pessimistic, but it shows the trade: small $u$ buys sparsity with stability risk.

Common defaults are $u = 0.1$ (many general-purpose codes) and $10^{-3}$ (Kundert's Sparse 1.3 used in SPICE3 and ngspice). The largest $u$ for which an entry $a$ stays eligible is $|a|/\max_k|a_{kj}|$. If a refactorization later finds a pivot that fails the test because values changed, a good solver re-pivots that column rather than continuing with an unstable choice.

## Going deeper

Duff, Erisman and Reid, "Direct Methods for Sparse Matrices" (2nd ed., Oxford, 2017), chapters 7 and 10; Kundert, "Sparse matrix techniques", in Ruehli (ed.), "Circuit Analysis, Simulation and Design" (North-Holland, 1986).
