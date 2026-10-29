# Finite difference discretization

Thermal analysis of a die, substrate noise and some field problems are partial differential equations. Finite differences replace derivatives by differences on a grid, turning the PDE into a sparse linear system that looks exactly like a resistor network. That equivalence is why thermal solvers reuse power-grid solvers, and why grid spacing controls both accuracy and cost.

## The second-difference stencil

Taylor expansion of $u(x \pm h)$ gives
$$u''(x) = \frac{u(x - h) - 2u(x) + u(x + h)}{h^2} - \frac{h^2}{12}u''''(\xi).$$
Dropping the error term gives the three-point stencil. Example: temperatures 40.0, 38.5 and 37.4 °C at spacing 0.5 mm: $T'' \approx (40.0 - 77.0 + 37.4)/0.25 = 1.6$ K/mm$^2$.

The stencil is exact for quadratics, since the fourth derivative vanishes. Positive $T''$ at a point means it is cooler than the average of its neighbours, so heat flows in. In 2D, adding the $x$ and $y$ stencils gives the five-point Laplacian, and in 3D the seven-point one.

## The five-point Laplacian on a thermal grid

Steady conduction with volumetric heat source $q$ (W/m$^3$) obeys $k\nabla^2T + q = 0$. On a square grid of spacing $h$, the five-point stencil at a node gives
$$T_c = \frac{T_N + T_S + T_E + T_W + qh^2/k}{4}.$$
Example: silicon, $k = 150$ W/(m K), $h = 100$ µm, $q = 3 \times 10^9$ W/m$^3$: $qh^2/k = 3 \times 10^9 \cdot 10^{-8}/150 = 0.2$ K. With neighbours at 60, 62, 61 and 63 °C, $T_c = (246 + 0.2)/4 = 61.55$ °C.

Multiplying the stencil equation through by the cell thickness shows it is KCL for a network of thermal conductances between neighbours with a heat current injected at the node. So a thermal grid is a resistive grid, and the same sparse solvers, CG and multigrid, apply.

## Truncation error and mesh refinement

The leading truncation error is $O(h^2)$, so halving $h$ cuts the discretization error by about 4. If a 50 µm grid shows a 2.4 K error at a hot spot compared with a reference, a 25 µm grid should show about 0.6 K.

Refinement is expensive: halving $h$ in 3D multiplies unknowns by 8, and the condition number grows like $h^{-2}$, so iterative solvers slow down too. Tools therefore refine only where gradients are steep, near hot spots, TSVs and interfaces, with nonuniform or adaptive meshes, and use the $h^2$ rate to extrapolate: two solutions at $h$ and $h/2$ give a Richardson estimate $(4T_{h/2} - T_h)/3$ that is more accurate than either. The rate also gives a test: if refinement does not reduce the error by about 4, the solution is not smooth there or the code has a bug.

## Going deeper

LeVeque, "Finite Difference Methods for Ordinary and Partial Differential Equations" (SIAM, 2007), chapters 2 and 3; Pedram and Nazarian, "Thermal modeling, analysis, and management in VLSI circuits: principles and methods", Proc. IEEE (2006).
