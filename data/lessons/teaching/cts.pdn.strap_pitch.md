# Sizing strap pitch from an IR budget

Straps hold the cell rails near the supply at regular intervals. Between two straps a rail sags under the current of the cells along it. The sag has a closed form in the continuum limit, and inverting it gives the largest strap pitch that meets an IR budget. This is how PDN planners turn a drop budget into a grid specification before any detailed analysis runs.

## The parabolic rail profile

Treat the rail between two straps at distance $P$ as a continuous resistor with resistance $r$ per unit length, drawing a uniform current $j$ per unit length, and both ends held at zero drop. Current conservation gives $d''(x) = -jr$, so
$$d(x) = \frac{jr}{2}\,x(P - x),$$
a parabola with its maximum $jrP^2/8$ at the midpoint. For $j = 0.01$ mA/µm, $r = 0.5\ \Omega$/µm and $P = 105$ µm, the worst drop is $0.01 \cdot 0.5 \cdot 11\,025/8 = 6.9$ mV. It is the continuous version of a doubly fed resistor ladder.

## The largest pitch for a budget

Invert the midpoint drop: $P_{max} = \sqrt{8\,\Delta V/(jr)}$. With a rail budget of 4.5 mV, $j = 0.02$ mA/µm and $r = 1.8\ \Omega$/µm, $P_{max} = \sqrt{36/0.036} = 31.6$ µm. Because the drop grows with $P^2$, halving the pitch quarters the rail drop while doubling the tracks the straps occupy. Planners split the total IR budget across levels: the straps, the vias and the rails each get a share, and the rail share sets the lowest strap pitch.

## The drop anywhere along the rail

The same formula gives the drop at any point: with $j = 0.03$ mA/µm, $r = 2\ \Omega$/µm and $P = 90$ µm, the drop 45 µm from a strap is $0.03 \cdot 45 \cdot 45 = 60.75$ mV. Drop near a strap is small and rises fast, so cells with high current are best placed near straps; IR-aware placement uses this. The uniform-current assumption is the weak point: clusters of high-power cells concentrate current, and the real profile must come from the full grid solve with per-instance currents.

## Going deeper

Kahng et al., "VLSI Physical Design" (2nd ed., 2022), the chapter on power and ground routing; Popovich, Mezhiba and Friedman, "Power Distribution Networks with On-Chip Decoupling Capacitors" (2008).
