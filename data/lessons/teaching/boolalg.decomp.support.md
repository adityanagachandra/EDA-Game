# Support and functional dependence

The support of a function is the set of variables it really depends on. Expressions routinely mention variables the function ignores, and netlists routinely compute signals that are functions of other signals. Finding true support shrinks every downstream problem: a BDD over fewer variables, a cut with fewer leaves, a SAT instance with fewer free inputs. Functional dependence, whether a target can be written as a function of given signals, is the question behind resubstitution in synthesis and behind reparameterization in verification.

## The true support of a function

Variable $x$ is in the support of $f$ exactly when $\partial f / \partial x = f_x \oplus f_{\overline{x}}$ is not identically 0, that is, when the two cofactors differ somewhere. A variable outside the support is called vacuous or redundant.

From a truth table, compare the two halves selected by each variable. For $f(a, b, c, d) = \sum m(4, 5, 6, 7, 12, 13, 14, 15)$ the ON rows are exactly those with $b = 1$: swapping $a$, $c$ or $d$ never changes the output, so the support is $\{b\}$ and $f = b$. With a bit-parallel table this test is one mask, one shift and one compare per variable.

## Vacuous variables in a redundant expression

Expressions hide vacuous variables. In $f = ab + a\overline{b} + c$ the first two terms merge to $a$, so $f = a + c$ and $b$ is vacuous, although it appears twice. Here $f_b = a + c = f_{\overline{b}}$, which confirms it.

Deciding whether a variable is vacuous in an arbitrary formula is as hard as an equivalence check (the cofactors must be proven equal), so tools do it semantically on BDDs or with SAT on $f_x \oplus f_{\overline{x}}$, not by reading the expression. Reporting a smaller support is also a correctness issue: a cut computed over a falsely large support simply wastes inputs, but one computed over a falsely small support produces a wrong circuit.

## When one function is a function of others

Given signals $f_1, \dots, f_k$ and a target $g$ over the same inputs, $g = h(f_1, \dots, f_k)$ for some $h$ exactly when no two input rows share the same values of $(f_1, \dots, f_k)$ but disagree on $g$. Group the rows by the value tuple; each group must be constant in $g$, and that constant is $h$ at the tuple.

Take $f_1 = a$, $f_2 = b$ and $g = a \oplus b \oplus c$ over $a, b, c$. Every pair $(f_1, f_2)$ occurs with $c = 0$ and $c = 1$, so all four groups contain both values of $g$, and $g$ is not a function of $f_1, f_2$. With $g = a \oplus b$ there are no conflicts and $h = f_1 \oplus f_2$. SAT-based resubstitution asks exactly this question with two copies of the circuit and a constraint that their $f$ values agree.

## Going deeper

Mishchenko and Brayton, "SAT-based complete don't-care computation for network optimization" (DATE 2005). Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms", on Boolean difference and support.
