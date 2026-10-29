# Sum-of-products and product-of-sums forms

Every Boolean function has two canonical two-level forms: the sum of its minterms and the product of its maxterms. They are unique, easy to build from a truth table and usually enormous, so no tool uses them as working representations. They matter as the reference point from which minimization starts, as the bridge between truth tables and CNF for SAT, and as a reminder that the same function can be cheap in one form and exponential in the other.

## Canonical sum of minterms

The canonical sum of products has one product term per ON row, each with all $n$ literals: row 6 over $a, b, c$ is $ab\overline{c}$. Its size is $|ON|$ terms and $n \cdot |ON|$ literals.

For $f(a, b, c) = ab + \overline{c}$, the ON-set is $m(0, 2, 4, 6, 7)$, so the canonical SOP is $\overline{a}\,\overline{b}\,\overline{c} + \overline{a}b\overline{c} + a\overline{b}\,\overline{c} + ab\overline{c} + abc$: 5 terms and 15 literals, against 2 terms and 3 literals for $ab + \overline{c}$. Minimization is the job of shrinking this list while keeping the function.

## Canonical product of maxterms

The canonical product of sums has one sum term (a maxterm, or clause) per OFF row, the clause that is 0 on exactly that row: row 5, $a = 1, b = 0, c = 1$, gives $(\overline{a} + b + \overline{c})$. Its size is $2^n - |ON|$ clauses and $n(2^n - |ON|)$ literals.

For $f = ab + \overline{c}$ the OFF rows are 1, 3 and 5, so the canonical POS has 3 clauses and 9 literals. A SAT solver reads CNF, which is a product of sums; the canonical form is the CNF you get directly from a truth table, and by De Morgan, the POS of $f$ is the complement of the SOP of $\overline{f}$.

## Evaluating a product of sums

A sum term is 0 exactly on the cube where all its literals are false, so a product of sums is 0 on the union of those cubes and 1 elsewhere. Counting ON rows means counting the complement of that union.

Take $f(a, b, c) = (a + b)(\overline{b} + c)$. The first clause is 0 on the cube $a = 0, b = 0$ (rows 0 and 1); the second on $b = 1, c = 0$ (rows 2 and 6). These cubes are disjoint, so 4 rows are OFF and $f$ is 1 on $8 - 4 = 4$ rows. The asymmetry is the point: SOP makes ON-set questions easy and OFF-set questions hard, CNF the reverse, and converting between them can grow the formula exponentially, as parity shows.

## Going deeper

Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", chapter 4. Knuth, "The Art of Computer Programming", volume 4A, section 7.1.1, on normal forms.
