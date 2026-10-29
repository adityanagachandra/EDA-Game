# Resolution and variable elimination

Resolution is the single inference rule behind clause-based reasoning: from $(x \lor A)$ and $(\neg x \lor B)$ infer $(A \lor B)$. It is refutation-complete, so a CNF is unsatisfiable exactly when resolution can derive the empty clause. CDCL solvers perform resolution implicitly during conflict analysis, their UNSAT proofs are resolution proofs in compressed form, and preprocessors eliminate variables by resolving out all of their clauses. Counting resolvents is how a tool decides whether an elimination shrinks or grows the formula.

## The resolution rule

Two clauses that clash on exactly one variable, one containing $x$ and the other $\neg x$, resolve on $x$ to the union of their remaining literals, with repeats merged. Resolving $(x_1 \lor x_2)$ with $(\neg x_1 \lor x_2 \lor x_3)$ gives $(x_2 \lor x_3)$: two literals, since $x_2$ appears in both. The resolvent is implied by the two parents, so adding it never changes the set of models.

Resolution can shorten clauses, which is how it eventually derives the empty clause: from $(x)$ and $(\neg x)$ it derives $()$, a contradiction. Conflict analysis in a CDCL solver resolves the conflicting clause with the reasons of the implied literals on the trail until one literal of the current level remains, producing the learned clause.

## Tautological resolvents

If two clauses clash on two variables, every resolvent contains some literal and its negation, a tautology that is always true and useless. Resolving $(x_1 \lor x_2)$ with $(\neg x_1 \lor \neg x_2)$ on $x_1$ gives $(x_2 \lor \neg x_2)$. Such resolvents are discarded.

When all clauses containing $x$ (the set $P$) are resolved against all clauses containing $\neg x$ (the set $N$), there are $|P| \cdot |N|$ pairs, and only the non-tautological, distinct resolvents matter. With $P = \{(x \lor a), (x \lor b)\}$ and $N = \{(\neg x \lor c), (\neg x \lor \neg a)\}$, the four pairs give $(a \lor c)$, $(a \lor \neg a)$, $(b \lor c)$ and $(b \lor \neg a)$; one is tautological, so 3 useful resolvents remain.

## Eliminating a variable

Davis-Putnam elimination removes a variable $x$ completely: delete every clause in $P$ and $N$ and add every useful resolvent. The new formula is satisfiable exactly when the old one was, because $\exists x\, F$ is equivalent to the set of resolvents on $x$ together with the clauses not mentioning $x$. In the example, 4 clauses are replaced by 3, a net gain of one.

Repeated elimination decides satisfiability on its own, but the clause count can explode, since $|P| \cdot |N|$ grows multiplicatively. Bounded variable elimination, introduced in SatELite and used in modern solvers, eliminates $x$ only when the number of added clauses does not exceed the number removed, often combined with subsumption to delete clauses implied by shorter ones.

## Going deeper

Davis and Putnam, "A computing procedure for quantification theory", Journal of the ACM (1960). Een and Biere, "Effective preprocessing in SAT through variable and clause elimination", SAT (2005).
