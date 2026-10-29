# Floating random-walk capacitance solvers

The floating random-walk method computes capacitance by Monte Carlo instead of meshing. Each walk starts near one conductor and hops through the dielectric until it lands on some conductor; averaging the walks' scores estimates a capacitance. It needs no mesh and no matrix, its memory is tiny, its accuracy is set by the number of walks, and it parallelizes almost perfectly. These properties make it a popular golden-accuracy engine in commercial extractors.

## Walks and scores

The method rests on the mean-value property of Laplace's equation: the potential at a point equals the average potential over a surrounding surface. A walk starts on a Gaussian surface around conductor $i$, jumps to a random point on the largest empty cube or sphere around its current position (sampled with the known Green's function of that shape), and repeats until it lands on a conductor. A walk ending on conductor $j$ contributes a score, a weight derived from Gauss's law at its start, to the estimate of $C_{ij}$; walks ending elsewhere contribute zero to that entry. The capacitance estimate is the sample mean. Seven walks scoring 2.2, 4, 0.9, 1.6, 3.7, 1.2 and 2.8 fF give 16.4/7 = 2.343 fF, a very rough estimate that real solvers refine with many thousands of walks.

## The standard error

Each walk is an independent sample, so the estimate's standard error is $\sigma/\sqrt{N}$, where $\sigma$ is the standard deviation of the per-walk score. With $N = 10000$ walks, mean 18 fF and $\sigma = 36$ fF, the standard error is $36/100 = 0.36$ fF, a relative error of 2%. The central limit theorem makes the estimate approximately normal, so the solver can report a confidence interval and stop when the interval is narrow enough. The per-walk standard deviation is often larger than the mean, because individual scores vary widely; it is the averaging that delivers accuracy. Variance-reduction techniques, such as better weighting or importance sampling, cut $\sigma$ and therefore the walk count for the same accuracy.

## Walks needed for a target

To reach a relative standard error $\varepsilon$, need $\sigma / \sqrt{N} \le \varepsilon \mu$, so $N \ge (\sigma / (\varepsilon \mu))^2$. With $\mu = 4.5$ fF, $\sigma = 11.25$ fF and a 2% target, $N \ge (11.25/0.09)^2 = 15625$ walks. The square is the key cost law: halving the error quadruples the work. That makes accuracy a natural tuning parameter for this engine, exactly the kind of Tier 1 setting an agent or flow can change per net class: loose for bulk signal nets, tight for clocks. Because walks are independent, the work splits across threads or machines with almost no communication, and an incremental change in geometry only requires re-running the walks for nearby conductors.

## Going deeper

Le Coz and Iverson, "A stochastic algorithm for high speed capacitance extraction in integrated circuits", Solid-State Electronics (1992); Rubinstein and Kroese, "Simulation and the Monte Carlo Method" (3rd ed., 2016), on variance reduction.
