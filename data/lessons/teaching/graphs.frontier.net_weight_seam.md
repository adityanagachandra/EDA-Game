# A seam for custom net weighting

The most common customer-specific request to a partitioner is to change what a cut net costs: timing-critical nets should cost more, clock nets nothing, buses as a group. The core already computes weighted cuts, so the right design is a narrow seam, a net-weight provider, through which those policies enter. Designing what the seam receives, what it returns and how its results are proved is a concrete exercise in agent-first engine design.

## Weights through the seam, cut in the core

A NetWeightProvider takes each net and its context and returns a non-negative weight. The core uses those weights in its unchanged weighted-cut bookkeeping and gain computations. A provider might return 4 for timing-critical nets, 0 for clock nets (clock-tree synthesis rebuilds them anyway) and 1 otherwise.

With that provider, a partition cutting two critical nets, one clock net and one ordinary net has weighted cut $4 + 4 + 0 + 1 = 9$, whatever the unweighted count of 4 says. The core neither knows nor cares why the weights are what they are; that separation is what lets each customer's policy be a small, testable plug-in.

## What the seam must expose and guarantee

The provider needs enough context for real policies: the net's identity and pins, its role (clock, reset, signal), and timing data such as slack or criticality from the timer, delivered through an observation hook. It should not see the partitioner's internals, such as the gain buckets.

Its contract must state that weights are non-negative (a negative weight would reward cutting a net), deterministic for identical inputs, and stable during a pass. It must also state how weights are represented: integer buckets need integer gains, so real-valued weights must be quantized or the core must use a priority queue. During coarsening, nets with identical pin sets merge by adding weights, so the provider is called once on the fine netlist and the core handles the rest.

## Proving a provider plug-in

A new provider is accepted by differential evidence, not by inspection. First, the neutral setting must reproduce the baseline exactly: with every weight 1 (or every criticality 0), the partition must be identical to the baseline under the same seed. Second, with the policy on, a report on a reference design must show the intended effect: cut critical nets drop, and the cost elsewhere is quantified.

Typical failures are scale problems and nondeterminism. A weight of 100 on macro nets can make those nets uncut while doubling every other cut, which is a Tier 1 retuning of the scale rather than a code bug. Two runs with the same seed that differ point to a provider reading an unordered map or the clock. Composite policies, such as timing weight times a congestion factor, should be stacks of providers combined by a declared rule, each tested separately.

## Going deeper

The OpenROAD documentation for its partitioning tools; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 2 on weighted partitioning.
