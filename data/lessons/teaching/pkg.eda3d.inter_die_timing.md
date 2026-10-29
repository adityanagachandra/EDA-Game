# Inter-die timing

A timing path that crosses dies passes through drivers, bumps or TSVs, and receivers, and its launch and capture clocks are distributed by separate clock trees on separate dies. The dies may come from different wafers, lots or even processes, so their variations are independent. Inter-die timing analysis must model the crossing delay, the clock latency on each die and the corner combinations that independence creates.

## Setup slack across a die boundary

Setup slack is the required time minus the arrival time. For a path launched on die A and captured on die B, the arrival is the launch clock latency $L_A$ plus the launch-die logic, the crossing (driver, bump or TSV, receiver) and the capture-die logic. The required time is the period plus the capture latency $L_B$ minus setup. With a 1000 ps period, $L_A = 300$, $L_B = 340$, 380 ps of logic on die A, 90 ps crossing, 260 ps on die B and 30 ps setup, the slack is $1000 + 340 - 30 - (300 + 380 + 90 + 260) = 280$ ps. Each die's clock tree is separate, so the difference $L_B - L_A$ is a skew that a die-level analysis alone never sees.

## Independent die corners

Dies made on different wafers do not share a process corner, so a path can see a slow launch die and a fast capture die. If die delays scale by $k_A$ and $k_B$, a single-corner analysis that assumes both dies slow overstates the delay of the capture-die logic by $(k_{slow} - k_B)$ times its delay. With 260 ps of capture-die logic, a slow factor of 1.15 and an actual capture factor of 0.9, the pessimism is $0.25 \times 260 = 65$ ps. Correct analysis enumerates corner combinations per die, which grows as the number of corners to the power of the number of dies, or uses statistical timing with independent die-level variables. Tools prune combinations that cannot be worst.

## Hold across dies

Hold checks need the earliest arrival and the latest capture clock. Hold slack is the minimum arrival minus the capture clock latency minus the hold time, all measured from the same edge: $\text{slack} = (L_A + d_{min}) - (L_B + t_{hold})$. With $L_A = 300$ ps, a minimum data delay of 95 ps including a fast crossing, $L_B = 380$ ps and 20 ps hold, the slack is $395 - 400 = -5$ ps, a violation. Skew between independently built clock trees on two dies is the usual cause, which is why 3D clocking either distributes one clock through dedicated vertical connections with matched delays, or retimes crossings through FIFOs or source-synchronous interfaces that tolerate skew.

## Going deeper

Lim, "Design for High Performance, Low Power, and Reliable 3D Integrated Circuits" (2013), chapters on 3D timing; Visweswariah and others, "First-order incremental block-based statistical timing analysis", DAC (2004).
