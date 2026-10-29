# Solving the density Poisson equation by FFT

Electrostatic placement solves Poisson's equation for the potential on the bin grid at every iteration. A spectral method does it in $O(m^2 \log m)$ time for an $m \times m$ grid: expand the density in cosine modes, divide each mode by its Laplacian eigenvalue and transform back. This step is why the method scales, and it is the kernel GPU placers accelerate with fast cosine transforms.

## Cosine modes diagonalise the Laplacian

On a rectangular grid with zero-flux (Neumann) boundaries, the density can be expanded in cosine modes $\cos(\omega_u x)\cos(\omega_v y)$, with $\omega_u = \pi u/m$ for unit bins on an $m$-bin grid. Each mode is an eigenfunction of the Laplacian: $\nabla^2[\cos(\omega_u x)\cos(\omega_v y)] = -(\omega_u^2 + \omega_v^2)\cos(\omega_u x)\cos(\omega_v y)$. So if the density has coefficient $a_{uv}$ for a mode, the potential solving $\nabla^2\psi = -\rho$ has coefficient $a_{uv}/(\omega_u^2 + \omega_v^2)$. With m = 16, u = 6, v = 3 and $a = 0.3$, $\omega_u = 1.1781$ and $\omega_v = 0.5890$, so the potential coefficient is $0.3/1.7349 = 0.1729$. The constant mode $u = v = 0$ is dropped, which is the neutrality condition.

## The field from the same coefficients

The field $E = -\nabla\psi$ follows by differentiating the modes, so it too is a sum of sine-cosine products with the same coefficients:
$$E_x = \sum_{u,v} \frac{a_{uv}\,\omega_u}{\omega_u^2 + \omega_v^2}\sin(\omega_u x)\cos(\omega_v y),$$
and similarly for $E_y$ with the roles swapped. Each field component is one inverse transform of a mixed sine and cosine type. Low-frequency modes, with small $\omega$, get the largest amplification: they describe large-scale imbalance, such as all cells on the left half, and produce long-range forces, while high-frequency modes describe local clumps and produce short-range ones. That built-in multiscale behaviour is a large part of why electrostatic spreading works well.

## Why the FFT makes it scale

A direct transform computes $m^2$ coefficients, each a sum over $m^2$ bins: $m^4$ operations. A fast cosine transform does the same work in about $m^2 \log_2(m^2)$ operations, a ratio of $m^2/(2\log_2 m)$, which for $m = 512$ is about 14564. One placement iteration needs a forward transform for the coefficients and inverse transforms for the potential and the two field components, all on the bin grid, independent of the number of cells. Bin counts are usually powers of two near the square root of the cell count. Accumulating the density onto the grid and sampling the field at cells are scatter and gather operations; on GPUs, DREAMPlace implements the cosine transforms via FFT routines.

## Going deeper

Lu et al., "ePlace" (ACM TODAES, 2015), section on the spectral Poisson solver; Lin et al., "DREAMPlace" (DAC 2019), on DCT-based density computation.
