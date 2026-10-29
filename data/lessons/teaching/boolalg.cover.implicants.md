# Cubes, implicants and prime implicants

An implicant of $f$ is a cube that lies entirely inside the ON-set, so it implies $f$. A prime implicant is an implicant that cannot be enlarged. Every minimum sum-of-products expression can be built from primes, which is why two-level minimizers from Quine-McCluskey to Espresso are organized around generating, expanding and selecting primes. Recognizing implicants and primes is also the cube-level view of SAT: a cube that implies a formula is a partial assignment that satisfies it.

## Implicants: cubes inside the ON-set

Cube $c$ is an implicant of $f$ when every minterm of $c$ is an ON row of $f$, written $c \le f$. List the cube's minterms and check each one; a single OFF row disqualifies it.

For $f(a, b, c) = ab + \overline{c}$ with ON-set $m(0, 2, 4, 6, 7)$: the cube $ab$ covers rows 6 and 7, both ON, so it is an implicant. The cube $b$ covers rows 2, 3, 6 and 7, and row 3 is OFF, so it is not. The minterm $ab\overline{c}$ (row 6) is an implicant too, as every ON minterm is.

## Prime implicants: maximal implicants

An implicant is prime when removing any single literal produces a cube that is not an implicant. Each literal removal doubles the cube, adding the rows that differ in that variable, so primality is checked with one test per literal.

$ab$ is prime: dropping $a$ gives $b$, which covers OFF row 3, and dropping $b$ gives $a$, which covers OFF row 5. $ab\overline{c}$ is not prime: dropping $a$ and $b$ leaves $\overline{c}$, still an implicant. The primes of this $f$ are exactly $ab$ and $\overline{c}$. In SAT terms a prime implicant is a minimal partial assignment that already forces the formula true.

## Recognizing primes in a cover

Espresso's EXPAND step takes each cube of a cover and drops literals as long as the cube stays inside ON plus don't cares, turning it into a prime. The order of literal removal matters: different orders can yield different primes from the same start, and the minimizer chooses the one that covers the most other cubes.

To tell whether a given candidate is prime, expand it as far as possible and compare. If $\overline{a}b\overline{c}$ expands to $\overline{c}$, it was not prime. A cube can be prime while not appearing in any minimum cover, and a minimum cover never needs a non-prime cube, because replacing a cube by a prime containing it never increases cost.

## Going deeper

Quine, "The problem of simplifying truth functions" (1952). McCluskey, "Minimization of Boolean functions" (1956). Brayton et al., "Logic Minimization Algorithms for VLSI Synthesis" (1984), chapter on EXPAND.
