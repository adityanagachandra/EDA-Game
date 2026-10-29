# Boolean versus algebraic division

Algebraic division is fast but treats a and its complement as unrelated and never uses $x \cdot x = x$. Boolean division can use any identity and any don't care, at much higher cost. Knowing what each can see tells a tool designer when cheap algebraic passes suffice and when a Boolean pass, such as resubstitution with don't cares, is worth its run time.

## What algebraic methods miss

Algebraic division works on SOPs as polynomials, so three classes of rewrite are invisible. Idempotence: $(a + b)(a + c) = a + bc$ needs $a \cdot a = a$. Complements: $a + \bar{a}b = a + b$ needs $a + \bar{a} = 1$. Absorption: $a + ab = a$. Purely algebraic identities, such as $ac + ad + bc + bd = (a + b)(c + d)$, are found. A useful test: if verifying the identity requires treating literals as Boolean rather than as independent polynomial variables, weak division cannot find it. Every algebraic factorization is Boolean-valid, so Boolean methods find a superset of the factorizations.

## Boolean division with don't cares

To divide $f$ by $g$ in a Boolean way, introduce a new variable $y$ standing for $g$, with the don't care $y \oplus g$ (points where $y$ and $g$ disagree never occur). Then minimize $f$ as a function of its inputs and $y$, using that don't care. The minimizer may rewrite $f$ in terms of $y$ wherever that is cheaper. For $f = a + bc$ and an existing node $g = a + b$, algebraic division fails ($f/a = 1$ and $f/b = c$ have an empty intersection), but Boolean minimization finds $f = y(a + c)$, which is correct because $(a + b)(a + c) = a + bc$.

## When to pay for Boolean methods

Algebraic methods run in near-linear time per division and scale to whole networks, so they are used for global restructuring. Boolean methods need don't-care computation and minimization or SAT, so they are applied locally: resubstitution tries a few divisors near each node, inside a window. Some structure only Boolean methods can see, for example multiplexer and XOR forms written in SOP, where $a$ and $\bar{a}$ appear in different cubes. For read-once functions algebraic factoring is already optimal. Modern AIG flows mostly replace algebraic passes with rewriting and truth-table based resubstitution, which are Boolean by construction but local.

## Going deeper

Brayton, Hachtel and Sangiovanni-Vincentelli, "Multilevel logic synthesis", Proceedings of the IEEE (1990); Hachtel and Somenzi (1996), on Boolean division.
