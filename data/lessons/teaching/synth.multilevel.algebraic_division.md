# Weak division of SOP expressions

Multi-level synthesis finds structure by dividing expressions. Algebraic (weak) division treats an SOP as a polynomial: literals are variables, a complement is just another variable, and no Boolean identities are used. It is fast and predictable, and it is the engine behind kernels, factoring and extraction.

## Dividing by a cube

To divide an SOP $F$ by a single cube $c$, keep every cube of $F$ that contains all literals of $c$, and remove those literals from it. The result is the quotient $F/c$; the cubes not containing $c$ form the remainder. For $F = abg + ad + ae + af + bd + be + bf$ and $c = b$, the cubes containing $b$ are $abg$, $bd$, $be$, $bf$, so $F/b = ag + d + e + f$ with 4 cubes, and the remainder is $ad + ae + af$. Then $F = b(ag + d + e + f) + ad + ae + af$ exactly. Division by a cube is a filter and a strip, linear in the number of cubes.

## Dividing by a multi-cube divisor

For a divisor $D = d_1 + \dots + d_k$, compute each cube quotient $F/d_i$ and intersect them: $Q = \bigcap_i F/d_i$. The remainder is $R = F - DQ$, the cubes of $F$ not produced by multiplying $D$ and $Q$. For $F = ad + ae + af + bd + cd + ce + cf$ and $D = a + c$: $F/a = d + e + f$, $F/c = d + e + f$, so $Q = d + e + f$, $DQ$ accounts for six cubes, and $R = bd$ with 2 literals. If $F = bef + cef + bg + dh$ and $D = b + c$, then $F/b = ef + g$ and $F/c = ef$, so $Q = ef$, one cube. With sorted cube lists the whole division runs in $O(n \log n)$ time.

## What algebraic means

Weak division is exact as polynomial division: $F = DQ + R$ holds as an identity of polynomials, and $Q$ is the largest expression for which $DQ$ is contained in $F$. It is well defined only when $D$ and $Q$ have disjoint supports, so that multiplying them back creates no $x \cdot x$ or $x \cdot \bar{x}$ terms. The price is blindness to Boolean facts: $(a + b)(a + c) = a + bc$ is invisible, because it needs $a \cdot a = a$, and $\bar{a}$ is unrelated to $a$. The gain is speed and a clean theory of divisors (kernels). Boolean division recovers some of the lost cases using don't cares, at a much higher cost.

## Going deeper

Brayton and McMullen, "The decomposition and factorization of Boolean expressions", ISCAS (1982); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), the chapter on multiple-level optimization.
