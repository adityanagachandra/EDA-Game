# Timing the clock-gate enable

A clock gate's enable has to settle before the gate's clock pin rises, and that clock arrives earlier than the clock at the flops the gate drives. The enable path is therefore tighter than an ordinary data path between the same flops. This check couples CTS and timing closure: where the gate sits in the tree decides how much time its enable has, and CTS sometimes clones gates to fix it.

## The gating setup check

Let the gated flops have latency $L_f$ and the gate's subtree delay $D_g$, so the gate's clock pin sees $L_f - D_g$. An enable launched by a flop at latency $L_f$ with clock-to-Q $t_{cq}$ and logic delay $d_{en}$ meets the gating setup $t_{su}$ when
$$\text{slack} = T + (L_f - D_g) - t_{su} - (L_f + t_{cq} + d_{en}) \ge 0.$$
With $T = 850$, $L_f = 480$, $D_g = 190$, $t_{su} = 50$, $t_{cq} = 40$ and $d_{en} = 250$ ps, the slack is $850 + 290 - 50 - 770 = 320$ ps. The flop latencies cancel, but $D_g$ does not: it is time stolen from the enable.

## The deepest allowed subtree

Since the slack equals $T - t_{su} - t_{cq} - d_{en} - D_g$, the subtree behind the gate may be at most $T - t_{su} - t_{cq} - d_{en}$ deep. With $T = 650$, $t_{su} = 54$, $t_{cq} = 61$ and $d_{en} = 350$ ps, $D_g$ may not exceed 185 ps. This bound tells CTS how far up the tree a gate may sit: a gate high in the tree has a large $D_g$, which saves power because more buffers stop toggling but leaves little room for slow enable logic.

## Cloning gates toward the flops

If the enable fails, CTS or the optimiser can clone the gate and place a copy close to each group of flops, so each copy has a short subtree. Cutting $D_g$ from 150 to 50 ps raises the slack by 100 ps, from 52 to 152 ps in one example. The price is more gates, each with its own always-toggling input, and less power saved per gate. The reverse transformation, merging clones up the tree, saves power where the enable has slack. Both are policy decisions that a CTS engine makes from timing and power data.

## Going deeper

Keating et al., "Low Power Methodology Manual" (2007); the OpenSTA documentation on clock-gating checks.
