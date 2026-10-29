# Random-pattern testability

Many faults are detected by almost any pattern, and some by almost none. The detection probability of a fault under uniformly random patterns decides how many random patterns a test needs, whether logic BIST is practical, and where ATPG must take over. The computations are small probability exercises, but their consequence, the long tail of random-pattern-resistant faults, shapes test architecture.

## Detection probability

If inputs are independent and uniformly random, a fault's **detection probability** $p$ is the fraction of the $2^n$ input vectors that detect it. On a small circuit it can be computed exactly by exhaustive simulation: if 12 of 32 vectors detect a fault, $p = 0.375$. Structure explains extremes. The output stuck-at-0 of a $k$-input AND fed by primary inputs needs all $k$ inputs at 1, so $p = 2^{-k}$; for $k = 16$ that is about $1.5 \times 10^{-5}$. Such faults are **random-pattern resistant**. Deep logic compounds the effect, because activation and each propagation condition multiply.

## Patterns for a target confidence

With $N$ independent random patterns, the fault escapes with probability $(1 - p)^N$. To make escape at most $\epsilon$,
$$N = \left\lceil \frac{\ln \epsilon}{\ln(1 - p)} \right\rceil \approx \frac{-\ln \epsilon}{p}.$$
For $p = 2^{-5}$ and $\epsilon = 0.05$, $N = 95$. For $p = 2^{-16}$ and the same confidence, $N \approx 196{,}000$. The number of patterns grows as $1/p$, which is why one hard fault can dominate a random test's length.

## Expected coverage

For a set of faults with detection probabilities $p_i$, the expected coverage after $N$ patterns is the mean of $1 - (1 - p_i)^N$. Easy faults saturate within tens of patterns, hard ones climb slowly, so coverage curves rise fast and then flatten. With faults at $2^{-1}$, $2^{-4}$, $2^{-8}$ and $2^{-8}$, after 100 patterns the expected coverage is about 66 percent, because each $2^{-8}$ fault is only about 32 percent likely to be caught. The remedies are deterministic ATPG for the tail, test points that raise $p$, and weighted random patterns that bias input probabilities.

## Going deeper

Bushnell and Agrawal (2000), chapter 15; Bardell, McAnney and Savir, "Built-In Test for VLSI: Pseudorandom Techniques" (1987).
