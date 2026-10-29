# Kernels and algebraic division

Which divisors should synthesis try? Brayton and McMullen answered with kernels: the cube-free quotients of an expression by cubes. Kernels are few enough to enumerate, and every multiple-cube common divisor of two expressions shows up as an intersection of their kernels. They are the candidate set for factoring and for extracting shared logic.

## Cube-free quotients and kernels

An expression is **cube-free** if no single cube divides it evenly, that is, no literal appears in all its cubes (and it has at least two cubes). $ab + c$ is cube-free; $ab + ac = a(b + c)$ is not. A **kernel** of $F$ is a cube-free quotient $F/c$, and the cube $c$ is its **co-kernel**. $F$ itself is a kernel (with co-kernel 1) when it is cube-free. For $F = abd + adh + cf + e$: $F$ is cube-free, so it is a kernel, and $F/(ad) = b + h$ is another; no other cube gives a cube-free quotient of two or more cubes, so $F$ has 2 kernels.

## Enumerating kernels and co-kernels

Kernels are found recursively: for each literal appearing in at least two cubes, divide by it, make the quotient cube-free by dividing out the largest common cube, record the kernel, and recurse on it, skipping literals already handled to avoid repeats. Several co-kernels can give the same kernel. For $F = bcdf + bcef + cdfg + cefg = cf(b + g)(d + e)$: dividing by $cf$ gives $bd + be + dg + eg$; by $bcf$ or $cfg$ gives $d + e$; by $cdf$ or $cef$ gives $b + g$. That is 3 distinct kernels and 5 co-kernels. $F$ itself is not a kernel, since $cf$ divides every cube.

## Levels and the intersection theorem

A kernel is **level 0** if it has no kernels other than itself; a level-$n$ kernel has kernels of level at most $n - 1$, with at least one of level $n - 1$. For $F = af + bc + cf + e$ the kernels are $F$, $a + c$ (co-kernel $f$) and $b + f$ (co-kernel $c$); the last two are level 0. Level-0 kernels are cheap to compute and often enough in practice. The key theorem: two expressions $F$ and $G$ have a common divisor of more than one cube if and only if some kernel of $F$ and some kernel of $G$ intersect in more than one cube. So searching kernel intersections finds every multiple-cube common divisor, which is how kernel extraction decides what to share.

## Going deeper

Brayton and McMullen, "The decomposition and factorization of Boolean expressions", ISCAS (1982); Brayton, Hachtel and Sangiovanni-Vincentelli, "Multilevel logic synthesis", Proceedings of the IEEE (1990).
