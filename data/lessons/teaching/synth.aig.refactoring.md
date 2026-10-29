# AIG refactoring

Refactoring is rewriting with larger windows. Instead of looking up a stored structure for a 4-input cut, it takes a cut of up to about ten leaves, computes the node's function over it, derives a new factored form from scratch, and keeps the result if it frees more nodes than it adds. It finds restructurings beyond the reach of 4-input rewriting.

## AIG size of a factored form

A factored form is a tree of binary ANDs and ORs over literals. Each binary operator becomes one AND node (an OR by De Morgan, with inversions free), so a factored form with $L$ literals needs $L - 1$ AND nodes. For $F = c(b + h)fag$, with 6 literals, the AIG needs 5 nodes. This is why factored-form literal count is the right area model for AIGs, and why refactoring computes a factored form rather than an SOP. Sharing between the new structure and the existing graph, found by structural hashing, can make the real cost lower than $L - 1$.

## Factoring versus SOP inside the window

Refactoring collapses the cone of a node over a large cut into a truth table, computes an irredundant SOP (for example by the Minato-Morreale ISOP algorithm on truth tables), and factors it algebraically. An SOP with $L$ literals built as a two-level AIG also costs $L - 1$ nodes, so factoring saves exactly the literal difference. For $bfh + bdf$ (6 literals, 5 nodes) the factored form $(h + d)fb$ has 4 literals and 3 nodes, a saving of 2. Both the SOP and its complement are tried, since one is often much smaller, and the output edge is complemented if the complement wins.

## Accepting a refactoring

The gain is MFFC size minus the nodes the new factored structure adds. If a node's MFFC has 5 nodes and the new form $f(g + a + d)(b + c)$ has 6 literals, it costs 5 nodes and the gain is 0; the plain pass rejects it and the zero-cost variant rfz accepts it. Large cuts make refactoring slower than rewriting, so ABC bounds the cut size and the number of nodes tried. Because the factored form is computed algebraically from an SOP, refactoring can miss XOR-rich structure; rewriting and resubstitution complement it, which is why scripts interleave rw, rf and resub.

## Going deeper

Mishchenko, Chatterjee and Brayton, DAC (2006), which describes refactoring with rewriting; Minato, "Fast generation of prime-irredundant covers from binary decision diagrams", IEICE Transactions (1993).
