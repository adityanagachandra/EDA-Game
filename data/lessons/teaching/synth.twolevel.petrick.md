# Petrick's method

Petrick's method solves the covering problem algebraically. Each column becomes a clause saying which rows can cover it, the product of all clauses says that every column is covered, and multiplying out with absorption lists every irredundant cover. It is exponential, but it is exact, it makes the structure of cyclic tables visible, and it is the Boolean view of what branch and bound searches.

## Writing the product of sums

Give each row (prime) a Boolean variable $P_i$, true if the row is selected. A column covered by rows $P_1$ and $P_3$ yields the clause $(P_1 + P_3)$: at least one must be selected. The covering condition is the product of all clauses, a CNF formula with only positive literals. Before expanding, simplify: duplicate clauses collapse, and a clause that contains another is implied by it and can be dropped, because $(P_1)(P_1 + P_2) = P_1$. For a table with clauses $(P_3 + P_4)$, $(P_2 + P_4)$, $(P_1 + P_4)$, $(P_1 + P_2)$ none contains another, so all four remain.

## Expanding into covers

Multiply the clauses out into a sum of products, applying $X \cdot X = X$ and absorption $X + XY = X$ as you go. Each surviving product is a set of rows that covers every column and from which no row can be removed: an **irredundant cover**. Example: $(P_1 + P_2)(P_1 + P_3)(P_2 + P_3) = (P_1 + P_2 P_3)(P_2 + P_3) = P_1 P_2 + P_1 P_3 + P_2 P_3$. There are three irredundant covers, each of two rows. Expansion can blow up exponentially in the number of clauses, so the method is practical only for small cyclic cores left after essential selection and dominance.

## Choosing the minimum

Among the products, choose the one of least cost: fewest variables for unit cost, or least total literal count when each $P_i$ is weighted by its prime's literals. Petrick's method shows every minimum at once, which reveals ties: a cyclic table such as $\sum m(0, 1, 2, 5, 6, 7)$ over three variables has two different minimum covers of three primes. The expanded SOP is exactly the set of minimal transversals of the column sets, so the method also connects covering to hypergraph transversals and to the SAT view of covering, where a pseudo-Boolean or MaxSAT solver replaces hand expansion on larger cores.

## Going deeper

Petrick, "A direct determination of the irredundant forms of a Boolean function from the set of prime implicants", AFCRC technical report (1956); De Micheli (1994), the chapter on two-level optimization.
