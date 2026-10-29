# Golden and differential checks for EM engines

Nothing in an agent-first engine is promoted without passing golden and differential tests. For EM engines the goldens are unusually strong, because physics supplies invariants that every correct answer must satisfy and analytic structures with exact answers. Differential reports on reference designs then prove that a change moved only what it was meant to move.

## Physical invariants as tests

Every correct electrostatic solve yields a Maxwell capacitance matrix that is symmetric ($C_{ij} = C_{ji}$, reciprocity), has positive diagonal entries, non-positive off-diagonal entries and non-negative row sums (the ground capacitances). S-parameter models of passive structures must be passive, $|S_{11}|^2 + |S_{21}|^2 \le 1$ for a two-port driven at port 1 with port 2 matched, and reciprocal, $S_{12} = S_{21}$.

These checks are cheap enough to run on every output, not only on goldens. Example: a solver reporting $C_{12} = -4.1$ fF and $C_{21} = -3.6$ fF has a 13% asymmetry, far beyond a converged solver's tolerance, and should be flagged. A passive two-port with $|S_{11}|^2 + |S_{21}|^2 = 1.03$ at one frequency creates energy there and can make a transient simulation blow up. Analytic goldens complement invariants: parallel plates, the coax with $C' = 2\pi\varepsilon/\ln(b/a)$, and a matched line.

## Convergence under refinement

A single match with a golden can be a coincidence. Convergence tests refine the discretization and check that results settle at the expected rate. Example: a capacitance changes by 4% on the first refinement and 0.2% on the next: the result is converging and the coarse mesh was inadequate. A change that does not shrink, or that grows, signals a bug, a singular integral handled badly, or a mesh that misses a feature.

For FDTD, a cavity's resonant frequency compared with its analytic value should approach it as the cell size falls, with error shrinking roughly as the square of the cell size for the standard second-order scheme. For MoM, panel refinement near edges matters most, since charge concentrates there. A golden suite records not one result but the sequence, so that an engine change that silently slows convergence is caught even when the finest-mesh answer still passes.

## Differential reports on a reference design

A differential report compares an engine change against the last promoted version on reference designs, net by net. Its baseline holds per-net R and C (and S-parameters where relevant), plus runtime and memory. The report must show that intended changes appear and nothing else moves beyond tolerance.

Example: an agent edits the fringe-correction plug-in for metal 3. The report should show only metal 3 nets moving, by amounts consistent with the correction. If a pure speed-up passes all goldens but shifts 10% of reference nets by 0.5%, promotion is blocked until the shift is explained. And identical totals are not equivalence: a build whose per-net capacitances differ by up to 3% in both directions can change timing on specific paths even though aggregates match, because errors in different nets do not cancel on any single path. Comparing net by net is the whole point.

## Going deeper

Pozar, "Microwave Engineering" (4th ed.), chapter 4 on passivity and reciprocity; Taflove and Hagness, "Computational Electrodynamics" (3rd ed.), on numerical dispersion and convergence.
