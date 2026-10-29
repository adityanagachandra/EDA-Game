# Bayesian optimization for analog sizing

Analog sizing chooses transistor widths, lengths, bias currents and capacitors to meet specifications such as gain, bandwidth and phase margin while minimizing power or area. Each evaluation is a SPICE run, often over corners. Bayesian optimization fits this setting because it needs few runs, but the constraints must be modeled as carefully as the objective.

## Modeling a constraint

A separate GP models each constrained quantity, such as phase margin. At a candidate sizing, if the phase-margin posterior is Gaussian with mean 63 degrees and standard deviation 4 degrees and the specification is at least 60 degrees, the probability of feasibility is $\Phi((63 - 60)/4) = \Phi(0.75) = 0.773$. The probability is low far from feasible data and near the boundary, and close to one deep inside the feasible region. Specifications checked at several process corners need the constraint at each corner, or a model of the worst corner, because a sizing feasible at the typical corner can fail at slow-slow.

## Constrained expected improvement

A standard way to fold constraints into the search multiplies expected improvement in the objective by the probability of feasibility, assuming the models are independent. If the best feasible power is 3.0 mW and a candidate's power posterior has mean 3.2 mW and standard deviation 0.4 mW, $z = -0.5$ and $\mathrm{EI} = -0.2 \cdot 0.309 + 0.4 \cdot 0.352 = 0.079$ mW. With feasibility probability 0.773, the constrained EI is 0.061 mW. A candidate that promises low power but is probably infeasible scores low, which keeps the search near the region the design can actually use.

## The feasible incumbent

The incumbent $f^*$ in constrained EI must be the best observed value among feasible runs. If runs gave 1.35 mW at phase margin 51 degrees and 3.94 mW at 57 degrees against a 55-degree specification, the incumbent is 3.94 mW, not 1.35. Using the infeasible value makes every feasible candidate look like no improvement, and the search stalls. If no run is yet feasible, the loop should first maximize the probability of feasibility. Sizing results from BO are still candidates: the final sizing needs full corner and Monte Carlo verification, and the layout parasitics change the numbers again.

## Going deeper

Gardner et al., "Bayesian Optimization with Inequality Constraints" (ICML 2014); Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), for the specifications being optimized.
