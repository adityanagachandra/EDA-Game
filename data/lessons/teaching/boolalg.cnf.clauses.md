# Conjunctive normal form and clauses

Conjunctive normal form is an AND of clauses, each clause an OR of literals. It is the input language of every SAT solver (the DIMACS format), so equivalence checking, bounded model checking, test generation and many scheduling and routing problems end up as CNF. CNF is good at saying what is forbidden: each clause rules out one cube of assignments. Reading clauses as forbidden cubes makes sizes, model counts and encodings easy to reason about.

## Clauses and assignments

A clause is satisfied when at least one of its literals is true, and a CNF is satisfied when every clause is. A clause is falsified only when all of its literals are false. Under a partial assignment a clause with one unassigned literal and the rest false is a unit clause, which forces that literal; this is unit propagation, the inner loop of CDCL solvers.

Take $(x_1 \lor \neg x_2) \land (x_2 \lor x_3) \land (\neg x_1 \lor \neg x_3)$ under $x_1 = 1$, $x_2 = 0$, $x_3 = 1$. The first clause holds through $x_1$, the second through $x_3$, and in the third both literals are false. One clause is falsified, so the assignment is not a model.

## What one clause forbids

A clause over $k$ distinct variables is false on exactly one assignment of those variables, the one that falsifies every literal. With $n$ variables in all, the other $n - k$ are free, so the clause forbids $2^{n-k}$ of the $2^n$ assignments and allows $2^n - 2^{n-k}$.

The clause $(\neg x_1 \lor x_4 \lor x_5)$ over $x_1, \dots, x_6$ forbids the cube $x_1 = 1, x_4 = 0, x_5 = 0$, which is $2^{6-3} = 8$ assignments. Short clauses forbid more: a unit clause removes half the space, a binary clause a quarter. That is why learned clauses are minimized before being stored, and why clause length is a quality metric for conflict analysis.

## CNF size of familiar functions

The canonical CNF of $f$ has one clause per OFF row, each with all $n$ variables, so $2^n - |ON|$ clauses. Most functions have much smaller CNFs, but some do not. Parity $x_1 \oplus \dots \oplus x_n$ needs $2^{n-1}$ clauses without extra variables: any clause that omits a variable forbids a cube containing assignments of both parities, so every clause must mention all $n$ variables and can exclude only one OFF row.

This is the motivating example for auxiliary variables. With the Tseitin encoding, a chain of two-input XORs encodes parity in $4(n - 1)$ clauses plus one unit clause, linear rather than exponential, at the price of $n - 1$ new variables.

## Going deeper

Biere, Heule, van Maaren and Walsh (eds.), "Handbook of Satisfiability" (2nd ed., 2021), chapters 1 and 2. Knuth, "The Art of Computer Programming", volume 4, fascicle 6, "Satisfiability".
