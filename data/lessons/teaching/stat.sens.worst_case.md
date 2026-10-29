# Worst-case versus statistical analysis

Worst-case analysis pushes every parameter to its bad limit at once and checks the circuit there. It is simple, deterministic and safe, and it is how corners began. It is also pessimistic in a way that can be quantified exactly, and that quantification is the argument for statistical methods in timing, analog sizing and yield.

## The worst-case stack

For a linear response $D = D_0 + \sum_i a_i\,\Delta p_i$, moving each parameter $k$ sigma in its delay-increasing direction gives
$$\Delta D_{wc} = k\sum_i |a_i|\,\sigma_i.$$
With three parameters contributing $|a_i|\sigma_i = 4, 3$ and 2 ps at $k = 3$, the worst-case delay increase is $3 \times 9 = 27$ ps. The sign of each sensitivity decides the direction each parameter is pushed, which is why a single slow corner can be wrong for a metric that depends on two parameters in opposite directions. The result is a guaranteed bound, provided each parameter really stays within $k$ sigma.

## The statistical spread

For independent Gaussian parameters the same linear response has $\sigma_D = \sqrt{\sum_i a_i^2\sigma_i^2}$, so the $k$-sigma statistical spread is $k\sqrt{\sum_i a_i^2\sigma_i^2}$. In the example, $\sigma_D = \sqrt{16 + 9 + 4} = 5.39$ ps and the 3-sigma spread is 16.2 ps, against 27 ps for the worst case. The two agree only when one parameter dominates or when all parameters are perfectly correlated. For $n$ equal independent contributions the worst case exceeds the statistical value by a factor $\sqrt{n}$, which grows without bound as more sources are modeled.

## How rare the worst case really is

Express the worst-case shift in units of the true output sigma:
$$\beta_{wc} = \frac{k\sum_i |a_i|\sigma_i}{\sqrt{\sum_i a_i^2\sigma_i^2}}.$$
In the example, $\beta_{wc} = 27/5.39 = 5.0$: a design checked at the "3-sigma" worst case is actually guarded to 5 sigma, a probability of $2.9 \times 10^{-7}$ rather than $1.35 \times 10^{-3}$. That margin costs area, power and frequency. The pessimism is not always waste: if parameters are correlated in a way the statistical model omits, or if the tails are heavier than Gaussian, the worst case protects against model error. The engineering choice is between a cheap bound with unknown slack and an accurate estimate that depends on the variation model being right.

## Going deeper

Orshansky, Nassif and Boning, "Design for Manufacturability and Statistical Design" (2008), the chapter on statistical circuit analysis; Director and Hachtel, "The simplicial approximation approach to design centering," IEEE Transactions on Circuits and Systems, 1977.
