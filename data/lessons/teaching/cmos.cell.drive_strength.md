# Drive strengths within a cell family

A library offers each function in several drive strengths, such as X1, X2, X4 and X8. Larger sizes drive loads faster but present more input capacitance to their own driver. Choosing a size is the most common timing optimization in synthesis and physical design, and its trade-off is a short RC calculation.

## Scaling within a family

Size Xk of a cell is, to first order, $k$ copies of X1 in parallel: its output resistance is $R_1/k$, and both its input capacitance and its output parasitic capacitance are $k$ times those of X1. Its delay driving a load $C_L$ is
$$t_k = \frac{R_1}{k}(k C_{p1} + C_L) = R_1 C_{p1} + \frac{R_1 C_L}{k}.$$
The first term is the parasitic delay, the same for every size. The second shrinks with $k$. With $R_1 = 8$ kOhm, $C_{p1} = 2$ fF and $C_L = 80$ fF: X1 takes $16 + 640 = 656$ ps, X4 takes $16 + 160 = 176$ ps and X16 takes $16 + 40 = 56$ ps. Logical effort is the same for every size, which is why it is a property of the topology.

## The smallest size that meets a target

A sizer usually wants the smallest cell that meets a delay or slew target, since smaller cells save area, leakage and input load. The delay $t_k$ decreases monotonically with $k$, so scan sizes in increasing order and take the first that meets the target. Continuing the example, a target of 200 ps rejects X1 (656 ps) and X2 (336 ps) and accepts X4 (176 ps). Libraries also give each output a max_capacitance and max_transition limit, and a size that meets delay but exceeds a limit is still rejected. In Liberty terms, each size has its own tables, so a real sizer looks up the actual delays rather than trusting the scaling rule.

## Upsizing loads the driver

Upsizing a cell multiplies its input capacitance, which slows the stage driving it. If the previous stage has resistance $R_0$ and drives a wire of capacitance $C_w$ into the cell, the two-stage delay is
$$T(k) = R_0(C_w + k C_{in1}) + R_1 C_{p1} + \frac{R_1 C_L}{k}.$$
The first term grows with $k$ and the last shrinks, so there is a best size, near $k = \sqrt{R_1 C_L / (R_0 C_{in1})}$ in the continuous case. With $R_0 = 6$ kOhm, $C_{in1} = 2$ fF, $R_1 = 8$ kOhm and $C_L = 80$ fF, the optimum is near $k = 7.3$, so X8 is the best discrete choice. This is logical-effort reasoning in absolute units, and it is why sizers evaluate a change together with its effect on the driver, not in isolation.

## Going deeper

Sutherland, Sproull and Harris, "Logical Effort" (1999), chapter 1; Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), chapter 8 on gate sizing.
