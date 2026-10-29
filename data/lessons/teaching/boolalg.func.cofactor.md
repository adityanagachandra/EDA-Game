# Cofactors of a Boolean function

The cofactor of $f$ with respect to a literal is the function left over when that variable is fixed. It is the most used operation in Boolean reasoning: BDD packages build every node from two cofactors, SAT solvers implicitly work in the cofactor selected by their current assignment, tautology checkers recurse on cofactors, and synthesis computes don't cares and decompositions from them. A tool designer needs to compute cofactors from expressions, from tables and from bit-parallel words without hesitation.

## Fixing one variable

The positive cofactor $f_x$ is $f$ with $x = 1$; the negative cofactor $f_{\overline{x}}$ is $f$ with $x = 0$. Both are functions of the remaining $n - 1$ variables. From an expression, substitute the constant and simplify with $x \cdot 1 = x$, $x \cdot 0 = 0$, $x + 1 = 1$ and $x + 0 = x$.

For $f(a, b, c) = ab + \overline{c}$: setting $a = 1$ gives $f_a = b + \overline{c}$, and setting $a = 0$ gives $f_{\overline{a}} = \overline{c}$. Over the two remaining variables $b, c$, $f_a$ is 1 on three of four rows and $f_{\overline{a}}$ on two.

## Cofactors from the truth table

With a table, a cofactor is a selection. Keep the rows where the variable has the fixed value, and drop that variable's column; the surviving rows, in order, are the cofactor's table. For $f = m(0, 2, 4, 6, 7)$ over $a, b, c$, the rows with $a = 1$ are 4, 5, 6 and 7 with values 1, 0, 1, 1. Renumbered over $b, c$ they are rows 0 to 3 of $f_a$, so $f_a = m(0, 2, 3)$: three ON rows, matching the expression above.

Selecting rows by one bit of the index is a fixed bit pattern, which is why cofactoring an explicit table is a couple of mask-and-shift instructions in practice.

## Cofactor with respect to a cube

Fixing several literals at once gives the cofactor with respect to a cube, for example $f_{a\overline{c}}$ with $a = 1$ and $c = 0$. It is the composition of single cofactors in any order, and its table keeps only the rows consistent with the cube. For $f = ab + \overline{c}$, $f_{ac} = b$ (one ON row over $b$) and $f_{a\overline{c}} = 1$.

Cube cofactors are what a SAT solver's partial assignment selects, and they answer containment questions: a cube $c$ is inside $f$ exactly when $f_c$ is the constant 1. Generalized cofactors (constrain and restrict) extend the idea from cubes to arbitrary care sets and are used to simplify BDDs.

## Going deeper

Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", chapter 2. Brayton, Hachtel, McMullen and Sangiovanni-Vincentelli, "Logic Minimization Algorithms for VLSI Synthesis", chapter 2.
