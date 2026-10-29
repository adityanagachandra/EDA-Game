# FPGA timing analysis

FPGA timing analysis runs ordinary static timing on a graph whose delays come from the device: LUT and flip-flop delays from the timing database, routing delays from the routed path or, before routing, from a placement delay model indexed by distance. Placement, routing and the final report all use the same propagation; what changes is how good the routing delays are.

## Path arrival from a delay model

Before routing, VPR estimates the delay of a connection from a table indexed by $(|\Delta x|, |\Delta y|)$, built by routing sample connections through an empty fabric. A register-to-register path's arrival at the capture flip-flop is the clock-to-q delay, plus every LUT delay, plus the table delay of every connection. For a path launch FF at $(5, 3)$, LUT1 at $(7, 5)$, LUT2 at $(4, 5)$, capture FF at $(3, 5)$, the connections have offsets $(2, 2)$, $(3, 0)$ and $(1, 0)$. If those table entries are 290, 220 and 120 ps, with clock-to-q 50 ps and LUTs of 170 ps, the arrival is $50 + 290 + 170 + 220 + 170 + 120 = 1020$ ps. The table ignores congestion, so it is optimistic.

## Setup slack

With an ideal clock of period $T$, setup slack is $T - t_{su} - \mathrm{arrival}$. A path with arrival 1180 ps, setup 100 ps and $T = 1470$ ps has slack $1470 - 100 - 1180 = 190$ ps. A negative slack means the path fails at this clock. In a full analysis the clock arrives with skew at each flip-flop, and its network has its own routed delay; FPGA clock networks are dedicated, balanced trees, so skew is small but not zero, and vendors' timing reports include clock-path delays and uncertainty. The same arrival-required arithmetic from static timing analysis applies unchanged; only the source of each delay is FPGA-specific.

## Maximum frequency

If a path is critical, the shortest legal period is its arrival plus setup, and the maximum frequency is the reciprocal: $F_{max} = 10^6 / (\mathrm{arrival} + t_{su})$ MHz with delays in ps. For arrival 1310 ps and setup 60 ps, $F_{max} = 10^6 / 1370 = 729.9$ MHz. Because placement-stage delays omit congestion and detours, the post-route $F_{max}$ is usually lower; the gap is why placers recompute criticalities from estimates and routers recompute them from real routes. Open flows report $F_{max}$ from nextpnr's or VPR's own timer, vendor flows from their signoff timers, and the two use different delay data, so cross-tool comparisons need one timer applied to both results.

## Going deeper

Marquardt, Betz and Rose, FPGA (2000), for the placement delay model; Murray et al., "Tatum: Parallel Timing Analysis for Faster Design Cycles and Improved Optimization", FPT (2018).
