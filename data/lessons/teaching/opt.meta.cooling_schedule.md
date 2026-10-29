# Cooling schedules and equilibrium

The cooling schedule decides how the temperature falls and how many moves run at each temperature. It sets both the run time and the quality of simulated annealing: cool too fast and the search freezes in a poor local minimum; too slowly and the run wastes time. Practical schedules are geometric, start from a temperature calibrated to an initial acceptance ratio, and adapt to the observed acceptance rate.

## Geometric cooling

The most common schedule multiplies the temperature by a constant $\alpha < 1$ after each batch of moves:
$$T_{k+1} = \alpha T_k, \qquad T_k = T_0 \alpha^k.$$
With $T_0 = 590$ and $\alpha = 0.9$: $T_4 = 590 \cdot 0.6561 = 387.1$. With $\alpha = 0.95$ the temperature takes about 14 steps to halve; with $\alpha = 0.8$ about 3.

Geometric cooling spends equal numbers of steps per factor of temperature, which suits the fact that cost landscapes have structure at many scales: big moves matter at high T and fine adjustments at low T. Theory guarantees convergence to a global optimum only for logarithmic schedules, $T_k \propto 1/\ln k$, far too slow for practice, so production annealers accept a good local optimum from a geometric or adaptive schedule.

## How many temperature steps

To cool from $T_0$ to a final temperature $T_f$ the number of steps is the smallest k with $T_0\alpha^k \le T_f$:
$$k = \left\lceil \frac{\ln(T_f/T_0)}{\ln \alpha} \right\rceil.$$
From $T_0 = 500$ to $T_f = 0.1$ with $\alpha = 0.85$: $\ln(0.0002)/\ln(0.85) = -8.517/-0.1625 = 52.4$, so 53 steps. The total move budget is steps times moves per temperature, often set proportional to the number of cells (VPR uses a multiple of $N^{4/3}$ moves per temperature for N blocks).

The stopping temperature is usually defined relative to the cost scale, for example when T falls below a small fraction of the average cost per net, or when the acceptance ratio stays near zero for several steps. Halving the final temperature adds only $\ln 2/|\ln\alpha|$ steps, about 4 for $\alpha = 0.85$, so the run time depends far more on $\alpha$ and on the moves per step.

## Initial temperature and adaptation

A good start accepts most uphill moves. Sample random moves, average the uphill cost changes $\overline{\Delta^+}$, and choose $T_0$ so the initial acceptance ratio for uphill moves is about $\chi_0$:
$$T_0 = -\frac{\overline{\Delta^+}}{\ln \chi_0}.$$
With $\overline{\Delta^+} = 22$ um and $\chi_0 = 0.9$: $T_0 = 22/0.1054 = 208.8$.

Adaptive schedules then pick $\alpha$ from the measured acceptance ratio: cool quickly while nearly everything is accepted or nearly nothing is, and slowly in the middle, where most of the cost reduction happens. Restarts reheat to a moderate temperature when the search stagnates. For an engine designer the schedule is a parameter set (initial ratio, $\alpha$ table, moves per temperature, stop rule) that should be exposed as configuration, since its best values depend on the design and the cost function.

## Going deeper

Aarts and Korst, "Simulated Annealing and Boltzmann Machines" (1989); Betz, Rose and Marquardt, "Architecture and CAD for Deep-Submicron FPGAs" (1999), for VPR's adaptive annealing schedule.
