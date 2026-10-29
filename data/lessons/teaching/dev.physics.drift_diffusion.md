# The drift-diffusion model

The drift-diffusion model is the set of coupled partial differential equations that device TCAD solves: Poisson's equation for the potential and continuity equations for electrons and holes, with currents given by drift plus diffusion. Almost every I-V curve that calibrates a compact model before silicon comes from this model. A tool designer should know what it contains, why it is hard to solve numerically, and where its physics stops.

## Poisson plus continuity

The model has three unknown fields: electrostatic potential $\psi$, electron density $n$ and hole density $p$. Poisson's equation sets the potential from the net charge,
$$\nabla \cdot (\varepsilon \nabla \psi) = -q(p - n + N_D^+ - N_A^-).$$
Continuity equations conserve each carrier,
$$\frac{\partial n}{\partial t} = \frac{1}{q}\nabla \cdot J_n + G - R, \qquad \frac{\partial p}{\partial t} = -\frac{1}{q}\nabla \cdot J_p + G - R,$$
and the current equations close the system: $J_n = q\mu_n n E + qD_n \nabla n$, with the analogous hole expression.

Material models plug in mobility, generation and recombination. Boundary conditions apply at contacts (fixed potential and equilibrium densities at ohmic contacts) and at insulator interfaces (zero normal current). Terminal currents, capacitances and thresholds are computed from the solution afterwards; they are outputs, not unknowns.

## Discretization and solution

The equations are stiff and nonlinear. Carrier densities depend exponentially on potential, and the drift term multiplies two unknowns. Simulators discretize on an unstructured mesh with the finite-volume (box) method, which conserves current between control volumes.

Plain central differences for the current between two nodes become unstable once the potential step between nodes exceeds a few $kT/q$; densities oscillate and go negative. The Scharfetter-Gummel scheme fixes this by assuming the exact exponential density profile along each edge for constant field, which keeps densities positive.

Densities span around 20 orders of magnitude in one device, so solvers scale variables or use quasi-Fermi potentials. Gummel's decoupled iteration (Poisson, then electrons, then holes) is cheap and robust at low bias. At high bias the coupling is strong and solvers switch to a fully coupled Newton method on the sparse Jacobian, using bias ramping for initial guesses.

## Where drift-diffusion stops

Drift-diffusion assumes carriers relax instantly to the local field and stay at the lattice temperature, with $D = \mu kT/q$. These assumptions fail in several places a modern device designer cares about.

In channels of a few tens of nanometres, carriers cross before scattering enough to relax, so velocity overshoot raises current above the saturation-velocity prediction. Energy-balance (hydrodynamic) models add a carrier temperature, and Monte Carlo transport simulates scattering explicitly. Hot-carrier reliability needs the energy distribution, which drift-diffusion does not provide.

Near a 1 nm gate dielectric, quantum confinement moves the inversion charge peak about a nanometre away from the interface. Density-gradient corrections add a quantum potential to reproduce this. Direct tunneling through dielectrics and band-to-band tunneling need separate models. In production TCAD, drift-diffusion with these corrections and calibrated parameters remains the workhorse, because it is fast enough for design-technology exploration.

## Going deeper

Selberherr, "Analysis and Simulation of Semiconductor Devices" (1984); Vasileska, Goodnick and Klimeck, "Computational Electronics" (2010); Scharfetter and Gummel, IEEE Transactions on Electron Devices (1969).
