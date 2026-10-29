# PVT robustness and parametric yield

A sized circuit must meet its specs over process corners, supply voltages and temperatures, and under local mismatch. Robust sizing turns those conditions into constraints: the worst corner must pass, and the parametric yield must reach a target. Counting the evaluations this requires, reading worst cases from corner tables and computing yield from a distribution are the routine calculations of robust analog design.

## Corner grids

Process corners describe global variation: TT, FF, SS and the skewed FS and SF, where the first letter is the NMOS speed and the second the PMOS. Supply is swept to its minimum and maximum, temperature from cold to hot, typically -40 to 125 degrees Celsius. A full grid of 5 process corners, 3 supplies and 3 temperatures has $5 \cdot 3 \cdot 3 = 45$ points. Running 100 mismatch Monte Carlo samples at each makes 4500 simulations per candidate. Designers prune the grid to the corners that are worst for each spec (slow and hot for bandwidth, fast and cold for power or stability), but pruning is a judgment that an automated sizer must learn or verify, because the worst corner can change as the sizing changes.

## Worst-case margin

Robust sizing optimizes the worst corner, not the typical one. Given phase margins of 66, 72, 58 and 70 degrees at four corners and a spec of 60, the worst case is 58 and the margin is -2 degrees: the design fails. The constraint for an optimizer is $\min_c \mathrm{PM}_c \ge 60$. A min over corners is non-smooth, since the active corner can switch, which hurts gradient methods; sizers handle it with constraint-per-corner formulations or smooth approximations. A design centred so that its worst-corner margin is positive but small is fragile against model error, so flows add a guard band. The margin to spec at the worst corner is also the quantity an agent or optimizer should report, not the nominal value.

## Parametric yield

Under mismatch and process variation, a performance is a random variable. If Monte Carlo shows DC gain approximately normal with mean 72 dB and standard deviation 1.5 dB, and the spec is gain at least 69 dB, the yield is $\Phi((72 - 69)/1.5) = \Phi(2) = 97.7\%$. Moving the mean to 73.5 dB gives $\Phi(3) = 99.87\%$. Several specs fail jointly, so the true yield is the probability that all pass, at most the smallest single-spec yield; estimating it needs joint samples. Design centring moves the nominal point away from the spec boundaries to maximize this probability. Because Monte Carlo is expensive, yield-aware sizers use sensitivity analysis, response surfaces or importance sampling to estimate it with fewer samples.

## Going deeper

Graeb, "Analog Design Centering and Sizing" (2007).
