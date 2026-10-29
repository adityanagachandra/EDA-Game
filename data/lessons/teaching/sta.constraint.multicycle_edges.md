# Multicycle paths and the hold adjustment

A multicycle path gives data more than one cycle to settle. The setup half of the exception is intuitive; the hold half is where mistakes happen, because moving the setup check also moves the default hold check. Getting the edge arithmetic right is a routine but essential skill, and timers implement it precisely as described here for paths within one clock.

## Moving the setup edge

For a path between flops on one clock of period $T$, set_multicycle_path $N$ -setup moves the setup capture edge from $T$ to $N T$. The setup slack becomes $S = N T - t_{su} - (t_{cq} + d_{max})$ with ideal clocks. For $T = 400$ ps and $N = 2$, a path with $t_{cq} = 68$, $d_{max} = 580$ and $t_{su} = 34$ ps has slack $800 - 34 - 648 = 118$ ps, where it would have failed by 282 ps as a single-cycle path. The design must really sample the data only every $N$ cycles, usually through an enable on the capture flop, or the relaxed check is a lie.

## The default hold check moves too

By default the hold check uses the capture edge one before the setup capture edge. After -setup $N$, that is $(N-1) T$, not 0. With only the setup multicycle, hold slack is $S = t_{cq} + d_{min} - ((N-1)T + t_h)$. For $T = 350$, $N = 2$, $t_{cq} = 38$, $d_{min} = 62$ and $t_h = 33$ ps, the hold check sits at 350 ps and the slack is $100 - 383 = -283$ ps: a huge, spurious hold violation. A hold fixer that trusted it would insert hundreds of picoseconds of buffers and then break the very setup margin the multicycle path was meant to grant.

## Pulling the hold check back

set_multicycle_path $M$ -hold moves the hold check $M$ cycles earlier, from $(N-1)T$ to $(N-1-M)T$. The usual pairing is $M = N - 1$, which returns the hold check to 0, the same-edge check of a single-cycle path. For $N = 4$ and $M = 3$ with $T = 300$ ps, the hold edge is at 0 and a path with $t_{cq} = 60$, $d_{min} = 39$ and $t_h = 21$ ps has hold slack 78 ps. A smaller $M$ leaves hold at an intermediate edge, which is occasionally intended. For paths between clocks of different periods, -start and -end choose whether the shift counts launch or capture clock cycles, and the arithmetic is done in those cycles.

## Going deeper

Bhasker and Chadha, "Static Timing Analysis for Nanometer Designs" (2009), chapter 8 on multicycle paths; the Synopsys Design Constraints format specification, set_multicycle_path.
