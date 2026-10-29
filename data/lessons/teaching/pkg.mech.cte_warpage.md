# CTE mismatch and warpage

Warpage is the out-of-plane bending of a package, wafer or panel caused by CTE mismatch between layers bonded at an elevated temperature. A warped package lifts its corner balls off the board during reflow, opens or bridges joints, raises thermal interface thickness and misaligns fine-pitch bonds. The bimaterial strip analysis of Timoshenko gives the curvature and is the first model every package mechanical engineer uses.

## Curvature of a symmetric bilayer

Two layers bonded stress-free and then cooled by $\Delta T$ bend toward the layer with the larger CTE, which shrinks more. Timoshenko's analysis gives the curvature $\kappa = 1/\rho$. When both layers have the same thickness and modulus, it reduces to $\kappa = 3\Delta\alpha\Delta T/(2h)$, with $h$ the total thickness. For $\Delta\alpha = 14$ ppm/K, $\Delta T = 100$ K and $h = 1$ mm, $\kappa = 3 \times 0.0014/0.002 = 2.1$ per metre, a radius of curvature of 0.476 m. Thinner stacks curve more, as $1/h$, which is why thin packages and thinned dies warp more than thick ones.

## Timoshenko's general formula

For layer thicknesses $t_1, t_2$, moduli $E_1, E_2$, $m = t_1/t_2$, $n = E_1/E_2$ and $h = t_1 + t_2$, Timoshenko's result is $\kappa = \frac{6\Delta\alpha\Delta T(1+m)^2}{h[3(1+m)^2 + (1+mn)(m^2 + 1/(mn))]}$. A 0.3 mm silicon die (130 GPa) on 0.5 mm of mold-like material (20 GPa) with $\Delta\alpha = 12.4$ ppm/K and $\Delta T = 100$ K gives $m = 0.6$, $n = 6.5$ and $\kappa = 2.22$ per metre, a radius of 0.45 m. The formula shows that a stiff thin layer on a compliant thick one bends the stack less than equal layers would, and that the curvature is linear in $\Delta\alpha\Delta T$, so it reverses sign when the assembly is heated past its stress-free temperature.

## From curvature to warpage

Package warpage is quoted as the out-of-plane deflection over a span, not as curvature. A spherical cap with curvature $\kappa$ over span $L$ deflects $\delta = \kappa L^2/8$ for small curvature. With $\kappa = 2.1$ per metre over a 20 mm span, $\delta = 2.1 \times 0.0004/8 = 105$ um. Board-assembly standards and coplanarity limits are often on the order of 100 um or less for large BGAs, so this package is at the edge. Warpage also changes with temperature through reflow, sometimes changing shape from concave to convex, which is why warpage is measured across the full reflow profile by shadow moiré, and why fine-pitch bonding needs flat, CTE-matched parts.

## Going deeper

Timoshenko, "Analysis of bi-metal thermostats", Journal of the Optical Society of America (1925); Suhir, "Structural Analysis in Microelectronic and Fiber-Optic Systems" (1991).
