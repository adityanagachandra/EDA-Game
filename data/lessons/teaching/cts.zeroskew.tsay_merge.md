# Tsay's zero-skew tapping point

Path length ignores how the downstream load slows a wire. Tsay's method keeps the bottom-up zero-skew merge but uses Elmore delay, so each subtree is summarised by its delay and its total capacitance. The tapping point then follows from a closed-form equation that is exact under Elmore. Tsay's merge is the standard building block of zero-skew clock routers and shows how capacitance, not only length, moves the balance point.

## The tapping-point formula

Merge subtree 1 (delay $t_1$, load $C_1$) and subtree 2 ($t_2$, $C_2$) with a wire of length $L$, resistance $r$ and capacitance $c$ per unit length. Put the tap at fraction $x$ from subtree 1, and model each side of the wire as a $\pi$ segment. Setting $r xL(c xL/2 + C_1) + t_1$ equal to $r(1-x)L(c(1-x)L/2 + C_2) + t_2$, the quadratic terms cancel and
$$x = \frac{(t_2 - t_1) + rL(C_2 + cL/2)}{rL(cL + C_1 + C_2)}.$$
With $L = 200$ µm, $r = 0.1\ \Omega$/µm, $c = 0.14$ fF/µm, $C_1 = 75$, $C_2 = 25$ fF, $t_1 = 35$ and $t_2 = 36$ ps: the numerator is $1 + 20 \cdot 39/1000 = 1.78$ ps, the denominator $20 \cdot 128/1000 = 2.56$ ps, so $x = 0.695$.

## The merged node

The merged subtree has delay equal to either side's: with the numbers above, $xL = 139.1$ µm, and $t_1 + r\,xL\,(c\,xL/2 + C_1)/1000 = 35 + 13.91 \cdot 84.73/1000 = 36.18$ ps; the subtree-2 side gives the same 36.18 ps. Its capacitance is $C_1 + C_2 + cL = 128$ fF. These two numbers, delay and load, are all the next merge needs, so the method climbs the topology in one bottom-up pass. Note that the heavier subtree pulls the tap toward itself: a large $C_1$ makes the wire on side 1 expensive, so $x$ shrinks.

## Elongation under Elmore delay

If $t_2 - t_1 > rL(cL/2 + C_1)$, the formula gives $x > 1$: subtree 2 is too slow even with the tap at its root. The tap moves onto root 2 and the wire to subtree 1 is lengthened to $L'$ with $rL'(cL'/2 + C_1) = t_2 - t_1$, a quadratic $(rc/2)L'^2 + rC_1L' - (t_2 - t_1) = 0$ with delays converted to the same units. For $r = 0.15$, $c = 0.18$, $C_1 = 20$ fF and a 2 ps gap, $0.0135L'^2 + 3L' - 2000 = 0$ gives $L' = 289.5$ µm. The case $x < 0$ is symmetric. Because Elmore delay grows quadratically with length, elongation under Elmore needs less extra wire than the linear model would suggest for large gaps.

## Going deeper

Tsay, "An exact zero-skew clock routing algorithm", IEEE Transactions on Computer-Aided Design (1993); Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the clock chapter.
