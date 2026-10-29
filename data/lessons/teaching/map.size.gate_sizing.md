# Gate sizing

After mapping, each cell's function is fixed but its drive strength is not. Gate sizing chooses sizes to meet timing at minimum area or power. Its smallest instance, one gate between a driver and a load, already shows the whole story: a convex trade-off between loading the driver and driving the load, an optimum in closed form, and a rounding loss when the library offers only discrete sizes.

## The two-stage trade-off

Let a gate of size s have input capacitance $c_0 s$ and drive resistance $r_0 / s$. A driver of resistance $R_d$ feeds it and it drives a load $C_L$. The two-stage delay is
$$D(s) = R_d c_0 s + p + \frac{r_0 C_L}{s}.$$
Setting $dD/ds = R_d c_0 - r_0 C_L / s^2 = 0$ gives $s^* = \sqrt{r_0 C_L / (R_d c_0)}$, and $D$ is convex for s > 0, so this is the minimum. With $R_d$ = 1.8 kOhm, $c_0$ = 1 fF, $r_0$ = 1.2 kOhm, $C_L$ = 52 fF: $s^* = \sqrt{62.4/1.8} = 5.89$.

## Discrete sizes

Libraries offer sizes like 1, 2, 4 and 8. Evaluate D at each and keep the minimum. With $R_d$ = 1.1 kOhm, $c_0$ = 4 fF, $r_0$ = 2.7 kOhm, p = 13 ps and $C_L$ = 10 fF: D(1) = 44.4, D(2) = 35.3, D(4) = 37.35 and D(8) = 51.58 ps, so size 2 wins. By convexity the best discrete size is one of the two neighbours of $s^*$, so only those need checking. In a full netlist the same trade-off appears at every gate, coupled through shared loads, and the problem becomes a discrete optimization over all gates.

## Continuous versus discrete sizing

Continuous sizing on posynomial delay models is a convex problem after a change of variables (geometric programming), so it can be solved exactly; TILOS-style greedy sizing and Lagrangian relaxation are the practical methods. Discrete sizing is NP-hard in general. The usual approach is to solve the continuous problem and round, then repair violations. Rounding costs little near the optimum because D is flat there: an optimum of 4.73 rounded to 4 may lose under 1 ps. It costs more when $s^*$ falls between widely spaced sizes or when rounding one gate shifts its neighbours' optima, which is why sizers iterate.

## Going deeper

Fishburn and Dunlop, "TILOS: a posynomial programming approach to transistor sizing", ICCAD 1985; Sapatnekar, "Timing" (2004), chapters on gate sizing; Kahng, Lienig, Markov and Hu, "VLSI Physical Design", chapter 8.
