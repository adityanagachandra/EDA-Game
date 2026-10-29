# Dynamic and domino logic

Dynamic logic replaces the pMOS pull-up network with a single clocked precharge device and evaluates through the nMOS network alone. It is faster and smaller than static CMOS, and it appears in high-performance datapaths, register files and memory peripheries. It is almost absent from standard-cell flows, and understanding why tells a tool designer exactly what the standard-cell abstraction assumes.

## Precharge and evaluate

A footed dynamic gate has a pMOS from $V_{DD}$ to the output node and an nMOS foot at the bottom of the pull-down network, both driven by the clock. While the clock is low (precharge), the pMOS charges the output to $V_{DD}$ and the foot is off, so no path to ground exists whatever the inputs. When the clock rises (evaluate), the precharge device turns off and the foot turns on; the output discharges if the nMOS network conducts and otherwise floats high on its capacitance. A dynamic NAND3 needs 3 series nMOS, the precharge pMOS and the foot: 5 transistors against 6 for static CMOS. Inputs drive only nMOS devices, so input capacitance and logical effort are lower, and the output makes only a falling transition during evaluation.

## Charge sharing

During evaluation the output node is floating if the pull-down network does not conduct. If some transistors in the stack turn on without completing a path to ground, the output shares its charge with uncharged internal nodes. By charge conservation, an output of capacitance $C_{out}$ precharged to $V_{DD}$ connected to an internal node $C_x$ at 0 V settles at
$$V = V_{DD}\frac{C_{out}}{C_{out} + C_x}.$$
With $C_{out} = 20$ fF, $C_x = 5$ fF and $V_{DD} = 1$ V, the output drops to 0.8 V. If that falls below the switching threshold of the gate it drives, the logic fails, even though the function said the output should stay high. Remedies are precharging internal nodes too, keeping internal capacitance small relative to the output, and adding a keeper.

## Monotonicity, domino and keepers

A dynamic output can fall during evaluation but cannot rise again until the next precharge, so inputs must only rise during evaluation. One dynamic gate cannot feed another directly, because the first output starts high and only falls. Domino logic adds a static inverter after each dynamic node, so every stage's output rises monotonically; the price is that each domino gate is non-inverting. A weak pMOS keeper fed back from the inverter's output holds the dynamic node high against leakage and noise, at the cost of a fight during evaluation. All of this, clock discipline, monotonic inputs, keeper ratios and charge-sharing checks, is outside what a Liberty arc and a static timer capture, which is why standard-cell flows stay with static CMOS.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), section 9.4; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 6.3.
