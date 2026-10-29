# Generative circuit design

Generative models can propose new designs rather than score given ones: analog topologies with sizings, layout patterns, or netlists with desired properties. Their outputs are only candidates. A generator is worth using when, after every candidate passes through validity checks and simulation, it delivers more accepted designs per unit of compute than simpler search.

## The acceptance funnel

Each stage filters the generator's output with a tool. Of 1000 generated analog circuits, a netlist checker may find 80 percent valid (connected, simulatable), SPICE may show 25 percent of those meet the specification, and 60 percent of those may be novel rather than near-copies of training examples or of each other. The batch yields $1000 \times 0.8 \times 0.25 \times 0.6 = 120$ useful designs. None of these judgments is made by the generator. A model that produces fluent but invalid netlists, or that memorizes training designs, looks productive until the funnel is applied.

## Cost per accepted design

Generation is cheap; verification is not. If generating costs 0.05 CPU-minutes and simulating a valid design 2 CPU-minutes, the batch above costs $1000 \times 0.05 + 800 \times 2 = 1650$ minutes, and the $1000 \times 0.8 \times 0.25 = 200$ spec-meeting designs cost 8.25 minutes each. The simulation term dominates, so a generator that raises the specification pass rate saves far more than one that generates faster. Invalid designs that the checker rejects before simulation cost almost nothing, which makes cheap validity checks the first filter in any generative flow.

## Comparing with simple search

The baseline is simple search, such as random sampling or a genetic algorithm in the same design space, at the same simulation cost. If random samples meet the specification 1 percent of the time at 2 minutes each, the baseline yields $0.01/2 = 0.005$ designs per minute, against $200/1650 = 0.121$ for the generator, a ratio of about 24. A fair comparison also counts the generator's training data and training compute, the baseline's ability to use the same validity checks, and whether accepted designs are diverse. Generative design adds candidates to a flow; specification checking, corners and layout stay with the engines.

## Going deeper

Goodfellow et al., "Generative Adversarial Nets" (NeurIPS 2014); Huang et al., "Machine Learning for Electronic Design Automation: A Survey," ACM TODAES (2021), section on analog design.
