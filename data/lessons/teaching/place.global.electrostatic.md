# Electrostatics-based placement

Electrostatics-based placement, introduced by ePlace and extended by RePlAce and DREAMPlace, models every cell as a positive charge equal to its area. Density overflow becomes electric potential energy, and the density gradient becomes the electric force that pushes charges apart. The analogy gives a global, smooth density penalty that handles standard cells and macros together.

## Cells as charges, density as energy

Each cell $i$ is a charge $q_i$ equal to its area, spread over its footprint. The charge density over the region defines a potential $\psi$ through Poisson's equation, $\nabla^2\psi = -\rho$, and the density penalty is the system's potential energy:
$$N(x) = \frac{1}{2}\sum_i q_i \psi_i.$$
For charges 4, 1 and 6 at potentials 5.6, 3.6 and 0.5, $N = (22.4 + 3.6 + 3)/2 = 14.5$. Energy is high when charge is concentrated, because nearby charges raise each other's potential, and low when it is spread evenly, so minimising $W + \lambda N$ trades wirelength against uniform spreading. To make the even distribution the minimum, the density's mean is subtracted (a neutral system).

## Density force as electric force

The gradient of the energy with respect to a cell's position is its charge times the electric field, $\nabla_i N = -q_i E_i$ with $E = -\nabla\psi$. So the density force on a cell is $q_i E_i$, pointing from dense regions to sparse ones. A cell of area 5 in a field $E = (1.5, 0.7)$ feels a force of magnitude $5\sqrt{1.5^2 + 0.7^2} = 8.276$. Larger cells, such as macros, feel proportionally larger forces, so mixed-size placement needs no special treatment, though a preconditioner that divides each cell's gradient by its pin count plus $\lambda$ times its charge keeps steps balanced across sizes. The field is computed once per iteration for the whole grid, then sampled at every cell.

## Fillers and the overflow target

If the cells cover less area than the target density allows, the electrostatic force would keep spreading them across all the whitespace, lengthening wires. ePlace inserts filler cells, charges with no connections, so that movable plus filler area equals the target density times the free area. With 24100 um$^2$ of free area, target 0.7 and 9640 um$^2$ of movable cells, the fillers need $16870 - 9640 = 7230$ um$^2$, which is 1446 fillers of 5 um$^2$. Fillers move freely and occupy whitespace, so real cells stop spreading once the density is uniform at target. They are removed before legalization. The placer stops when the overflow of real cells reaches its target, about 0.1.

## Going deeper

Lu et al., "ePlace: electrostatics-based placement using fast Fourier transform and Nesterov's method" (ACM TODAES, 2015); Cheng, Kahng, Kang and Wang, "RePlAce" (IEEE TCAD, 2019).
