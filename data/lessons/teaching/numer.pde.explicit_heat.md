# Explicit time stepping of heat flow

Transient thermal analysis, such as a hot spot heating during a power burst, integrates the heat equation in time. Discretizing space first turns it into a large stiff ODE system, and the explicit time-step limit shows concretely why fine thermal meshes demand implicit integration, just like circuits.

## Method of lines: the grid becomes an ODE system

Discretize $T_t = \alpha T_{xx}$ in space only, with diffusivity $\alpha = k/(\rho c_p)$ in mm$^2$/s:
$$\frac{dT_i}{dt} = \frac{\alpha}{h^2}(T_{i-1} - 2T_i + T_{i+1}).$$
Each node becomes an ODE, and the grid becomes a system $T' = AT + s$ with a sparse matrix $A$, an RC network in which each node has a thermal capacitance. Example: $\alpha = 90$ mm$^2$/s (close to silicon at room temperature), $h = 0.5$ mm, temperatures 60, 64, 62 °C: $dT/dt = 90 \cdot (-6)/0.25 = -2160$ K/s at the middle node.

Any ODE integrator can now be applied. The eigenvalues of the 1D operator lie between about $-4\alpha/h^2$ and 0; the most negative one sets the stiffness.

## The explicit stability limit

Forward Euler in time (FTCS) is stable only if $h\lambda$ stays in its stability disk, which for the most negative eigenvalue $-4\alpha/h^2$ in 1D gives
$$r = \frac{\alpha\,\Delta t}{h^2} \le \frac12, \qquad\text{and in } d \text{ dimensions } r \le \frac{1}{2d}.$$
With $\alpha = 90$ mm$^2$/s and $h = 10$ µm in 3D: $\Delta t \le (0.01)^2/(6 \cdot 90)$ s $= 0.185$ µs.

The limit scales like $h^2$: halving the mesh spacing quarters the allowed step and multiplies the step count by 4, on top of 8 times the unknowns. A 10 ms power transient at that step needs 54,000 steps over every node of the mesh. Fine features such as TSVs or thin bond layers make this hopeless for explicit methods.

## One explicit step and why implicit is preferred

The FTCS update is $T_i \leftarrow T_i + r(T_{i-1} - 2T_i + T_{i+1})$. With $r = 0.4$ and temperatures 50, 56, 53 °C: $T_i = 56 + 0.4(-9) = 52.4$ °C. For $r \le 1/2$ the new value is a weighted average of the old ones with non-negative weights, so no new maxima appear, the discrete maximum principle; for $r > 1/2$ the middle weight goes negative and errors oscillate and grow.

Thermal solvers for chips therefore use implicit steps: backward Euler, unconditionally stable and damping, or Crank-Nicolson (the trapezoidal rule), second order but prone to the same ringing as in circuits when steps are large. Each implicit step is a sparse SPD solve, the same kind as a steady thermal or IR-drop solve, so the solver investment carries over.

## Going deeper

LeVeque, "Finite Difference Methods for Ordinary and Partial Differential Equations" (SIAM, 2007), chapter 9; Huang et al., "HotSpot: a compact thermal modeling methodology for early-stage VLSI design", IEEE Trans. VLSI (2006).
