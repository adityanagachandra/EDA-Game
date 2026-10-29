# Elmore delay of a transistor stack

A series transistor stack is a small RC ladder: every internal node between stacked devices carries diffusion capacitance that must also be discharged. The Elmore delay of the stack shows why tall stacks are slow, why input ordering matters, and what a characterizer is measuring when it reports different delays for different input pins of the same NAND.

## The stack as an RC ladder

Take the pull-down of a NAND$n$: $n$ nMOS in series from the output to ground, each of resistance $R_s$ when on. Number the internal nodes from the ground end, so node $j$ sits above $j$ transistors and has capacitance $C_j$; the output node, above all $n$, has $C_{out}$. When every transistor is on, the stack is an RC chain from the output down to ground, and each node discharges through the resistance between it and ground. That is the same structure as a wire's RC tree, rooted at ground instead of at a driver, so the Elmore formula applies: each capacitance is weighted by the resistance it shares with the path to ground.

## Elmore delay of the falling output

In the worst case every internal node starts charged high, so all of them must be discharged with the output. Node $j$ discharges through $j R_s$ and the output through $n R_s$:
$$T = \sum_{j=1}^{n-1} j R_s C_j + n R_s C_{out}.$$
Example, NAND3 with $R_s = 4$ kOhm, $C_1 = 3$ fF, $C_2 = 2$ fF and $C_{out} = 15$ fF: the internal terms are $4 \cdot 3 + 8 \cdot 2 = 28$ ps and the output term is $12 \cdot 15 = 180$ ps, for 208 ps. If the stack is sized with width $n$ per device, $R_s = R/n$ and the output term is $R C_{out}$, matching an inverter, but the internal terms remain, which is the parasitic delay a tall stack adds.

## Input ordering and pin-dependent delay

Which input arrives last decides which internal nodes are still charged. If the transistor nearest ground switches last, every node above it is still high and the full sum applies. If the transistor nearest the output switches last, all the devices below it are already on, the internal nodes have already discharged, and only the output term remains: the saving is $\sum_j j R_s C_j$. That is why designers place the latest-arriving signal on the device nearest the output, and why a characterized NAND3 has different arc delays from pins A, B and C. Liberty captures this as separate timing arcs per input pin; a synthesis or ECO tool that swaps logically equivalent pins to improve timing is exploiting exactly this asymmetry.

## Going deeper

Weste and Harris, "CMOS VLSI Design" (4th ed.), section 4.3.6; Rabaey, Chandrakasan and Nikolic, "Digital Integrated Circuits" (2nd ed.), section 6.2.
