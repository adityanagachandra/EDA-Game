# The DerateProvider seam

Timing derating is the clearest example of where an agent-first timer needs a seam. The propagation mathematics is settled, but every team margins differently: flat OCV, AOCV tables, POCV with LVF sigmas, and custom margin stacks from silicon correlation. Putting all of that behind one narrow interface keeps the core fixed and turns most customer requests into configuration or small plug-ins.

## The question the core asks

The core, graph propagation and the delay-model API, asks one question at one point: what derate applies to this arc in this context? A DerateProvider answers it. The context must carry everything any scheme might key on: the arc and its cell, the cell's Vt class, early or late, rise or fall, logic depth, distance, slew, corner, clock or data, and path group. A flat OCV or AOCV provider returns a multiplicative factor; a POCV provider returns a mean and a sigma adjustment, which the core combines statistically. Skills such as FlatOCV, AOCV and POCV sit behind the seam, and a CustomMargin layer stacks on top. A seam that omits a field, say distance, cannot host a scheme that needs it.

## Computing with a custom provider

A Tier 2 provider expresses a rule no parameter can: for example, add 0.03 to the late derate factor on data arcs through HVT cells when the path's logic depth is below 5, and otherwise fall through to the AOCV table. For a path of depth 4 with an AOCV value of 1.04, a 60 ps HVT clock arc keeps 1.04, a 50 ps HVT data arc gets 1.07, a 20 ps HVT data arc gets 1.07 and a 20 ps SVT data arc keeps 1.04. The late path delay is $60 \times 1.04 + 70 \times 1.07 + 20 \times 1.04 = 158.1$ ps. At depth 7 the rule would not fire, and every arc would use the table value. The common mistakes are applying the rule to clock arcs, ignoring the depth condition, or multiplying factors instead of adding to them.

## Keeping the seam fast and safe

The provider sits next to the timer's hottest loop, so the interface follows strict rules. Calls are coarse, per batch of arcs rather than per arc through a virtual call, so the core's kernel can still be inlined and vectorized. The extension ABI is frozen, so providers built against it survive core releases. Providers receive plain data, hold no hidden global state, and are deterministic and seedable, because reproducible timing and differential tests depend on identical inputs giving identical derates. Scheme-specific logic, such as the AOCV table lookup, lives in its provider rather than in the propagation loop. A new provider is accepted only when a differential timing report on a reference design shows that only the intended arcs moved, by the intended amounts.

## Going deeper

The OpenSTA source on how derates are applied during delay calculation.
