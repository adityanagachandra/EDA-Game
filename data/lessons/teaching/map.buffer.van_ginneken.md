# Van Ginneken buffer insertion

Long or heavily loaded nets are too slow even with well-sized gates. Van Ginneken's algorithm finds the buffer placement that maximizes the required time at the driver over a given set of candidate sites, exactly under the Elmore model, by dynamic programming over (load, required time) candidates from the sinks to the source. It is the standard post-mapping buffering engine, and it ties mapping to physical design.

## The dynamic program

Start at each sink with the candidate (sink load, sink required time). Walk toward the driver. Through a wire segment apply $C' = C + C_w$, $q' = q - R(C_w/2 + C)$. At a buffer site, add a buffered copy $(C_b, q - d_b - R_b C)$ of every candidate. At a branch, merge pairs with $C_1 + C_2$ and $\min(q_1, q_2)$. After every step, prune dominated candidates. At the driver, the answer is $\max(q - R_d C)$ over the remaining candidates. Each candidate records its buffer decisions, so the placement is recovered by tracing back.

## A two-pin net by hand

Driver ($R_d$ = 2 kOhm), segment, site 1, segment, sink (6 fF, 512 ps); each segment 0.4 kOhm and 54 fF; buffer 6 fF, 1 kOhm, 18 ps. Through the downstream segment: (60, 512 - 0.4 x 33 = 498.8). At site 1 keep (60, 498.8) and add the buffered (6, 498.8 - 18 - 60 = 420.8). Through the upstream segment: (114, 498.8 - 0.4 x 87 = 464.0) and (60, 420.8 - 0.4 x 33 = 407.6). At the driver: 464.0 - 228 = 236.0 and 407.6 - 120 = 287.6. Buffering wins by 51.6 ps, because it shields the driver from 54 fF of wire.

## Optimality and its limits

The algorithm is exact for its model: given sites, one or more buffer types, Elmore wire delay and linear buffer delay, it maximizes the driver's required time. Its limits are the model's: real buffer delay depends on input slew, which Elmore ignores; buffers cost area and power, which the basic version ignores; and the sites must be chosen in advance, usually every so many microns along a Steiner tree from routing estimates. Extensions add slew constraints, cost-bounded variants that trade buffers for slack, and inverters with polarity tracking. In OpenROAD, repair_design and repair_timing perform buffering in this family after placement.

## Going deeper

van Ginneken, "Buffer placement in distributed RC-tree networks for minimal Elmore delay", ISCAS 1990; Alpert, Mehta and Sapatnekar (eds.), "Handbook of Algorithms for Physical Design Automation" (2008), the chapters on buffer insertion; the OpenROAD documentation for the Resizer.
