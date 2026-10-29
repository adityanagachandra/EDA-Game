# Karnaugh-map minimization

A Karnaugh map lays out the truth table of a function of up to four or five variables so that logically adjacent minterms are geometrically adjacent. Grouping 1-cells into power-of-two rectangles finds implicants by eye. Synthesis tools never use maps, but the map is the cleanest picture of primes, essentials and covering, the concepts every exact and heuristic minimizer implements.

## Adjacency and implicant groups

Rows and columns are labeled in Gray-code order (00, 01, 11, 10), so neighboring cells differ in one variable, and the map wraps around at its edges. A rectangle of $2^k$ cells that are all 1 (or don't care) is an **implicant** with $n - k$ literals: the variables that stay constant across it. For $f(a,b,c,d) = \sum m(0, 4, 8, 12)$ the four cells form the column $c = 0, d = 0$, a group of 4 giving the implicant $\bar{c}\bar{d}$. Bigger groups mean fewer literals. The largest group through a minterm is the prime implicant with fewest literals that contains it, and finding it is the human version of Espresso's expand step.

## Choosing a minimum cover

A minimum sum of products uses as few product terms as possible, and among those as few literals. The procedure on a map: first take every group that is the only way to cover some 1-cell (an essential prime), then cover the remaining 1s with as few additional primes as possible. For $f = \sum m(1, 4, 5, 7, 9, 15)$ over $a, b, c, d$: minterm 9 needs $\bar{b}\bar{c}d$, minterm 4 needs $\bar{a}b\bar{c}$, and minterm 15 needs $bcd$; together they cover all six 1-cells, so the minimum has 3 products. Every product in a minimum cover can be taken prime, since enlarging a group never adds cost.

## Don't cares and the limits of maps

A don't care cell (written d or X) may be included in a group if it helps and ignored otherwise. With $f = \sum m(1, 5, 6, 11) + d(8)$ the don't care at 8 does not join any useful group, so the minimum stays $\bar{a}\bar{c}d + \bar{a}bc\bar{d} + a\bar{b}cd$ with 11 literals; in other functions one don't care can merge two pairs into a quad and save two literals. Maps stop being practical beyond five or six variables, because adjacency is no longer visible in two dimensions. The underlying problem, choosing primes to cover the ON-set, is what Quine-McCluskey and unate covering solve exactly and what Espresso solves heuristically.

## Going deeper

Karnaugh, "The map method for synthesis of combinational logic circuits", AIEE Transactions (1953); Katz and Borriello, "Contemporary Logic Design" (2nd ed.), on two-level simplification.
