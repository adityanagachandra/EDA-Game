# Coupling capacitance

Two wires side by side form a capacitor through the dielectric between their facing sidewalls. This coupling capacitance does not go to ground: its other terminal is another signal, so its effect on delay and noise depends on what that neighbour is doing. As wires got taller and closer, coupling became the largest component of many nets' capacitance, and extracting it per neighbour, rather than lumping it into ground, is what makes signal-integrity analysis possible.

## Facing sidewalls

The facing sidewalls of two parallel wires of thickness $T$, spaced $s$ apart along a run of length $L$, behave roughly like parallel plates: $C_c \approx \varepsilon_0 \varepsilon_r T L / s$. Two wires 0.05 um thick at 0.05 um spacing for 30 um in $\varepsilon_r = 2.5$ give $0.008854 \times 2.5 \times 0.05 \times 30 / 0.05 = 0.66$ fF. Fringing between top corners adds to this, so field solvers report more. The formula shows why coupling grew with scaling: wires are kept tall to limit resistance (aspect ratios of 2 or more), so $T/s$ rises as pitch shrinks. Coupling is then often more than half of a minimum-pitch wire's total capacitance.

## Spacing dependence

Coupling per length falls as spacing grows, roughly as $1/s$ in the plate model, faster when the extra space lets field lines escape to ground instead. Characterization tables give coupling at several spacings, and the extractor interpolates. If a table gives 95 aF/um at 0.04 um and the model is $1/s$, two wires at 0.08 um couple 47.5 aF/um; over 250 um that is 11.875 fF. Doubling the spacing halves coupling but does not shrink total capacitance by as much, because some of the field lines that ended on the neighbour now end on the ground below as fringe. This is the physical basis of the spacing-based crosstalk fix that routers apply to critical nets.

## Parallel run length

Two wires couple only where they run side by side. The extractor finds, for each pair of neighbouring segments, the interval along the track direction where both exist: the overlap of their extents. A victim from $x = 0$ to 90 um and an aggressor from $x = 10$ to 150 um on the adjacent track overlap from 10 to 90, 80 um. At 150 aF/um this gives 12 fF. Computing overlaps for every pair is a geometric sweep: sort segments by position along each track and scan neighbouring tracks together, which is $O(n \log n)$ rather than all pairs. End effects, where the field wraps around a line end, add a small correction that pattern libraries handle separately.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), sections on interconnect parasitics; Sakurai, "Closed-form expressions for interconnection delay, coupling, and crosstalk in VLSIs", IEEE Transactions on Electron Devices (1993).
