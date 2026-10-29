# LFSR periods and primitive polynomials

A pattern generator is only useful if it cycles through many states before repeating. The period of an LFSR is set entirely by its characteristic polynomial: a primitive polynomial gives the maximal period $2^n - 1$, visiting every nonzero state, while other polynomials give shorter cycles that depend on the seed. Choosing primitive polynomials, and knowing how many exist, is the core design decision for BIST generators.

## Maximal period

An $n$-stage LFSR has $2^n$ states, and the all-zero state maps to itself, so the longest possible cycle visits all $2^n - 1$ nonzero states. That **maximal-length** sequence (an m-sequence) occurs exactly when the characteristic polynomial is **primitive**: irreducible over GF(2), with a root that generates the multiplicative group of GF($2^n$). A 5-stage LFSR with a primitive polynomial repeats after 31 clocks; a 32-stage one after about $4.3 \times 10^9$. m-sequences have near-balanced numbers of 0s and 1s and good run statistics, which is why they approximate random patterns well.

## Period by stepping

If the polynomial is not primitive, the period is the **order** of $x$ modulo the polynomial (for an irreducible one) or a combination over its factors, and it may depend on the seed. The direct method is to step the register from the seed until the seed recurs. With 4 stages and feedback $Q_1 \oplus Q_2 \oplus Q_3 \oplus Q_4$, the polynomial is $x^4 + x^3 + x^2 + x + 1$; from 1000 the states are 1100, 0110, 0011, 0001, 1000, so the period is 5. That polynomial is irreducible but not primitive, since it divides $x^5 - 1$. A maximal period of 15 would need, for example, $x^4 + x + 1$.

## How many primitive polynomials

The primitive elements of GF($2^n$) are the generators of its cyclic multiplicative group of order $2^n - 1$, and there are $\varphi(2^n - 1)$ of them, where $\varphi$ is Euler's totient. Each primitive polynomial of degree $n$ has $n$ of them as roots (the conjugates $\alpha, \alpha^2, \alpha^4, \ldots$), so the number of primitive polynomials is $\varphi(2^n - 1)/n$. For $n = 5$: $\varphi(31)/5 = 30/5 = 6$. For $n = 4$: $\varphi(15)/4 = 8/4 = 2$, namely $x^4 + x + 1$ and $x^4 + x^3 + 1$. Tables of primitive trinomials are preferred in hardware because they need a single XOR gate.

## Going deeper

Lidl and Niederreiter, "Finite Fields" (1997), chapter 3; Golomb, "Shift Register Sequences" (1967); Bushnell and Agrawal (2000), chapter 15.
