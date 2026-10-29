# Multi-voltage design and voltage islands

Dynamic power scales with the square of the supply voltage, so running non-critical logic at a lower voltage is the strongest single power lever. A voltage island is a region with its own supply; dynamic voltage and frequency scaling (DVFS) changes a domain's voltage and clock at run time. Both create new work for the flow: separate grids, level shifters on every crossing, and timing at several voltage corners. A tool designer needs to quantify both the gain and the speed cost.

## Voltage scaling of dynamic power

With frequency and activity unchanged, $P \propto V^2$. A block burning 65 mW at 0.8 V burns $65 \cdot (0.65/0.8)^2 = 42.9$ mW at 0.65 V, a 34% saving from a 19% voltage cut. Leakage falls with voltage too, often more than linearly. The catch is speed: lower voltage makes every gate slower, so only logic with timing margin can move to the island without changing the frequency.

## Dynamic voltage and frequency scaling

When a workload does not need full speed, DVFS lowers both $V$ and $f$. Power scales as $V^2 f$: moving from 0.8 V at 2.25 GHz to 0.6 V at 1.125 GHz gives $(0.6/0.8)^2 \cdot 0.5 = 0.281$ of the original, a 71.9% reduction, while the work rate halves. Energy per operation scales as $V^2$, so DVFS saves energy, not just power. Each operating point is a separate timing and IR signoff scenario, which multiplies the analysis work.

## The speed cost

Gate delay rises as voltage approaches the threshold. The alpha-power law models it as $t_d \propto V/(V - V_t)^\alpha$, with $\alpha$ between 1 and 2 for short-channel devices. With $V_t = 0.4$ V and $\alpha = 1.3$, lowering from 0.85 to 0.65 V multiplies delay by $\frac{0.65/0.25^{1.3}}{0.85/0.45^{1.3}} = 1.64$. Near threshold the delay grows quickly and variation grows with it. This curve sets how low an island can go for a given frequency, and it is why signoff characterises libraries at each supply voltage.

## Going deeper

Keating et al., "Low Power Methodology Manual" (2007); Sakurai and Newton, "Alpha-power law MOSFET model and its applications to CMOS inverter delay", IEEE Journal of Solid-State Circuits (1990).
