# PPA closure

Closing a design means meeting its performance target within its power and area budgets, all measured by signoff tools: PPA closure. Every fix trades one against another. Upsizing and low-Vt swaps buy speed with power; spacing and shielding buy SI margin with routing area; voltage buys everything with power and reliability. The final weeks of a project are a search over such trades, and the right choice depends on which constraint binds where.

## The trade-offs

Each closure move has a cost profile. Swapping a cell to a lower threshold voltage speeds it but multiplies its leakage, often several times. Upsizing speeds a gate but raises its input capacitance, its area and its switching power, and can slow its driver. Buffering long nets trades area and power for delay. Spacing or shielding a net reduces coupling and crosstalk without adding leakage, but consumes routing tracks. Raising the supply voltage speeds everything but raises dynamic power with the square of voltage and stresses EM and IR budgets. A design with large positive slack everywhere is not optimal: it is over-designed, and closure includes recovering power and area on non-critical paths.

## Choosing a fix

Given a failing endpoint and candidate fixes with known slack gains and costs, the closure choice is the cheapest fix that actually closes. With WNS of -40 ps and candidates of +55 ps for 0.2 uW of leakage, +30 ps for 1.7 uW, +43 ps for 0.7 uW and +49 ps for 0.7 uW, the 30 ps fix fails to close and is out, and the 55 ps upsizing at 0.2 uW is cheapest among the rest. Real optimizers solve this over many endpoints at once, with interacting fixes on shared logic, using sensitivity-guided greedy moves or Lagrangian relaxation, and they re-time incrementally after each batch. The cost function encodes the project's priorities, so it is a parameter, not a constant of the tool.

## Closure with margins

Closure is only real at signoff accuracy. Optimization often runs on fast extraction and timing; paths closed with zero margin in fast mode can reopen at signoff when the accurate extractor shifts their slack. A disciplined flow keeps margin proportional to the fast mode's error on near-critical paths, or recurses to accurate extraction for them. Power recovery spends positive slack, so it must respect a margin too, or it creates new near-critical paths. And every fix changes parasitics, which changes switching power, so power analysis must use post-ECO parasitics. The value of a trade depends on context: 4% more area in a frequency-limiting block is worth it if the die and floorplan absorb it, and worthless if they do not.

## Going deeper

Sapatnekar, "Timing" (2004), chapters on gate sizing and optimization; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter on timing closure.
