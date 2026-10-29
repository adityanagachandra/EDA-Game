# Unate and binate functions

A function is positive unate in $x$ when raising $x$ from 0 to 1 can never lower the output, negative unate when it can never raise it, and binate when it does both on different rows. Unateness is a monotonicity property per variable. Tautology checkers and two-level minimizers exploit unate covers to stop recursion early, BDD and SAT heuristics split on the most binate variable first, and static timing uses the unateness of a cell arc to decide which edge propagates.

## Comparing the two cofactors

$f$ is positive unate in $x$ when $f_{\overline{x}} \le f_x$, meaning every ON row of the negative cofactor is an ON row of the positive cofactor. It is negative unate when $f_x \le f_{\overline{x}}$, and binate when neither containment holds. If both hold, the cofactors are equal and $f$ does not depend on $x$ at all.

For $f(a, b, c) = ab + \overline{c}$: $f_{\overline{a}} = \overline{c}$ and $f_a = b + \overline{c}$, and $\overline{c} \le b + \overline{c}$, so $f$ is positive unate in $a$. For $c$: $f_c = ab \le 1 = f_{\overline{c}}$, so $f$ is negative unate in $c$. For $a \oplus b$, $f_a = \overline{b}$ and $f_{\overline{a}} = b$ are incomparable, so XOR is binate in both inputs.

## Classifying every variable

Classifying each variable takes $2n$ containment tests of $(n-1)$-variable cofactors, which is cheap for a table and easy on a BDD. A function unate in every variable is called unate; if every variable is positive unate the function is monotone. A NAND is negative unate in every input, so in static timing a rising input can only cause a falling output; an XOR is binate, so both input edges must be propagated to both output edges.

Unate recursive algorithms take the count seriously: a cover that is unate in every variable is a tautology exactly when it contains the all-dash cube, so recursion stops there. Splitting on the most binate variable drives the subproblems toward unateness fastest.

## Syntactic versus semantic unateness

A cover in which each variable appears in only one polarity is unate, but the converse fails: a variable can appear in both polarities in an expression and still be unate as a function. Consider $f = ab + \overline{a}b + ac$. Both $a$ and $\overline{a}$ appear, yet $f = b + ac$, and $f_{\overline{a}} = b \le b + c = f_a$, so $f$ is positive unate in $a$.

Cheap syntactic checks are used as a fast filter; the semantic test on cofactors is the one that decides. Confusing the two gives wrong timing senses and missed recursion shortcuts.

## Going deeper

Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis" (1984), chapter 3, on the unate recursive paradigm.
