# The differential pair

The differential pair is the input stage of nearly every op-amp, comparator and OTA. It amplifies the difference between two inputs and rejects their common part, which makes it robust to supply noise and interference. Its symmetry is the reason analog layout tools treat differential pairs as their most important matching constraint.

## Splitting the tail current

Two matched transistors share a source node fed by a tail current $I_{SS}$. With equal inputs each carries $I_{SS}/2$. A small differential input $v_{id}$ moves current from one side to the other: $\Delta i = g_m v_{id}/2$ in each device, where $g_m$ is the transconductance of one device at $I_{SS}/2$. In gm/Id terms, $g_{m1} = (g_m/I_D)(I_{SS}/2)$: a 200 microampere tail at $g_m/I_D = 15$ gives $g_{m1} = 1.5$ mS. A common-mode input moves both gates together; ideally the tail holds the total current fixed and nothing changes at the outputs. Finite tail impedance and mismatch let some common-mode signal through, which is what CMRR measures.

## The five-transistor OTA

Loading the pair with a current mirror converts the differential current into a single-ended output. The mirror copies one side's current change to the output node, where it adds to the other side's opposite change, so the full $g_{m1}v_{id}$ flows into the output resistance $r_{o2} \parallel r_{o4}$. The differential gain is $A_d = g_{m1}(r_{o2} \parallel r_{o4})$. With $g_{m1} = 0.5$ mS and both output resistances 200 kilohms, $A_d = 0.5 \cdot 100 = 50$ V/V. This stage is the first stage of the classic two-stage op-amp. The mirror's own pole at its gate node introduces a pole-zero doublet, usually at high enough frequency to ignore in first-order design.

## Large-signal range

The linear model holds only for small inputs. In the square law, one device carries the whole tail current once $|v_{id}| \ge \sqrt{2}\,V_{OV}$, where $V_{OV}$ is the balanced overdrive. A pair at $V_{OV} = 0.15$ V saturates beyond 212 mV of differential input. Lower overdrive (higher $g_m/I_D$) gives more $g_m$ per unit current but a narrower linear range, so input pairs of comparators and low-noise amplifiers sit at high $g_m/I_D$ while pairs that must handle large swings sit lower. Source degeneration widens the range at the cost of gain. A sizing tool that optimizes only small-signal specs can produce a pair that slews or distorts badly for realistic inputs; large-signal specs such as slew rate and linearity must appear in the formulation too.

## Going deeper

Razavi, "Design of Analog CMOS Integrated Circuits" (2nd ed., 2017), chapter 4.
