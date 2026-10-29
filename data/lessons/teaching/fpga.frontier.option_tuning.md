# Agents tuning nextpnr and VPR options

Open FPGA tools expose dozens of options: placer effort, timing trade-off, router cost schedules, seeds. An agent can tune them with the harness loop, running experiments and promoting the best configuration. Because placement is randomized, the agent must treat results as samples, average over seeds and compare against noise, or it will promote lucky runs.

## Results are distributions

A placer seed changes the random moves and so the final placement; on a design near its limit, Fmax can vary by several percent across seeds. An agent evaluating one configuration therefore runs it on several seeds and logs the mean. With six seeds giving 131.5, 131.5, 134, 133.5, 131.5 and 126.5 MHz, the mean is 131.42 MHz. A single seed would conflate the effect of the configuration with the luck of the seed. The same is true of minimum channel width and run time. Logging every seed's result, not just the mean, also lets the agent see variance, which is itself a property of the configuration: a setting with a slightly lower mean and much smaller spread may be the better choice for production.

## Comparing configurations

The agent compares configurations by their means against a baseline run under the same seeds and benchmarks. With a baseline of 180, 171.5 and 176 MHz (mean 175.83) and three configurations averaging 167, 181.83 and 173.5 MHz, the best improves on the baseline by 6.00 MHz. The comparison is per design, on fixed tool versions and fixed constraints, and every run is logged so the experiment is reproducible. In the harness loop, a configuration change is Tier 1: no recompile, so an agent can run many experiments cheaply, in parallel, inside a sandbox with limits on CPU, memory and time. The expensive part is the runs, not the decisions.

## Gating promotion on noise

Before promoting a configuration, the agent requires its improvement to exceed the noise. A simple gate demands that the mean beat the baseline mean by more than the baseline's seed-to-seed range. With a baseline of 113.5, 115.5, 119 and 116.5 MHz (mean 116.13, range 5.5), candidates averaging 114, 118.5, 125.5 and 124.13 MHz, only the last two clear $116.13 + 5.5$. A statistical test on the seed samples is the principled version. Without such a gate, an agent that tries many configurations will eventually find one that looks better by chance, which is the multiple-comparisons trap. The promoted configuration then goes through golden tests on held-out designs before entering the shared library.

## Going deeper

The nextpnr and VPR command-line documentation for seeds and effort options.
