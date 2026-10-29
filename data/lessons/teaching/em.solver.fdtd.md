# Finite-difference time-domain and the Courant limit

The finite-difference time-domain method (FDTD) discretizes Maxwell's curl equations on a staggered grid and steps electric and magnetic fields forward in time. One run with a broadband pulse yields the response over a wide frequency range, which makes it a workhorse for package, connector and antenna analysis. Its cost is set by two constraints: cell size for accuracy, and the Courant limit for stability.

## The Yee grid and cell size

Yee's scheme places electric and magnetic field components on interleaved grids offset by half a cell, and updates them in alternating half time steps (leapfrog), so each curl is a central difference. The grid must resolve the shortest wavelength present; a common rule is at least 10 to 20 cells per wavelength, since numerical dispersion makes waves travel at slightly wrong speeds on coarse grids.

The shortest wavelength is $\lambda_{min} = c/(f_{max}\sqrt{\varepsilon_r})$, so $\Delta x \le \lambda_{min}/N$. Example: up to 20 GHz in a dielectric with $\varepsilon_r = 4$, $\lambda_{min} = 7.49$ mm, and 20 cells per wavelength allow $\Delta x = 0.375$ mm. Fine geometric features, thin traces or small gaps, often force much smaller cells than the wavelength rule, which is where FDTD's cost comes from.

## The Courant stability limit

Explicit time stepping is stable only if a wave cannot cross more than one cell per step. For cubic cells of side $\Delta x$ in 3D the Courant-Friedrichs-Lewy condition is
$$\Delta t \le \frac{\Delta x}{c\sqrt 3},$$
and $\Delta x/(c\sqrt 2)$ for square cells in 2D; $c$ is the fastest wave speed anywhere in the grid, usually that of vacuum or air.

Example: 0.1 mm cubic cells give $\Delta t \le 0.1 \times 10^{-3}/(3.00 \times 10^8 \cdot 1.732) = 192.6$ fs. In 2D the limit would be 235.9 fs. Exceeding the limit by even a little makes the fields grow exponentially. Because the smallest cell sets the step for the whole grid, local refinement is expensive unless the solver uses subgridding or implicit schemes.

## Time steps and run cost

The number of time steps to simulate a duration $T$ is $\lceil T/\Delta t\rceil$, and each step touches every cell, so the cost is roughly cells times steps. Halving the cell size in 3D multiplies the cell count by 8 and the step count by 2, a factor of 16 in work.

Example: 0.1 mm cells at 0.95 of the Courant limit give $\Delta t = 182.95$ fs, so 1 ns takes 5466 steps. The run must last long enough for the response to decay: high-Q structures ring for many periods, and cutting the run short leaves spectral artifacts after the Fourier transform. Absorbing boundaries, such as perfectly matched layers, terminate the grid without reflections and add cells at the edges. These are the levers an agent tuning an FDTD run trades against accuracy.

## Going deeper

Taflove and Hagness, "Computational Electrodynamics: The Finite-Difference Time-Domain Method" (3rd ed.); Yee, "Numerical solution of initial boundary value problems involving Maxwell's equations in isotropic media", IEEE Transactions on Antennas and Propagation, 1966.
