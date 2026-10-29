# Finite volume and flux conservation

The finite volume method integrates the conservation law over each cell, so heat leaving one cell exactly enters its neighbour. On chips and packages, where silicon, copper, oxide and thermal interface materials with conductivities differing by factors of hundreds meet, that exact conservation and the right interface rule matter more than formal order of accuracy.

## Balancing fluxes over a control volume

Integrate $-\nabla\cdot(k\nabla T) = q$ over a cell and apply the divergence theorem: the sum of heat flows out through the faces equals the heat generated inside. Approximating each face flow as a conductance times a temperature difference,
$$q_f = \frac{kA_f}{d}(T_P - T_N),$$
with $A_f$ the face area and $d$ the distance between cell centres. Example: $k = 150$ W/(m K), face 0.1 mm$^2$, centres 0.2 mm apart, temperatures 70 and 66 °C: $q_f = 150 \cdot 0.5 \cdot 4 = 300$ mW (with millimetre units, W/(m K) times mm gives mW/K).

Because the flux leaving one cell through a face is the same number that enters the neighbour, the discrete system conserves energy exactly, whatever the mesh.

## Harmonic averaging at material interfaces

When a face separates materials with conductivities $k_1$ and $k_2$, at distances $d_1$ and $d_2$ from the face, the two half-paths are thermal resistances in series:
$$G_f = \frac{A}{d_1/k_1 + d_2/k_2}.$$
With equal distances this gives the harmonic mean $k_f = 2k_1k_2/(k_1 + k_2)$. For an oxide-like layer, $k_1 = 1.4$, against silicon, $k_2 = 150$: $k_f = 2.77$ W/(m K), close to the poor conductor, as physics requires. The arithmetic mean, 75.7, would overstate heat flow through the barrier by a factor of 27.

This is the most common error in hand-written thermal solvers. It also applies to thin interface materials and bond layers, which can dominate the thermal resistance of a 3D stack.

## Conservation on nonuniform meshes

With face conductances $G_f$ to neighbouring cells at temperatures $T_f$ and power $P$ generated in the cell, the balance $\sum_f G_f(T_P - T_f) = P$ gives
$$T_P = \frac{\sum_f G_fT_f + P}{\sum_f G_f}.$$
Example: conductances 2, 1 and 3 mW/K to cells at 60, 70 and 50 °C, and 40 mW generated: $T_P = (120 + 70 + 150 + 40)/6 = 63.3$ °C.

Nothing here needs a uniform grid: cells can be thin in the stacking direction and wide laterally, which suits layered chip stacks, and the matrix stays symmetric and diagonally dominant. Compact thermal models of chips are essentially coarse finite-volume networks, with each block or layer a cell, and fine-grained thermal simulators refine the same structure.

## Going deeper

Patankar, "Numerical Heat Transfer and Fluid Flow" (Hemisphere, 1980), chapters 3 and 4; LeVeque, "Finite Volume Methods for Hyperbolic Problems" (Cambridge, 2002), chapter 1 for the conservation viewpoint.
