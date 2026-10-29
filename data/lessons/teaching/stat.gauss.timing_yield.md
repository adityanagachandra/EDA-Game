# Timing yield from a delay distribution

A statistical timer's output is not a slack but a distribution, and the number that matters to the business is the fraction of manufactured dies that meet the target frequency. That timing yield, and the clock period that achieves a target yield, are direct Gaussian computations once the circuit-delay distribution is known, and they show why the dependence between critical paths matters.

## Yield at a given period

If the circuit delay $D$ is Gaussian with mean $\mu$ and standard deviation $\sigma$, the probability a die meets clock period $T$ is
$$Y(T) = P(D \le T) = \Phi\!\left(\frac{T - \mu}{\sigma}\right).$$
With $\mu = 800$ ps, $\sigma = 25$ ps and $T = 840$ ps, $z = 1.6$ and $Y = 94.5$ percent. The same die population at $T = 800$ ps yields 50 percent: designing to the mean delay ships half the parts slow. Statistical timing replaces a single worst-case slack with this curve, and a yield-aware optimizer works on the curve rather than on one corner.

## Period for a target yield

Inverting gives the period that achieves yield $Y$: $T = \mu + z_Y\,\sigma$ with $z_Y = \Phi^{-1}(Y)$. For 99 percent, $z = 2.326$; for 99.87 percent, $z = 3$. With the numbers above, a 99 percent timing yield needs $T = 800 + 2.326 \times 25 = 858$ ps. The guard-band $z_Y\sigma$ grows only slowly with the yield target, but it grows without limit, and a design that reports only a mean delay hides it. Parametric binning, selling faster dies at higher frequency, uses the same inversion at several targets.

## Several critical paths

A chip meets timing only if every critical path does, so the circuit delay is the maximum of path delays. If two paths are independent, $Y = Y_1 Y_2$: two paths each at 95 percent give 90.25 percent. If they are perfectly correlated, the chip yield is the smaller of the two. Real paths sit in between, which is why the statistical maximum must carry correlation. A design with hundreds of near-critical paths that are mostly independent loses much more yield than its worst path suggests, which is the statistical form of the wall of slack seen after aggressive optimization: equalizing many paths to the same nominal slack makes the maximum later and wider.

## Going deeper

Blaauw, Chopra, Srivastava and Scheffer, "Statistical timing analysis: from basic principles to state of the art," IEEE TCAD, 2008.
