# Fast and accurate engines in a modeling loop

Every analysis and implementation engine can exist in a fast version and an accurate version: a fast placer and a full one, a congestion estimator and a global router, a timer with lumped wires and a signoff timer with extracted parasitics. An agent runs the fast versions continuously, watches their metrics and variance, and recurses into the accurate engines only where the fast results cannot be trusted. This loop is where most of an agent-first flow's compute is spent.

## Fast and accurate versions of each engine

A fast engine trades fidelity for speed in a specific, known way. A fast placer may stop global placement at a looser overflow and skip detailed placement. A congestion estimator uses probabilistic or pattern-based routing demand instead of routing. A fast timer may use wire-load estimates or Elmore delay on estimated trees instead of extracted parasitics and may skip crosstalk. A fast extractor uses pattern tables instead of a field solver. A surrogate model, trained on accurate results, is another way to build a fast mode. Each is useful only within a regime where its error against the accurate engine has been measured.

## The cost of a modeling loop

Suppose an agent runs $N$ iterations of the fast loop at $t_f$ each and confirms with an accurate run of $t_a$ every $k$ iterations. The loop takes $N t_f + \lfloor N/k \rfloor t_a$. With $t_f = 6$ minutes, $t_a = 90$ minutes and $k = 10$, a budget of 600 minutes fits $N = 40$ iterations: $240 + 4 \times 90 = 600$. A forty-first iteration would need 606 minutes. Running every iteration accurately would allow only 6. The ratio $t_a / t_f$ and the confirmation interval $k$ decide how many experiments fit in a day.

## Watching metrics and variance

The agent decides from the loop's numbers. If fast-loop WNS is improving steadily with low variance across seeds, the current skill is working and the agent continues. If a metric plateaus, the agent tries another skill or customizes one. If the variance of a metric across seeds or iterations grows, or the design enters a regime where the fast engine has not been calibrated (very high congestion, a new layer stack), the fast numbers stop being informative and the agent recurses: it runs the accurate engine on that region or step. The recursion decision is made from data, which is why the fast engine must report its uncertainty as well as its estimate.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), on congestion estimation; the agent-first EDA notes in this curriculum, section on fast and accurate engines.
