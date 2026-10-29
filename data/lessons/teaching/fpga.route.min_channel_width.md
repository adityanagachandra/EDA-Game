# Minimum channel width search

The minimum channel width $W_{min}$ is the smallest number of tracks per channel at which a circuit routes on a given architecture. It is the standard measure of routability in FPGA research: architectures and routers are compared by $W_{min}$ and by delay at a relaxed width. Finding it means routing many times, so the search procedure matters.

## Lower bounds from channel density

Before detailed routing, each channel's demand gives a lower bound. In a channel, nets crossing the same column need different tracks, so the channel needs at least its density: the maximum number of nets crossing any one column. With nets spanning columns 4 to 7, 6 to 9, 5 to 10 and 7 to 9, column 7 is crossed by all four, so that channel needs 4 tracks. $W$ must cover the densest channel, so $W_{min}$ is at least the maximum density over all channels. A real router may need more, because switch-box patterns and connection-box flexibility limit which tracks a net can switch between, as with the track domains of disjoint switch boxes.

## Searching for the minimum

VPR finds $W_{min}$ by routing repeatedly. Start at some width; while routing fails, double $W$; once a width succeeds, binary-search between the last failing and the first succeeding width, routing at the midpoint each time, until they differ by 1. If routing succeeds exactly when $W \ge 28$ and the search starts at $W = 2$, it tries 2, 4, 8, 16 and 32 (five attempts), then 24, 28, 26 and 27 (four more), nine routings in all. Each attempt is a complete PathFinder run, and failing attempts are the slowest because they run to the iteration limit. Routability is assumed monotone in $W$, which holds in practice, though the router's heuristics can make it noisy near the threshold.

## Low-stress routing

Routing at exactly $W_{min}$ stresses the router: nets take long detours and the critical path grows. Architecture studies therefore report delay at a relaxed width, commonly $1.3 \, W_{min}$, which models a device with some routing headroom. Unidirectional architectures add a parity rule, because wires come in pairs driven in opposite directions, so $W$ must be even. With $W_{min} = 66$, $1.3 \times 66 = 85.8$ rounds up to 86, which is already even. A paper that reports a lower $W_{min}$ without delay at a low-stress width hides a possible trade: a router can buy routability with delay. Both numbers, together with run time, are needed to judge a change.

## Going deeper

Betz, Rose and Marquardt (1999), chapters 6 and 7; Rose et al., "The VTR Project: Architecture and CAD for FPGAs from Verilog to Routing", FPGA (2012).
