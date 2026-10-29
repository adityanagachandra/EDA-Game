# Preprocessing: variable elimination and subsumption

Before search, and sometimes between restarts, a SAT solver rewrites the formula into an equisatisfiable one that is smaller or easier. Bounded variable elimination and subsumption are the two workhorses. CNF produced by Tseitin encoding of circuits is full of variables that can be eliminated cheaply, so preprocessing can remove a large share of the formula before the search even starts.

## Eliminating a variable by resolution

Let $P$ be the clauses containing $x$ and $N$ the clauses containing $\neg x$. Replacing $P \cup N$ by all resolvents of a clause of $P$ with a clause of $N$ on $x$ gives a formula without $x$ that is satisfiable exactly when the original is (Davis-Putnam elimination). Resolvents that contain some literal and its negation are tautologies and are dropped. Example: $P = \{(x \vee a), (x \vee b \vee \neg c)\}$ and $N = \{(\neg x \vee \neg a), (\neg x \vee c)\}$. The four pairs resolve to $(a \vee \neg a)$, a tautology; $(a \vee c)$; $(b \vee \neg c \vee \neg a)$; and $(b \vee \neg c \vee c)$, a tautology. So elimination produces 2 non-tautological resolvents. A model of the new formula extends to $x$: set $x$ so that every clause of $P$ and $N$ is satisfied, which the resolvents guarantee is possible.

## When elimination pays

Eliminating $x$ can multiply clauses: $|P| \cdot |N|$ resolvents replace $|P| + |N|$ clauses. Bounded variable elimination, as in the SatELite preprocessor, eliminates $x$ only if the count does not grow. In the example, 2 resolvents replace $2 + 2 = 4$ clauses, a net change of $2 - 4 = -2$, so $x$ is eliminated. With $|P| = |N| = 3$ and no tautologies, 9 resolvents would replace 6 clauses, a net change of $+3$, and $x$ would stay. Tseitin variables of gates with few fan-outs often pass the test, which is why circuit-derived CNF shrinks so much. The solver records the eliminated clauses so that, after finding a model of the reduced formula, it can reconstruct values for the eliminated variables in reverse order.

## Subsumption

Clause $C$ subsumes clause $D$ if every literal of $C$ is in $D$. Then $D$ is implied by $C$ and can be deleted without changing the set of models. Among $(a \vee b)$, $(a \vee b \vee \neg c)$, $(\neg a \vee c)$ and $(b)$, the clause $(b)$ subsumes both $(a \vee b)$ and $(a \vee b \vee \neg c)$, and $(a \vee b)$ also subsumes $(a \vee b \vee \neg c)$, so two clauses can go; count each deletable clause once even when several clauses subsume it. Subsumption checks use occurrence lists and signatures (bit masks of the variables) to avoid comparing every pair. A related rule, self-subsuming resolution, shortens a clause: from $(a \vee b)$ and $(\neg a \vee b \vee c)$ the resolvent $(b \vee c)$ subsumes the second clause, which can be replaced by it.

## Going deeper

Een and Biere, "Effective Preprocessing in SAT through Variable and Clause Elimination" (SAT 2005); Biere et al. (eds.), "Handbook of Satisfiability" (2nd ed., 2021), the chapter on preprocessing.
